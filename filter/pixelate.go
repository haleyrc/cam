package filter

import (
	"image"
	"image/color"

	"github.com/haleyrc/cam/draw"
	"github.com/haleyrc/cam/histogram"
	"github.com/haleyrc/cam/math"
)

func NewPixelateFilter(blockSize int) Pixelate {
	return Pixelate{
		hist:      histogram.New(),
		blockSize: blockSize,
	}
}

type Pixelate struct {
	hist      histogram.Histogram
	blockSize int
}

func (p *Pixelate) OnKey(key int) {
	switch key {
	case 61: // =
		p.SetBlockSize(p.BlockSize() * 2)
	case 45: // -
		p.SetBlockSize(p.BlockSize() / 2)
	}
}

func (p *Pixelate) SetBlockSize(s int) {
	p.blockSize = math.ClampLower(s, 4)
}

func (p Pixelate) BlockSize() int {
	return p.blockSize
}

func (p Pixelate) Filter(img *image.RGBA) {
	bounds := img.Bounds()

	for x := 0; x < bounds.Max.X; x += p.blockSize {
		for y := 0; y < bounds.Max.Y; y += p.blockSize {
			right := math.ClampUpper(x+p.blockSize, bounds.Max.X)
			bottom := math.ClampUpper(y+p.blockSize, bounds.Max.Y)
			p.hist.FromConstrained(img, x, y, right, bottom)
			avgR, avgG, avgB := p.hist.Average()
			draw.BoxAt(img, x, y, right, bottom, color.RGBA{
				R: uint8(avgR),
				G: uint8(avgG),
				B: uint8(avgB),
				A: 1,
			})
		}
	}
}
