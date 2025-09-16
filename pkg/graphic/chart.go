package graphic

import (
	"cmp"
	"errors"
	"fmt"
	"image"
	"slices"

	"github.com/fogleman/gg"
	"github.com/oka4shi/kusamochi/pkg/graphic/font/bizudpgothicregular"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

type BarChartParams struct {
	width      float64
	height     float64
	marginX    float64
	marginY    float64
	lineWidth  float64
	gap        float64
	labelWidth float64
}

func DrawBarChart(data Ranking, title string, count int) (*image.Image, error) {
	sizes := BarChartParams{
		width:      1200,
		height:     800,
		marginX:    50,
		marginY:    100,
		lineWidth:  2,
		gap:        5,
		labelWidth: 150,
	}

	if len(data.Ranking) == 0 {
		return nil, errors.New("ranking is empty")
	}

	highestValue := slices.MaxFunc(data.Ranking, func(a, b RankingItem) int {
		return cmp.Compare(a.Value, b.Value)
	}).Value
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

	// Draw title
	dc.SetRGB(1, 1, 1)
	dc.DrawStringAnchored(title,
		sizes.width/2,
		sizes.marginY/2,
		0.5,
		0.5)

	marginLeft := sizes.marginX + sizes.labelWidth + sizes.gap*2
	areaX := sizes.width - marginLeft - sizes.marginX
	areaY := sizes.height - 2*sizes.marginY

	scaleGap := areaX / (float64(scaleLimit)/float64(scaleInterval) + 0.5)

	// Draw vertical scales
	for i := 0; i <= scaleLimit; i += scaleInterval {
		dc.SetRGB(0.25, 0.25, 0.25)
		dc.SetLineWidth(sizes.lineWidth)
		x := marginLeft + scaleGap*float64(i/scaleInterval)
		dc.DrawLine(x, sizes.marginY, x, sizes.height-sizes.marginY)
		dc.Stroke()

		// Draw scale text
		dc.SetRGB(1, 1, 1)
		dc.DrawStringAnchored(fmt.Sprintf("%d", i),
			x,
			sizes.height-sizes.marginY+sizes.gap,
			0.5,
			1,
		)
	}

	// Draw bars
	barHeight := areaY / float64(len(data.Ranking))
	paddingY := barHeight * 0.1

	for i, r := range data.Ranking {
		x := marginLeft
		y := sizes.marginY + (barHeight)*float64(i) + paddingY
		actualBarHeight := barHeight - paddingY*2
		barWidth := float64(r.Value)/float64(scaleInterval)*scaleGap + sizes.lineWidth
		color := colors[i%len(colors)]

		// Draw bar
		dc.SetRGB(color[0], color[1], color[2])
		dc.DrawRectangle(x, y, barWidth, actualBarHeight)
		dc.Fill()

		dc.SetRGB(1, 1, 1)
		// Draw name
		dc.DrawStringAnchored(r.Name,
			x-sizes.gap,
			y+actualBarHeight/2,
			1,
			0.5,
		)
		// Draw value
		dc.DrawStringAnchored(fmt.Sprintf("%d", r.Value),
			x+barWidth+sizes.gap,
			y+actualBarHeight/2,
			0,
			0.5,
		)
	}

	image := dc.Image()
	return &image, nil
}
