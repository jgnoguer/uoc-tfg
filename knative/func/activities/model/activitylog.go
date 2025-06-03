package model

import "time"

type ActivityLog struct {
	Id          string         `json:"id"`
	Status      ActivityStatus `json:"status"`
	UpdatedAt   time.Time      `json:"updateTime" db:"update_time"`
	UpdateBy    string         `json:"updateBy" db:"updated_by"`
	Description string         `json:"description" db:"description"`
}
