package function

import (
	"fmt"
	"github.com/blackjack/webcam"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

// Handle an HTTP Request.
func Handle(w http.ResponseWriter, r *http.Request) {

	godotenv.Load()

	w.Header().Set("Content-Type", "application/json")

	switch r.Method {
	case "POST":
		addAgent(w, r)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}

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
