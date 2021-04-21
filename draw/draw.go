package draw

import (
	"image"
	"image/color"
)

func BoxAt(img *image.RGBA, left, top, right, bottom int, c color.RGBA) {
	c.A = 1.0
	for x := left; x < right; x++ {
		for y := top; y < bottom; y++ {
			img.SetRGBA(x, y, c)
		}
	}
}
