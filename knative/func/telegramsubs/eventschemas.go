package main

// defines the Data of CloudEvent with type=dev.jgnoguer.knative.uoc.imageadded
type TelegramMsg struct {
	MsgId string `json:"message_id,omitempty"`

	Date         int `json:"date"`
	TelegramFrom struct {
		// Msg holds the message from the event
		Id        int    `json:"id"`
		Username  string `json:"username"`
		Bot       bool   `json:"bot"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		IsBot     bool   `json:"is_bot"`
	} `json:"from"`
	Text string `json:"text"`
	Chat struct {
		Id                          string `json:"id"`
		Title                       string `json:"title"`
		Type                        string `json:"type"`
		AllMembersAreAdministrators bool   `json:"all_members_are_administrators"`
	} `json:"chat"`
}
