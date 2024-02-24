package main

import (
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
	client, err := createClient()
	if err != nil {
		fmt.Println(err)
		return
	}

	contributions, err := getLastWeekContributions(client, time.Now(), "oka4shi")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%v", contributions)

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
