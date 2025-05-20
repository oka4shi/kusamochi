package graphic

import (
	"errors"
	"fmt"
	"image"
	"slices"

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

type RankingGraphParams struct {
	minX      float64
	minY      float64
	marginX   float64
	marginY   float64
	radius    float64
	lineWidth float64
	gapX      float64
	gapY      float64
}

func DrawRanking(rankings []Ranking, count int) (*image.Image, error) {
	sizes := RankingGraphParams{
		minX:      1200,
		minY:      800,
		marginX:   200,
		marginY:   100,
		radius:    30,
		lineWidth: 5,
		gapX:      80,
		gapY:      30,
	}

	if len(rankings) == 0 {
		return nil, errors.New("rankings is empty")
	}

	rankingByName := lineUpRankingByName(rankings, count)

	Xcount := len(rankingByName[0].History)
	x := max(float64(sizes.minX), 2*sizes.marginX+float64(Xcount)*(sizes.radius+sizes.gapX)-sizes.gapX)
	y := max(float64(sizes.minY), calcYPos(count, sizes)+sizes.marginY-sizes.gapY)
	dc := gg.NewContext(int(x), int(y))
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawRectangle(0, 0, x, y)
	dc.Fill()

	dc.SetLineWidth(sizes.lineWidth)

	colors := []string{
		"#FF0000",
		"#00FF00",
		"#0000FF",
	}

	for i, r := range rankingByName {
		for j, h := range r.History {
			if j == 0 {
				dc.SetRGB(1, 1, 1)
				dc.DrawStringAnchored(r.Name, x-sizes.marginX+20, calcYPos(h.Rank, sizes)+sizes.radius, 0, 0.5)
			}
			if h.Rank != -1 {
				dc.SetHexColor(colors[i%len(colors)])
				dc.DrawCircle(calcXPos(j, sizes, x)-sizes.radius, calcYPos(h.Rank, sizes)+sizes.radius, sizes.radius)
				dc.Stroke()
			}
		}
	}

	image := dc.Image()

	return &image, nil
}

func calcXPos(count int, sizes RankingGraphParams, xsize float64) float64 {
	return float64(xsize - (sizes.marginX + float64(count)*(sizes.radius*2+sizes.gapX)))
}

func calcYPos(count int, sizes RankingGraphParams) float64 {
	return float64(sizes.marginY + float64(count)*(sizes.radius*2+sizes.gapY))
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

		for k, r := range rankingByName {
			rank := slices.IndexFunc(ranking.Ranking, func(n RankingItem) bool {
				return n.Name == r.Name
			})
			rankingByName[k].History = append(r.History, rankingItem{
				Rank:  rank,
				Value: ranking.Ranking[rank].Value,
			})
		}
	}
	return rankingByName
}
