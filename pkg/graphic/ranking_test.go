package graphic

import (
	"bytes"
	"fmt"
	"image/png"
	"io"
	"os"
	"testing"

	"github.com/oka4shi/kusamochi/pkg/discord"
)

func TestDrawRanking(t *testing.T) {
	img, err := DrawRanking([]Ranking{
		{
			Time: "4/1",
			Ranking: []RankingItem{
				{Name: "Alice", Value: "100"},
				{Name: "Bob", Value: "200"},
				{Name: "Charlie", Value: "300"},
			},
		},
		{
			Time: "3/1",
			Ranking: []RankingItem{
				{Name: "Bob", Value: "200"},
				{Name: "Alice", Value: "100"},
				{Name: "Charlie", Value: "300"},
			},
		},
		{
			Time: "2/1",
			Ranking: []RankingItem{
				{Name: "Bob", Value: "200"},
				{Name: "Alice", Value: "100"},
				{Name: "Charlie", Value: "300"},
			},
		},
		{
			Time: "1/1",
			Ranking: []RankingItem{
				{Name: "Charlie", Value: "300"},
				{Name: "Alice", Value: "100"},
				{Name: "Bob", Value: "200"},
			},
		},
	}, 3)
	if err != nil {
		t.Fatal(err)
	}

	stream := new(bytes.Buffer)

	png.Encode(stream, *img)

	hookURL := os.Getenv("KUSAMOCHI_WEBHOOK_URL")
	if hookURL == "" {
		return
	}

	files := []discord.File{
		{
			Name:        "ranking.png",
			Content:     stream,
			Description: "ファイルです",
		},
	}

	res, err := discord.PostWithFiles(hookURL, "画像をアップロードします", files)
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("Status: ", res.Status)
	fmt.Printf("Status: %+v\n", string(body))

}
