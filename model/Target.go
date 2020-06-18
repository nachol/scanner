package model

import (
	"time"
)

/*
Target : List of targets of a scope
*/
type Target struct {
	ID        int
	Name      string
	ProgramID uint
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}
