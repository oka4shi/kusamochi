package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/oka4shi/kusamochi/pkg/github"
	"github.com/oka4shi/kusamochi/pkg/webhook"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type dateRange struct {
	From time.Time
	To   time.Time
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
	defer f.Close()

	jsonStr, err := io.ReadAll(f)
	if err != nil {
		log.Fatalln(err)
	}

	var users []string
	err = json.Unmarshal([]byte(jsonStr), &users)
	if err != nil {
		log.Fatalln(err)
	}

	var duration dateRange
	var data []person
	var skipped []string
	for _, user := range users {
		contributions, err := github.GetLastWeekContributions(client, time.Now(), user)

		var errList gqlerror.List
		if errors.As(err, &errList) {
			log.Println(err)
			skipped = append(skipped, user+"さん")
			continue
		} else if err != nil {
			// TODO: check graphql.HTTPError

			body := "GitHub APIの呼び出しに失敗しました"
			_, err2 := webhook.Post(hookURL, body)
			if err2 != nil {
				log.Println(err)
				log.Fatalln(err2)
			}

			log.Fatalln(err)
		}

		var contributionsSum int
		for _, v := range contributions {
			contributionsSum += v.ContributionCount
		}
		data = append(data, person{
			Name:          user,
			Contributions: contributionsSum,
		})

		duration = dateRange{
			From: contributions[0].Date,
			To:   contributions[len(contributions)-1].Date,
		}
	}

	sort.Slice(data, func(i, j int) bool { return data[i].Contributions > data[j].Contributions })

	body := ""
	body += fmt.Sprintf("先週(%s～%s)のGitHubのContribution数ランキングをお知らせします！\n\n", formatDate(&duration.From), formatDate(&duration.To))
	for i, p := range data {
		body += fmt.Sprintf("%d位: %s (%d contributions)\n", i+1, p.Name, p.Contributions)
	}
	if len(skipped) > 0 {
		body += fmt.Sprintf("\n%v のデータは取得に失敗したためランキングに含まれていません", strings.Join(skipped, "、"))
	}

	response, err := webhook.Post(hookURL, body)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Response status of Webhook:", response.Status)
}

func formatDate(t *time.Time) string {
	return t.Format("2006/01/02")
}

//go:generate go run github.com/Khan/genqlient ../../genqlient.yaml
