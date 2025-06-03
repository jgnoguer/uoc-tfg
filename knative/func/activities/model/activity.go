package model

import "time"

type Activity struct {
	Id          string         `json:"id"`
	ShortCode   string         `json:"shortcode" db:"shortcode"`
	Description string         `json:"description" db:"description"`
	Type        string         `json:"type"`
	Status      ActivityStatus `json:"status"`
	CreatedAt   time.Time      `json:"creationTime" db:"created_at"`
	UpdatedAt   time.Time      `json:"updateTime" db:"updated_at"`
}

type ActivityStatusUpdate struct {
	Id          string         `json:"id"`
	Status      ActivityStatus `json:"status"`
	UpdatedTime time.Time      `json:"updateTime" db:"updated_at"`
	Issuer      string         `json:"issuer" db:"updated_by"`
	Description string         `json:"description"`
}

type ActivityStatus int16

const (
	Scheduled ActivityStatus = iota
	Started
	Finished
	Cancelled
)

func (s ActivityStatus) String() string {
	switch s {
	case Scheduled:
		return "Scheduled"
	case Started:
		return "Started"
	case Finished:
		return "Finished"
	case Cancelled:
		return "Cancelled"
	}
	return "unknown"
}
