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
	"strconv"
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
	var err error
	ghToken := os.Getenv("KUSAMOCHI_GITHUB_TOKEN")
	if ghToken == "" {
		err = fmt.Errorf("must set KUSAMOCHI_GITHUB_TOKEN")
		log.Fatalln(err)
	}
	client := github.CreateGraphQLClient(ghToken)

	hookURL := os.Getenv("KUSAMOCHI_WEBHOOK_URL")
	if hookURL == "" {
		err = fmt.Errorf("must set KUSAMOCHI_WEBHOOK_URL")
		log.Fatalln(err)
	}

	jsonPath := os.Getenv("KUSAMOCHI_JSON_PATH")
	if jsonPath == "" {
		jsonPath = "users.json"
	}

	f, err := os.Open(jsonPath)
	if err != nil {
		log.Fatalln(err)
	}

	jsonStr, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}
	err = f.Close()
	if err != nil {
		fmt.Errorf("failed to close file: %v", err)
	}

	var users []string
	err = json.Unmarshal([]byte(jsonStr), &users)
	if err != nil {
		log.Fatalln(err)
	}

	var duration dateRange
	var data []weeklyData
	var skipped []string
	for i, user := range users {
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
				log.Println(err)
				log.Fatalln(err2)
			}

			log.Fatalln(err)
		}

		for j, c := range contributions {
			var contributionsSum int
			for _, v := range c {
				contributionsSum += v.ContributionCount
			}
			p := person{
				Name:          user,
				Contributions: contributionsSum,
			}

			if i == 0 {
				data = append(data, weeklyData{
					Time: dateRange{
						From: c[0].Date,
						To:   c[len(c)-1].Date,
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

	body := ""
	body += fmt.Sprintf("先週(%s～%s)のGitHubのContribution数ランキングをお知らせします！\n\n", formatDate(&duration.From), formatDate(&duration.To))
	for i, p := range data[0].Data {
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
				Value: strconv.Itoa(p.Contributions),
			})
		}

		rankings = append(rankings, graphic.Ranking{
			Time:    d.Time.From.Format("1/2"),
			Ranking: r,
		})
	}
	image, err := graphic.DrawRanking(rankings, len(rankings[0].Ranking))
	if err != nil {
		body += "画像は生成に失敗したため送信しません"
		response, err := discord.Post(hookURL, body)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Response status of Webhook:", response.Status)
	}

	stream := new(bytes.Buffer)
	png.Encode(stream, *image)
	files := []discord.File{
		{
			Name:        "ranking.png",
			Description: "ランキングの推移の画像",
			Content:     stream,
		},
	}
	discord.PostWithFiles(hookURL, body, files)
}

func formatDate(t *time.Time) string {
	return t.Format("2006/01/02")
}

//go:generate go run github.com/Khan/genqlient ../../genqlient.yaml
