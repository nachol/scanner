package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nachol/scanner/model"
	"github.com/nachol/scanner/scan"
	"go.mongodb.org/mongo-driver/bson"
)

// var CollectionProgram *mongo.Collection

func RunScan(c *gin.Context) {
	id := c.PostForm("programID")
	var program model.Program
	filter := bson.D{{"name", id}}
	err := CollectionProgram.FindOne(context.TODO(), filter).Decode(&program)
	if err != nil {
		log.Fatal(err)
	}

	var scan scan.Scan
	c.Bind(&scan)
	options := c.PostFormMap("options")
	scan.Options = options
	scan.Threads = program.Threads
	scan.LoadModes()
	_, err = scan.Run()
	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was a problem running the scan",
		})

	}
	program.Scans = append(program.Scans, scan)
	if err != nil {
		log.Panicln(err)
	}
	update := bson.D{
		{"$set", bson.D{
			{"Scans", program.Scans},
		}},
	}
	updateResult, err := CollectionProgram.UpdateOne(context.TODO(), filter, update)

	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"message": "Work Done",
		"result":  updateResult,
	})
}

func ViewTerminal(c *gin.Context) {
	id := c.Param("id")
	// scanName := c.Param("scanName")

	var program model.Program

	// filter := bson.M{
	// 	"name": bson.M{"$eq": id},
	// 	"Scans": bson.M{
	// 		"$elemMatch": bson.M{
	// 			"name": bson.M{"$eq": scanName},
	// 		},
	// 	},
	// }
	filter := bson.D{{"name", id}}
	// err = collectionProgram.FindOne(context.TODO(), filter).Decode(&program)
	err := CollectionProgram.FindOne(context.TODO(), filter).Decode(&program)

	if err != nil {
		log.Fatal(err)
	}
	c.HTML(http.StatusOK, "terminal.tmpl", gin.H{
		"raw": program.Scans[0].Raw,
	})
}
