package model

import "time"

type Group struct {
	Id          string      `json:"id"`
	Name        string      `json:"name" db:"name"`
	Description string      `json:"description"`
	Members     []string    `json:"members"`
	Type        GroupType   `json:"type"`
	Status      GroupStatus `json:"status"`
	CreatedAt   time.Time   `json:"creationTime" db:"created_at"`
	UpdatedAt   time.Time   `json:"updateTime" db:"updated_at"`
}
type GroupType int16

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	Undefined GroupType = iota
	ActivityAnimalGroup
)

func (s GroupType) String() string {
	switch s {
	case ActivityAnimalGroup:
		return "ActivityAnimalGroup"
	}
	return "unknown"
}

type GroupStatus int16

const (
	Inactive GroupStatus = iota
	Active
)

func (s GroupStatus) String() string {
	switch s {
	case Active:
		return "Active"
	case Inactive:
		return "Inactive"
	}
	return "unknown"
}
