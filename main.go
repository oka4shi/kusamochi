package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"time"

	"github.com/Khan/genqlient/graphql"
)

type dateRange struct {
	To   time.Time
	From time.Time
}

type person struct {
	Name          string
	Contributions int
}

type authedTransport struct {
	key     string
	wrapped http.RoundTripper
}

func (t *authedTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("Authorization", "bearer "+t.key)
	return t.wrapped.RoundTrip(req)
}

func main() {
	client, err := createClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	f, err := os.Open("users.json")
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
	for i, user := range users {
		contributions, err := getLastWeekContributions(client, time.Now(), user)
		if err != nil {
			log.Println(err)
		}

		var contributionsSum int
		for _, v := range contributions {
			contributionsSum += v.ContributionCount
		}
		data = append(data, person{
			Name:          user,
			Contributions: contributionsSum,
		})
		if i == 0 {
			duration = dateRange{
				From: contributions[0].Date,
				To:   contributions[len(contributions)-1].Date,
			}
		}
	}

	sort.Slice(data, func(i, j int) bool { return data[i].Contributions > data[j].Contributions })

	fmt.Printf("先週(%s～%s)のGitHubのContribution数ランキングをお知らせします！\n\n", formatDate(&duration.From), formatDate(&duration.To))
	for i, p := range data {
		fmt.Printf("%d位: %s (%d contributions)\n", i+1, p.Name, p.Contributions)
	}

}

func formatDate(t *time.Time) string {
	return t.Format("2006/01/02")
}

func createClient() (graphql.Client, error) {
	var err error
	ghToken := os.Getenv("KUSAMOCHI_GITHUB_TOKEN")
	if ghToken == "" {
		err = fmt.Errorf("must set KUSAMOCHI_GITHUB_TOKEN")
		return nil, err
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     ghToken,
			wrapped: http.DefaultTransport,
		},
	}

	return graphql.NewClient("https://api.github.com/graphql", &httpClient), nil
}

//go:generate go run github.com/Khan/genqlient genqlient.yaml
