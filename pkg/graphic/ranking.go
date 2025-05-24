package graphic

import (
	"errors"
	"image"
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
	minX          float64
	minY          float64
	marginX       float64
	marginY       float64
	radius        float64
	lineWidth     float64
	boldLineWidth float64
	gapX          float64
	gapY          float64
}

func DrawRanking(rankings []Ranking, count int) (*image.Image, error) {
	sizes := RankingGraphParams{
		minX:          900,
		minY:          600,
		marginX:       200,
		marginY:       100,
		radius:        30,
		lineWidth:     2,
		boldLineWidth: 5,
		gapX:          300,
		gapY:          30,
	}

	if len(rankings) == 0 {
		return nil, errors.New("rankings is empty")
	}

	rankingByName := lineUpRankingByName(rankings, count)

	Xcount := len(rankingByName[0].History)
	width := max(float64(sizes.minX), 2*sizes.marginX+float64(Xcount)*(sizes.radius+sizes.gapX)-sizes.gapX)
	height := max(float64(sizes.minY), calcYPos(count, sizes)+sizes.marginY-sizes.gapY)
	dc := gg.NewContext(int(width), int(height))

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

	dc.SetRGB(0.1, 0.1, 0.1)
	dc.DrawRectangle(0, 0, width, height)
	dc.Fill()

	for i, r := range rankings {
		dc.SetRGB(0.25, 0.25, 0.25)
		dc.SetLineWidth(sizes.lineWidth)
		x := calcXPos(i, sizes, width) - sizes.radius
		dc.DrawLine(x, sizes.marginY, x, height-sizes.marginY)
		dc.Stroke()

		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(r.Time, x, height-sizes.marginY, 0.5, 1)
	}

	for i := range len(rankings[0].Ranking) {
		dc.SetRGB(0.25, 0.25, 0.25)
		dc.SetLineWidth(sizes.lineWidth)
		y := calcYPos(i, sizes) + sizes.radius
		dc.DrawLine(0, y, width-sizes.marginX, y)
		dc.Stroke()

	}

	for i, r := range rankingByName {
		for j, h := range r.History {
			if j == 0 {
				dc.SetRGB(1, 1, 1)
				dc.DrawStringAnchored(r.Name, width-sizes.marginX+20, calcYPos(h.Rank, sizes)+sizes.radius, 0, 0.5)
			}
			if h.Rank != -1 && h.Rank < count {
				x := calcXPos(j, sizes, width) - sizes.radius
				y := calcYPos(h.Rank, sizes) + sizes.radius

				// Draw a value
				dc.SetRGB(1, 1, 1)
				dc.DrawStringAnchored(h.Value, x, y, 0.5, 0.5)

				// Draw a circle
				c := colors[i%len(colors)]
				dc.SetRGB(c[0], c[1], c[2])
				dc.SetLineWidth(sizes.boldLineWidth)
				dc.DrawCircle(x, y, sizes.radius)
				dc.Stroke()

				if j != 0 {
					prevRank := r.History[j-1].Rank
					if prevRank != -1 && prevRank < count {
						drawTranverseLine(dc, sizes, width, h.Rank, prevRank, j)
					}
				}
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

func drawTranverseLine(dc *gg.Context, sizes RankingGraphParams, canvasWidth float64, rank int, prevRank int, count int) {
	if rank == prevRank {
		y := calcYPos(rank, sizes) + sizes.radius
		dc.DrawLine(calcXPos(count, sizes, canvasWidth), y,
			calcXPos(count-1, sizes, canvasWidth)-sizes.radius*2, y)
		dc.Stroke()
	} else {
		startX := calcXPos(count-1, sizes, canvasWidth) - sizes.radius*2
		startY := calcYPos(prevRank, sizes) + sizes.radius
		endX := calcXPos(count, sizes, canvasWidth)
		endY := calcYPos(rank, sizes) + sizes.radius

		dc.MoveTo(startX, startY)
		dc.CubicTo((startX+endX)/2, startY, (startX+endX)/2, endY, endX, endY)
		dc.Stroke()
	}

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
