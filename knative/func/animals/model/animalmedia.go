package model

import "time"

type AnimalMedia struct {
	Id          string    `json:"id"`
	MediaId     string    `json:"mediaId" db:"media_id"`
	Description string    `json:"description"`
	Type        MediaType `json:"type"`
	CreatedAt   time.Time `json:"creationTime" db:"created_at"`
}
type MediaType int16

const (
	// since iota starts with 0, the first value
	// defined here will be the default
	Unknown MediaType = iota
	Image
	Video
	Sound
)

func (s MediaType) String() string {
	switch s {
	case Image:
		return "Image"
	case Video:
		return "Video"
	case Sound:
		return "Sound"
	}
	return "unknown"
}
