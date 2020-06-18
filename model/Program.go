package model

import (
	"context"
	"log"

	"github.com/nachol/scanner/scan"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	// scan "scanner/scan"
	"time"
)

var CollectionProgram *mongo.Collection

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

type Programs []*Program

type Domain struct {
	Name string `bson:"name"`
}

func (p *Program) getTargets() []Target {
	return p.Targets
}

func GetPrograms() Programs {
	var progs = []*Program{}

	cur, err := CollectionProgram.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatal(err)
	}
	cur.All(context.TODO(), &progs)

	return progs
}

func CreateProgram(p *Program) (*mongo.InsertOneResult, error) {
	res, err := CollectionProgram.InsertOne(context.TODO(), p)
	if err != nil {
		return nil, err
	}
	return res, nil
}
