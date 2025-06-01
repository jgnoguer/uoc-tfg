package function

// defines the Data of CloudEvent with type=dev.jgnoguer.knative.uoc.imageadded
type ImageAdded struct {
	// Msg holds the message from the event
	MediaId     string `json:"mediaId,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Size        int64  `json:"size,omitempty"`
	ContentType string `json:"contentType,omitempty"`
}
