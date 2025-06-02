package main

import (
	"context"
	"fmt"
	"github.com/blackjack/webcam"
	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/kelseyhightower/envconfig"
	"log"
	"os"
)

type Receiver struct {
	client cloudevents.Client

	// If the K_SINK environment variable is set, then events are sent there,
	// otherwise we simply reply to the inbound request.
	Target string `envconfig:"TELEGRAM_BROKER"`
}

// Request is the structure of the event we expect to receive.
type Request struct {
	Name string `json:"name"`
}

// Response is the structure of the event we send in response to requests.
type Response struct {
	Message string `json:"message,omitempty"`
}

// handle shared the logic for producing the Response event from the Request.
func handle(req Request) Response {
	return Response{Message: fmt.Sprintf("Sensor was triggered, %s", req.Name)}
}

func (recv *Receiver) ReceiveAndSend(ctx context.Context, event cloudevents.Event) cloudevents.Result {
	// Here is where your code to process the event will go.
	// In this example we will log the event msg
	log.Printf("Event received. \n%s\n", event)
	data := &SensorTriggered{}
	if err := event.DataAs(data); err != nil {
		log.Printf("Error while extracting cloudevent Data: %s\n", err.Error())
		return cloudevents.NewHTTPResult(400, "failed to convert data: %s", err)
	}
	log.Printf("Received event %q", data.Msg)

	r := cloudevents.NewEvent(cloudevents.VersionV1)
	r.SetType("dev.jgnoguer.knative.uoc.sensortriggered")
	r.SetSource("dev.jgnoguer.knative.uoc/sensorresponse")

	if err := r.SetData(cloudevents.TextPlain, fmt.Sprintf("The sensor %s reported temperature %f. Msg: %s", data.SensorID, data.Temperature, data.Msg)); err != nil {
		return cloudevents.NewHTTPResult(500, "failed to set response data: %s", err)
	}

	ctx = cloudevents.ContextWithTarget(ctx, recv.Target)

	//go takePhoto()

	return recv.client.Send(ctx, r)

}

func takePhoto() {
	videoDevice := os.Getenv("VIDEO_DEVICE")
	// ...
	cam, err := webcam.Open("/dev/" + videoDevice) // Open webcam
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
	log.Print("sensorcontrol sample started.")
	client, err := cloudevents.NewDefaultClient()
	if err != nil {
		log.Fatal(err.Error())
	}

	r := Receiver{client: client}
	if err := envconfig.Process("", &r); err != nil {
		log.Fatal(err.Error())
	}

	// Depending on whether targeting data has been supplied,
	// we will either reply with our response or send it on to
	// an event sink.
	var receiver interface{} // the SDK reflects on the signature.
	receiver = r.ReceiveAndSend

	if err := client.StartReceiver(context.Background(), receiver); err != nil {
		log.Fatal(err)
	}
}
