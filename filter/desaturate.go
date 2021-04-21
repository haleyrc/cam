package filter

import (
	"image"
)

type Channel int

const (
	None Channel = iota
	Red
	Green
	Blue
	Alpha
)

type Desaturator struct {
	only Channel
}

func (d *Desaturator) OnKey(key int) {
	switch key {
	case 114: // r
		if d.only == Red {
			d.only = None
			return
		}
		d.only = Red
	case 103: // g
		if d.only == Green {
			d.only = None
			return
		}
		d.only = Green
	case 98: // b
		if d.only == Blue {
			d.only = None
			return
		}
		d.only = Blue
	}
}

func (d Desaturator) Filter(img *image.RGBA) {
	bounds := img.Bounds()

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			c := img.RGBAAt(x, y)
			switch d.only {
			case Red:
				c.G = 0
				c.B = 0
			case Green:
				c.R = 0
				c.B = 0
			case Blue:
				c.R = 0
				c.G = 0
			}
			img.SetRGBA(x, y, c)
		}
	}
}
