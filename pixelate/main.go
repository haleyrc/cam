package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"

	"gocv.io/x/gocv"
)

func main() {
	deviceID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		log.Fatalf("error opening capture device: %d: %v\n", deviceID, err)
	}
	defer webcam.Close()

	window := gocv.NewWindow("Test")
	defer window.Close()
	window.SetWindowProperty(gocv.WindowPropertyFullscreen, gocv.WindowFullscreen)

	frame := gocv.NewMat()
	defer frame.Close()

	hist := NewHistogram()
	step := 16
	pixelate := true

	for {
		if ok := webcam.Read(&frame); !ok {
			fmt.Printf("cannot read device %d: %v\n", deviceID, err)
			return
		}
		if frame.Empty() {
			fmt.Println("empty")
			continue
		}

		img, err := frame.ToImage()
		if err != nil {
			fmt.Printf("cannot convert to image: %v\n", err)
			return
		}
		rgbaImg := img.(*image.RGBA)
		bounds := img.Bounds()

		if pixelate {
			for x := 0; x < bounds.Max.X; x += step {
				for y := 0; y < bounds.Max.Y; y += step {
					right := constrainUpper(x+step, bounds.Max.X)
					bottom := constrainUpper(y+step, bounds.Max.Y)
					hist.FromConstrained(img, x, y, right, bottom)
					avgR, avgG, avgB := hist.Average()
					boxAt(rgbaImg, x, y, right, bottom, color.RGBA{
						R: uint8(avgR),
						G: uint8(avgG),
						B: uint8(avgB),
						A: 1,
					})
				}
			}
		}

		newMat, err := gocv.ImageToMatRGB(rgbaImg)
		if err != nil {
			fmt.Printf("cannot convert img to mat: %v\n", err)
			return
		}
		window.IMShow(newMat)

		key := window.WaitKey(1)
		switch key {
		case 27:
			return
		case 61: // =
			shortSide := math.Min(float64(bounds.Max.X), float64(bounds.Max.Y))
			step = constrainUpper(step*2, int(shortSide))
		case 45: // -
			step = constrainLower(step/2, 4)
		case 112: // p
			pixelate = !pixelate
		default:
			if key != -1 {
				fmt.Println(key)
			}
		}
	}
}

func constrainLower(val, bound int) int {
	if val < bound {
		return bound
	}
	return val
}

func constrainUpper(val, bound int) int {
	if val > bound {
		return bound
	}
	return val
}

func boxAt(img *image.RGBA, left, top, right, bottom int, c color.RGBA) {
	c.A = 1.0
	for x := left; x < right; x++ {
		for y := top; y < bottom; y++ {
			img.SetRGBA(x, y, c)
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

func NewHistogram() Histogram {
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
