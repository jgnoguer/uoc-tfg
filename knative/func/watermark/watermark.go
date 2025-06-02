package main

import (
	"context"
	"fmt"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/sunshineplan/imgconv"
	"image"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

func receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	// Here is where your code to process the event will go.
	// In this example we will log the event msg
	log.Printf("Event received. \n%s\n", event)
	data := &ImageAdded{}
	if err := event.DataAs(data); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Image added with %q", data.MediaId)

	gatherImage(context.Background(), data.MediaId)

	return nil, cloudevents.NewHTTPResult(201, "Accepted")
}

func main() {
	log.Print("Watermark sample started.")
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), receive))
}

func gatherImage(ctx context.Context, imageId string) {
	endpoint := os.Getenv("MEDIASTORAGE_ENDPOINT")
	storagefolder := os.Getenv("STORAGE_FOLDER")
	resp, err := http.Get(endpoint + "/" + imageId)
	if err != nil {
		log.Fatalln(err)
	}
	filePath := filepath.Join(storagefolder, imageId)
	newFile, err := os.Create(filePath)
	_, errCopy := io.Copy(newFile, resp.Body)
	if errCopy != nil {
		log.Fatalln(err)
	}
	log.Print("Got the image.")
	wattermarkTest(filePath, filepath.Join(storagefolder, imageId+"_watermark"), "perroagua.png")

}

func wattermarkTest(inputImage string, outputImage string, wattermarkImageFile string) {
	// Open a test image.
	srcImage, err := imgconv.Open(inputImage)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}
	markImage, err := imgconv.Open(wattermarkImageFile)
	if err != nil {
		log.Fatalf("failed to open wattermark image: %v", err)
	}
	markImage400 := imgconv.Resize(markImage, &imgconv.ResizeOption{Width: 200})
	// Resize srcImage to width = 800px preserving the aspect ratio.
	dstImage800 := imgconv.Resize(srcImage, &imgconv.ResizeOption{Width: 1200})
	dstImage := imgconv.Watermark(dstImage800, &imgconv.WatermarkOption{Mark: markImage400, Opacity: 120, Random: false,
		Offset: image.Pt(-520, 270)})

	newFile, err := os.Create(outputImage)

	// Write the resulting image
	if err := imgconv.Write(newFile, dstImage, &imgconv.FormatOption{Format: imgconv.JPEG}); err != nil {
		log.Fatalf("failed to write image: %v", err)
	}
	defer newFile.Close()

}

func getImageDimension(filepath string) (int, int) {

	f, err := os.Open(filepath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	config, _, err := image.DecodeConfig(f)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Width: %d\nHeight: %d\n", config.Width, config.Height)
	return config.Width, config.Height
}
