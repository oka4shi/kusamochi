package graphic

import (
	"errors"
	"fmt"
	"image"
	"math"
	"slices"

	"github.com/fogleman/gg"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/oka4shi/kusamochi/pkg/graphic/font/bizudpgothicregular"
)

type Ranking struct {
	Time    string
	Ranking []RankingItem
}

type RankingItem struct {
	Name  string
	Value int
}

type rankingItem struct {
	Rank  int
	Value int
}

type rankingByNameType struct {
	Name    string
	History []rankingItem
}

type RankingGraphParams struct {
	width         float64
	height        float64
	marginX       float64
	marginY       float64
	radius        float64
	lineWidth     float64
	boldLineWidth float64
}

func DrawRanking(rankings []Ranking, count int) (*image.Image, error) {
	sizes := RankingGraphParams{
		width:         1200,
		height:        800,
		marginX:       200,
		marginY:       100,
		radius:        30,
		lineWidth:     2,
		boldLineWidth: 5,
	}

	if len(rankings) == 0 {
		return nil, errors.New("rankings is empty")
	}

	rankingByName := lineUpRankingByName(rankings, count)

	highestValue := getHighestValue(rankings)
	scaleInterval := calcScaleInterval(highestValue)
	scaleLimit := calcScaleLimit(highestValue, scaleInterval)

	dc := gg.NewContext(int(sizes.width), int(sizes.height))

	colors := GenerateColors(count)

	f, err := opentype.Parse(bizudpgothicregular.TTF)
	if err != nil {
		return nil, err
	}
	face, err := opentype.NewFace(f, &opentype.FaceOptions{
		Size:    18,
		DPI:     72,
		Hinting: font.HintingNone,
	})
	if err != nil {
		return nil, err
	}
	dc.SetFontFace(face)

	// Draw background
	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawRectangle(0, 0, sizes.width, sizes.height)
	dc.Fill()

	// Draw horizontal scales
	for i := 0; i <= scaleLimit; i += scaleInterval {
		dc.SetRGB(0.25, 0.25, 0.25)
		dc.SetLineWidth(sizes.lineWidth)
		y := calcYPos(i, scaleLimit, sizes)
		dc.DrawLine(sizes.marginX, y, sizes.width-sizes.marginX, y)
		dc.Stroke()

		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(fmt.Sprintf("%d", i),
			sizes.marginX,
			y,
			1,
			0.5,
		)
	}

	// Draw vertical scales
	for i, r := range rankings {
		dc.SetRGB(0.25, 0.25, 0.25)
		dc.SetLineWidth(sizes.lineWidth)
		x := calcXPos(i, len(rankings), sizes)
		dc.DrawLine(x, sizes.marginY, x, sizes.height-sizes.marginY)
		dc.Stroke()

		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(r.Time, x, sizes.height-sizes.marginY, 0.5, 1)
	}

	// Draw scale value text
	for i, r := range rankingByName {
		c := colors[i%len(colors)]
		dc.SetRGB(c[0], c[1], c[2])

		// Skip if all values are less than 5% of scale limit
		shouldHide := true
		for _, data := range r.History {
			if float64(data.Value) >= float64(scaleLimit)*0.05 {
				shouldHide = false
				break
			}
		}
		if shouldHide {
			continue
		}

		//dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(r.Name, sizes.width-sizes.marginX, calcYPos(r.History[0].Value, scaleLimit, sizes), 0, 0.5)

		for j := 1; j < len(r.History); j++ {
			prevData := r.History[j-1]
			data := r.History[j]

			// Draw line
			startX := calcXPos(j-1, len(rankings), sizes)
			endX := calcXPos(j, len(rankings), sizes)

			if data.Value == prevData.Value {
				y := calcYPos(data.Value, scaleLimit, sizes)
				dc.DrawLine(startX, y, endX, y)
				dc.Stroke()
			} else {
				startY := calcYPos(prevData.Value, scaleLimit, sizes)
				endY := calcYPos(data.Value, scaleLimit, sizes)

				fmt.Printf("Draw line: %s (%f,%f) - (%f,%f)\n", r.Name, startX, startY, endX, endY)

				dc.MoveTo(startX, startY)
				dc.CubicTo((startX+endX)/2, startY, (startX+endX)/2, endY, endX, endY)
				dc.Stroke()
			}
		}
	}

	image := dc.Image()

	return &image, nil
}

func calcXPos(count int, total int, sizes RankingGraphParams) float64 {
	gapX := (sizes.width - 2*sizes.marginX) / (float64(total - 1))
	return sizes.width - (sizes.marginX + float64(count)*gapX)
}

func calcYPos(value int, scaleLimit int, sizes RankingGraphParams) float64 {
	return sizes.height - (sizes.marginY + float64(value)/float64(scaleLimit)*(sizes.height-2*sizes.marginY))
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

func getHighestValue(rankings []Ranking) int {
	maxValue := -1
	for _, r := range rankings {
		v := r.Ranking[0].Value
		if v > maxValue {
			maxValue = v
		}
	}
	return maxValue
}

func calcScaleInterval(highestValue int) int {
	if highestValue <= 20 {
		return 1
	}
	if highestValue <= 50 {
		return 5
	}
	return int(math.Ceil(float64(highestValue)/100) * 10)
}

func calcScaleLimit(highestValue, interval int) int {
	return int(math.Ceil(float64(highestValue)/float64(interval)) * float64(interval))
}
