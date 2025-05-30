package wattermark

import (
	"fmt"
	"github.com/disintegration/imaging"
	"image"
	"log"
	"math"
	"strconv"
	"strings"
)

func AddWaterMark(bgImg string, watermark string, watermarkSize string) {

	outName := fmt.Sprintf("watermark-new-%s", watermark)

	src := openImage(bgImg)

	markFit := resizeImage(watermark, watermarkSize)

	bgDimensions := src.Bounds().Max
	markDimensions := markFit.Bounds().Max

	bgAspectRatio := math.Round(float64(bgDimensions.X) / float64(bgDimensions.Y))

	xPos, yPos := calcWaterMarkPosition(bgDimensions, markDimensions, bgAspectRatio)

	placeImage(outName, bgImg, watermark, watermarkSize, fmt.Sprintf("%dx%d", xPos, yPos))

	fmt.Printf("Added watermark '%s' to image '%s' with dimensions %s.\n", watermark, bgImg, watermarkSize)
}

func placeImage(outName, bgImg, markImg, markDimensions, locationDimensions string) {

	// Coordinate to super-impose on. e.g. 200x500
	locationX, locationY := parseCoordinates(locationDimensions, "x")

	src := openImage(bgImg)

	// Resize the watermark to fit these dimensions, preserving aspect ratio.
	markFit := resizeImage(markImg, markDimensions)

	// Place the watermark over the background in the location
	dst := imaging.Paste(src, markFit, image.Pt(locationX, locationY))

	err := imaging.Save(dst, outName)

	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	fmt.Printf("Placed image '%s' on '%s'.\n", markImg, bgImg)
}

func resizeImage(image, dimensions string) image.Image {

	width, height := parseCoordinates(dimensions, "x")

	src := openImage(image)

	return imaging.Fit(src, width, height, imaging.Lanczos)
}

func openImage(name string) image.Image {
	src, err := imaging.Open(name)
	if err != nil {
		log.Fatalf("failed to open image %s: %v", name, err)
	}
	return src
}

func parseCoordinates(input, delimiter string) (int, int) {

	arr := strings.Split(input, delimiter)

	// convert a string to an int
	x, err := strconv.Atoi(arr[0])

	if err != nil {
		log.Fatalf("failed to parse x coordinate: %v", err)
	}

	y, err := strconv.Atoi(arr[1])

	if err != nil {
		log.Fatalf("failed to parse y coordinate: %v", err)
	}

	return x, y
}

// Subtracts the dimensions of the watermark and padding based on the background’s aspect ratio
func calcWaterMarkPosition(bgDimensions, markDimensions image.Point, aspectRatio float64) (int, int) {
	bgX := bgDimensions.X
	bgY := bgDimensions.Y
	markX := markDimensions.X
	markY := markDimensions.Y
	padding := 20 * int(aspectRatio)
	return bgX - markX - padding, bgY - markY - padding
}
