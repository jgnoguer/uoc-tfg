package main

// defines the Data of CloudEvent with type=dev.jgnoguer.knative.uoc.imageadded
type SensorTriggered struct {
	// Msg holds the message from the event
	Msg         string  `json:"msg,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	SensorID    string  `json:"sensorID,omitempty"`
}
