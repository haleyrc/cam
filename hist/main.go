package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
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

	frame := gocv.NewMat()
	defer frame.Close()

	hist := NewHistogram()

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

		hist.From(img)

		avgR, avgG, avgB := hist.Average()
		rHeight := int(100 * avgR / 256)
		gHeight := int(100 * avgG / 256)
		bHeight := int(100 * avgB / 256)

		rgbaImg := img.(*image.RGBA)
		boxAt(rgbaImg, 10, 10+(100-rHeight), 30, rHeight, color.RGBA{R: 255.0})
		boxAt(rgbaImg, 40, 10+(100-gHeight), 30, gHeight, color.RGBA{G: 255.0})
		boxAt(rgbaImg, 70, 10+(100-bHeight), 30, bHeight, color.RGBA{B: 255.0})

		newMat, err := gocv.ImageToMatRGB(rgbaImg)
		if err != nil {
			fmt.Printf("cannot convert img to mat: %v\n", err)
			return
		}
		window.IMShow(newMat)
		if window.WaitKey(1) >= 0 {
			break
		}
	}
}

func boxAt(img *image.RGBA, left, top int, w, h int, c color.RGBA) {
	c.A = 1.0
	for x := left; x < left+w; x++ {
		for y := top; y < top+h; y++ {
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
