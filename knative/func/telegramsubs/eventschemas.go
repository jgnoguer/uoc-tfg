package main

// defines the Data of CloudEvent with type=dev.jgnoguer.knative.uoc.imageadded
type TelegramMsg struct {
	// Msg holds the message from the event
	Msg string `json:"msg,omitempty"`
}
