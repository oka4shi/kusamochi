package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"image/png"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/oka4shi/kusamochi/pkg/discord"
	"github.com/oka4shi/kusamochi/pkg/github"
	"github.com/oka4shi/kusamochi/pkg/graphic"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type dateRange struct {
	From time.Time
	To   time.Time
}

type weeklyData struct {
	Time dateRange
	Data []person
}

type person struct {
	Name          string
	Contributions int
}

func main() {
	err := run()
	if err != nil {
		log.Fatalln(err)
	}
}

func run() error {
	var err error
	ghToken := os.Getenv("KUSAMOCHI_GITHUB_TOKEN")
	if ghToken == "" {
		return errors.New("must set KUSAMOCHI_GITHUB_TOKEN")
	}
	client := github.CreateGraphQLClient(ghToken)

	hookURL := os.Getenv("KUSAMOCHI_WEBHOOK_URL")
	if hookURL == "" {
		return errors.New("must set KUSAMOCHI_WEBHOOK_URL")
	}

	jsonPath := os.Getenv("KUSAMOCHI_JSON_PATH")
	if jsonPath == "" {
		jsonPath = "users.json"
	}

	f, err := os.Open(jsonPath)
	if err != nil {
		return err
	}

	jsonStr, err := io.ReadAll(f)
	if err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		log.Printf("failed to close file: %v", err)
	}

	var users []string
	err = json.Unmarshal([]byte(jsonStr), &users)
	if err != nil {
		return err
	}

	var data []weeklyData
	var skipped []string

	for _, user := range users {
		contributions, err := github.GetWeeklyContributions(client, time.Now(), user, 6)

		var errList gqlerror.List
		if errors.As(err, &errList) {
			log.Println(err)
			skipped = append(skipped, user+"さん")
			continue
		} else if err != nil {
			// TODO: check graphql.HTTPError

			body := "GitHub APIの呼び出しに失敗しました"
			_, err2 := discord.Post(hookURL, body)
			if err2 != nil {
				return fmt.Errorf("failed to send error message to Discord: %w", err2)
			}

			return err
		}

		for j, weekly := range contributions {
			contributionsSum := 0
			for _, day := range weekly {
				contributionsSum += day.ContributionCount
			}
			p := person{
				Name:          user,
				Contributions: contributionsSum,
			}

			if len(data) <= j {
				data = append(data, weeklyData{
					Time: dateRange{
						From: weekly[0].Date,
						To:   weekly[len(weekly)-1].Date,
					},
					Data: []person{p},
				})
			} else {
				data[j].Data = append(data[j].Data, p)
			}
		}
	}

	for i := range data {
		sort.Slice(data[i].Data, func(j, k int) bool { return data[i].Data[j].Contributions > data[i].Data[k].Contributions })
	}

	latest := data[len(data)-1]

	body := ""
	body += fmt.Sprintf("先週(%s～%s)のGitHubのContribution数ランキングをお知らせします！\n\n", formatDate(&latest.Time.From), formatDate(&latest.Time.To))
	for i, p := range latest.Data {
		body += fmt.Sprintf("%d位: %s (%d contributions)\n", i+1, p.Name, p.Contributions)
	}
	if len(skipped) > 0 {
		body += fmt.Sprintf("\n%v のデータは取得に失敗したためランキングに含まれていません", strings.Join(skipped, "、"))
	}

	// Draw a ranking graph
	rankings := []graphic.Ranking{}
	for i := range data {
		d := data[len(data)-i-1]

		r := []graphic.RankingItem{}
		for _, p := range d.Data {
			r = append(r, graphic.RankingItem{
				Name:  p.Name,
				Value: p.Contributions,
			})
		}

		rankings = append(rankings, graphic.Ranking{
			Time:    d.Time.From.Format("1/2"),
			Ranking: r,
		})
	}

	stream, err := getRankingPng(rankings)
	if err != nil {
		body += "画像は生成に失敗したため送信しません"
		_, err := discord.Post(hookURL, body)
		if err != nil {
			return fmt.Errorf("failed to send error message to Discord: %w", err)
		}
	}

	files := []discord.File{
		{
			Name:        "ranking.png",
			Description: "ランキングの推移の画像",
			Content:     stream,
		},
	}
	resp, err := discord.PostWithFiles(hookURL, body, files)
	if err != nil {
		return fmt.Errorf("failed to send error message to Discord: %w", err)
	}
	if err := resp.Body.Close(); err != nil {
		log.Printf("failed to close response body: %v", err)
	}

	return nil
}

func formatDate(t *time.Time) string {
	return t.Format("2006/01/02")
}

func getRankingPng(rankings []graphic.Ranking) (*bytes.Buffer, error) {
	image, err := graphic.DrawRanking(rankings, len(rankings[0].Ranking))
	if err != nil {
		return nil, err
	}

	stream := new(bytes.Buffer)
	if err := png.Encode(stream, *image); err != nil {
		return nil, err
	}

	return stream, nil
}

//go:generate go run github.com/Khan/genqlient ../../genqlient.yaml
