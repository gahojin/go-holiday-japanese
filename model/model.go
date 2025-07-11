package model

import (
	"time"
)

type Name struct {
	Ja string `json:"ja"`
	En string `json:"en"`
}

type Holiday struct {
	Name
	Date time.Time `json:"date"`
}
