package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/Khan/genqlient/graphql"
)

type dateRange struct {
	To   time.Time
	From time.Time
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
	var err error
	defer func() {
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}()

	ghToken := os.Getenv("KUSAMOCHI_GITHUB_TOKEN")
	if ghToken == "" {
		err = fmt.Errorf("must set KUSAMOCHI_GITHUB_TOKEN")
		return
	}

	httpClient := http.Client{
		Transport: &authedTransport{
			key:     ghToken,
			wrapped: http.DefaultTransport,
		},
	}

	graphqlClient := graphql.NewClient("https://api.github.com/graphql", &httpClient)
	contributions, err := getLastWeekContributions(graphqlClient, "oka4shi")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v", contributions)

}

func getDateRange(daysBefore int) dateRange {
	now := time.Now()

	return dateRange{
		To:   now,
		From: now.AddDate(0, 0, -daysBefore),
	}
}

type weeklyContributions = getUserContributionsUserContributionsCollectionContributionCalendarWeeksContributionCalendarWeek

func getLastWeekContributions(c graphql.Client, user string) (weeklyContributions, error) {
	r := getDateRange(7)

	var resp *getUserContributionsResponse
	resp, err := getUserContributions(context.Background(), c, user, r.To, r.From)
	if err != nil {
		return weeklyContributions{}, err
	}

	return resp.User.ContributionsCollection.ContributionCalendar.Weeks[0], nil
}

//go:generate go run github.com/Khan/genqlient genqlient.yaml
