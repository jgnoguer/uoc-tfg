package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	cloudevents "github.com/cloudevents/sdk-go/v2"
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

	gatherImage(data.MediaId)

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

func gatherImage(imageId string) {
	endpoint := os.Getenv("MEDIASTORAGE_ENDPOINT")
	storagefolder := os.Getenv("STORAGE_FOLDER")
	resp, err := http.Get(endpoint + "/" + imageId)
	if err != nil {
		log.Fatalln(err)
	}
	newFile, err := os.Create(filepath.Join(storagefolder, imageId))
	_, errCopy := io.Copy(newFile, resp.Body)
	if errCopy != nil {
		log.Fatalln(err)
	}
	log.Print("Got the image.")

}
