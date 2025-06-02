package main

import (
	"context"
	"fmt"
	"github.com/blackjack/webcam"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/google/uuid"
	"log"
	"net/http"
	"os"
)

func receive(ctx context.Context, event cloudevents.Event) (*cloudevents.Event, cloudevents.Result) {
	// Here is where your code to process the event will go.
	// In this example we will log the event msg
	log.Printf("Event received. \n%s\n", event)
	data := &SensorTriggered{}
	if err := event.DataAs(data); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return nil, cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Image added from received event %q", data.Msg)

	// Respond with another event (optional)
	// This is optional and is intended to show how to respond back with another event after processing.
	// The response will go back into the knative eventing system just like any other event
	newEvent := cloudevents.NewEvent()
	// Setting the ID here is not necessary. When using NewDefaultClient the ID is set
	// automatically. We set the ID anyway so it appears in the log.
	newEvent.SetID(uuid.New().String())
	newEvent.SetSource("knative/eventing/samples/hello-world")
	newEvent.SetType("dev.jgnoguer.knative.uoc.hifromknative")
	if err := newEvent.SetData(cloudevents.ApplicationJSON, HiFromKnative{Msg: "Hi from helloworld-go app!"}); err != nil {
		return nil, cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}
	log.Printf("Responding with event\n%s\n", newEvent)
	return &newEvent, nil
}

func addAgent(w http.ResponseWriter, r *http.Request) {
	// ...
	cam, err := webcam.Open("/dev/video0") // Open webcam
	if err != nil {
		panic(err.Error())
	}
	defer cam.Close()
	// ...
	// Setup webcam image format and frame size here (see examples or documentation)
	// ...
	err = cam.StartStreaming()
	if err != nil {
		panic(err.Error())
	}

	timeout := 300
	err = cam.WaitForFrame(uint32(timeout))

	switch err.(type) {
	case nil:
	case *webcam.Timeout:
		fmt.Fprint(os.Stderr, err.Error())
	default:
		panic(err.Error())
	}

	frame, err := cam.ReadFrame()
	if len(frame) != 0 {
		// Process frame
		fmt.Println("Have a frame")
		f, err := os.Create("havephoto.jpeg")
		check(err)
		n2, err := f.Write(frame)
		check(err)
		fmt.Printf("wrote %d bytes\n", n2)

	} else if err != nil {
		panic(err.Error())
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
	log.Print("Snapshotter sample started.")
	c, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatalf("failed to create client, %v", err)
	}
	log.Fatal(c.StartReceiver(context.Background(), receive))
}
