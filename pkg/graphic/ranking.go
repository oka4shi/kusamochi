package graphic

import (
	"errors"
	"image"

	"github.com/fogleman/gg"
)

type Ranking struct {
	Time    string
	Ranking []RankingItem
}

type RankingItem struct {
	Name  string
	Value string
}

type rankingItem struct {
	Rank  int
	Value string
}

type rankingByNameType struct {
	Name    string
	History []rankingItem
}

func DrawRanking(rankings []Ranking, count int) (*image.Image, error) {
	if len(rankings) == 0 {
		return nil, errors.New("Rankings is empty")
	}

	// rankingByName := lineUpRankingByName(rankings, count)

	dc := gg.NewContext(800, 600)

	dc.SetLineWidth(10)
	dc.DrawCircle(400, 300, 200)
	dc.SetRGB(0, 0, 0)
	dc.Stroke()

	image := dc.Image()

	return &image, nil
}

func lineUpRankingByName(rankings []Ranking, count int) []rankingByNameType {
	rankingByName := []rankingByNameType{}
	for i := range count {
		item := rankings[0].Ranking[i]
		rankingByName = append(rankingByName,
			rankingByNameType{
				Name: item.Name,
				History: []rankingItem{
					{
						Rank:  i,
						Value: item.Value,
					},
				},
			},
		)
	}

	for i, ranking := range rankings {
		if i == 0 {
			continue
		}

		for j := range count {
			item := ranking.Ranking[j]
			for _, r := range rankingByName {
				if r.Name == item.Name {
					r.History = append(r.History, rankingItem{
						Rank:  j,
						Value: item.Value,
					})
				} else {
					r.History = append(r.History, rankingItem{
						Rank:  -1,
						Value: item.Value,
					})
				}
			}
		}
	}
	return rankingByName
}
