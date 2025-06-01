package main

// HelloWorld defines the Data of CloudEvent with type=dev.jgnoguer.knative.uoc.imageadded
type ImageAdded struct {
	// Msg holds the message from the event
	MediaId string `json:"mediaId,omitempty"`
}

// HiFromKnative defines the Data of CloudEvent with type=dev.jgnoguer.knative.uoc.hifromknative
type HiFromKnative struct {
	// Msg holds the message from the event
	Msg string `json:"msg,omitempty"`
}
