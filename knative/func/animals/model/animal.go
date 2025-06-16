package model

import "time"

type Animal struct {
	Id          string       `json:"id"`
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Breed       string       `json:"breed"`
	Type        AnimalType   `json:"type"`
	Status      AnimalStatus `json:"status"`
	CreatedAt   time.Time    `json:"creationTime"`
	UpdatedAt   time.Time    `json:"updateTime"`
}
type AnimalType int16

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	Undefined AnimalType = iota
	Dog
	Cat
)

func (s AnimalType) String() string {
	switch s {
	case Dog:
		return "Dog"
	case Cat:
		return "Cat"
	}
	return "unknown"
}

type AnimalStatus int16

const (
	Admitted AnimalStatus = iota
	Hosted
	Adopted
)

func (s AnimalStatus) String() string {
	switch s {
	case Admitted:
		return "Admitted"
	case Hosted:
		return "Hosted"
	case Adopted:
		return "Adopted"
	}
	return "unknown"
}
