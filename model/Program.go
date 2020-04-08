package model

import (
	"time"
)

/*
Program : Program Of H1
*/
type Program struct {
	ID        int
	Name      string
	Targets   []Target `gorm:"foreignkey:ProgramId"`
	Threads   int
	URL       string
	Logo      string
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

func (p *Program) getTargets() []Target {
	return p.Targets
}
