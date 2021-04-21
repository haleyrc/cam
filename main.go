package main

import (
	"log"
	"os"
	"strconv"

	"github.com/haleyrc/cam/camera"
	"github.com/haleyrc/cam/filter"
)

func main() {
	deviceID, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	desaturator := filter.Desaturator{}
	pixelator := filter.NewPixelateFilter(16)
	cam, err := camera.NewViewer(deviceID, &desaturator, &pixelator)
	if err != nil {
		log.Println(err)
		cam.Close()
	}

	if err := cam.Run(); err != nil {
		log.Println(err)
		cam.Close()
	}
}
