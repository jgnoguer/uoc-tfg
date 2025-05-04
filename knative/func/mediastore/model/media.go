package model

import "time"

type Media struct {
	Id          string    `json:"id"`
	Name        string    `json:"name"`
	ContentType string    `json:"contentType" db:"contenttype"`
	Location    string    `json:"-"`
	CreatedAt   time.Time `json:"uploadTime"`
	Status      int16     `json:"status"`
	Size        int64     `json:"size"`
}
