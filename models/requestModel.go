package models

import (
	"time"

	"gorm.io/gorm"
)

type Request struct {
	gorm.Model
	UserID        uint
	Address       string
	WorkType      string
	WorkScale     string
	Time          time.Time
	Status        string
	StatusMessage string
	Workers       string
}
