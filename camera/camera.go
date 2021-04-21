package camera

import (
	"fmt"
	"image"

	"github.com/pkg/errors"
	"gocv.io/x/gocv"
)

type Filter interface {
	Filter(img *image.RGBA)
	OnKey(key int)
}

func NewViewer(deviceID int, filters ...Filter) (Viewer, error) {
	webcam, err := gocv.VideoCaptureDevice(deviceID)
	if err != nil {
		return Viewer{}, errors.Errorf("error opening capture device: %d: %v\n", deviceID, err)
	}

	return Viewer{
		webcam:         webcam,
		filtersEnabled: true,
		filters:        filters,
		fullscreen:     true,
	}, nil
}

type Viewer struct {
	webcam         *gocv.VideoCapture
	filtersEnabled bool
	filters        []Filter
	fullscreen     bool
}

func (v Viewer) Close() {
	v.webcam.Close()
}

func (v Viewer) Run() error {
	window := gocv.NewWindow("Test")
	defer window.Close()
	window.SetWindowProperty(gocv.WindowPropertyFullscreen, gocv.WindowFullscreen)

	frame := gocv.NewMat()
	defer frame.Close()

	for {
		if ok := v.webcam.Read(&frame); !ok || frame.Empty() {
			fmt.Println("dropped frame: could not read frame data")
			continue
		}

		img, err := frame.ToImage()
		if err != nil {
			fmt.Printf("dropped frame: cannot convert frame to image: %v\n", err)
			continue
		}
		rgbaImg := img.(*image.RGBA)

		if v.filtersEnabled {
			for _, filter := range v.filters {
				filter.Filter(rgbaImg)
			}
		}

		newMat, err := gocv.ImageToMatRGB(rgbaImg)
		if err != nil {
			fmt.Printf("dropped frame: cannot convert image to mat: %v\n", err)
			continue
		}
		window.IMShow(newMat)

		key := window.WaitKey(1)
		switch key {
		case 27: // Esc
			return nil
		case 102: // f
			v.filtersEnabled = !v.filtersEnabled
		case 200: // f11
			if v.fullscreen {
				window.SetWindowProperty(gocv.WindowPropertyFullscreen, gocv.WindowNormal)
			} else {
				window.SetWindowProperty(gocv.WindowPropertyFullscreen, gocv.WindowFullscreen)
			}
			v.fullscreen = !v.fullscreen
		default:
			if key != -1 {
				fmt.Println(key)
			}
		}
		for _, filter := range v.filters {
			filter.OnKey(key)
		}
	}
}
