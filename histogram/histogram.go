package histogram

import "image"

func New() Histogram {
	return Histogram{
		Red:   make([]int64, 256),
		Green: make([]int64, 256),
		Blue:  make([]int64, 256),
	}
}

type Histogram struct {
	Red   []int64
	Green []int64
	Blue  []int64
}

func (h Histogram) Average() (float64, float64, float64) {
	return avg(h.Red), avg(h.Green), avg(h.Blue)
}

func (h *Histogram) FromConstrained(img image.Image, left, top, right, bottom int) {
	h.Red = make([]int64, 256)
	h.Blue = make([]int64, 256)
	h.Green = make([]int64, 256)

	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			h.Red[r>>8]++
			h.Blue[b>>8]++
			h.Green[g>>8]++
		}
	}
}

func (h *Histogram) From(img image.Image) {
	bounds := img.Bounds()

	h.Red = make([]int64, 256)
	h.Blue = make([]int64, 256)
	h.Green = make([]int64, 256)

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := img.At(x, y).RGBA()

			h.Red[r>>8]++
			h.Blue[b>>8]++
			h.Green[g>>8]++
		}
	}
}

func avg(vals []int64) float64 {
	total := float64(sum(vals))
	avg := 0.0
	for i := range vals {
		avg += (float64(vals[i]) / total) * float64(i)
	}
	return avg
}

func sum(vals []int64) int64 {
	total := int64(0)
	for _, val := range vals {
		total += val
	}
	return total
}
