package model

import "time"

type Agent struct {
	Id        string      `json:"id"`
	FirstName string      `json:"firstname" db:"firstname"`
	LastName  string      `json:"lastname" db:"lastname"`
	Type      AgentType   `json:"type"`
	Status    AgentStatus `json:"status"`
	CreatedAt time.Time   `json:"creationTime"`
}
type AgentType int16

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	Undefined AgentType = iota
	Vet
	Staff
	Volunteer
	Trainer
)

func (s AgentType) String() string {
	switch s {
	case Vet:
		return "Vet"
	case Staff:
		return "Staff"
	case Volunteer:
		return "Volunteer"
	case Trainer:
		return "Dog trainer"
	}
	return "unknown"
}

type AgentStatus int16

const (
	Inactive AgentStatus = iota
	Active
)

func (s AgentStatus) String() string {
	switch s {
	case Active:
		return "Active"
	case Inactive:
		return "Inactive"
	}
	return "unknown"
}
