package model

import (
	"github.com/nachol/scanner/scan"
	// scan "scanner/scan"
	"time"
)

/*
Program : Program Of H1
*/
type Program struct {
	Name      string      `bson:"name"`
	Targets   []Target    `bson:"targets"`
	Threads   int         `bson:"threads"`
	URL       string      `bson:"url"`
	Logo      string      `bson:"logo"`
	CreatedAt time.Time   `bson:"created"`
	Domains   []Domain    `bson:"domains"`
	Scans     []scan.Scan `bson:"Scans"`
}

type Domain struct {
	Name string `bson:"name"`
}

func (p *Program) getTargets() []Target {
	return p.Targets
}
