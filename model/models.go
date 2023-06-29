package model

import "time"

type User struct {
	Id       int    `json:"-" db:"id"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type Audit struct {
	EventDate time.Time `json:"event_date" db:"event_date"`
	Event     string    `json:"event" db:"event"`
}
