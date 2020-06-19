package model

import (
	"context"
	"log"

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
	Name      string    `bson:"name"`
	Targets   []Target  `bson:"targets"`
	Threads   int       `bson:"threads"`
	URL       string    `bson:"url"`
	Logo      string    `bson:"logo"`
	CreatedAt time.Time `bson:"created"`
	Domains   []Domain  `bson:"domains"`
	Scans     []Scan    `bson:"Scans"`
}

type Programs []*Program

type Domain struct {
	Name string `bson:"name"`
}

func (p *Program) getTargets() []Target {
	return p.Targets
}

func (p *Program) UpdateScans(scan *Scan) (res *mongo.UpdateResult, e error) {
	p.Scans = append(p.Scans, *scan)
	update := bson.D{
		{"$set", bson.D{
			{"Scans", p.Scans},
		}},
	}
	filter := bson.D{{"name", p.Name}}

	res, err := CollectionProgram.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		return nil, err
	}
	return res, nil
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

func GetProgramById(id string) (program *Program, e error) {
	filter := bson.D{{"name", id}}
	err := CollectionProgram.FindOne(context.TODO(), filter).Decode(&program)
	if err != nil {
		return nil, err
	}

	return program, nil
}

func CreateProgram(p *Program) (*mongo.InsertOneResult, error) {
	res, err := CollectionProgram.InsertOne(context.TODO(), p)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func DeleteProgram(id string) error {
	filter := bson.D{{"name", id}}
	_, err := CollectionProgram.DeleteOne(context.TODO(), filter)
	if err != nil {
		return err
	}
	return nil
}
