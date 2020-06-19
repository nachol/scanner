package handlers

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nachol/scanner/model"
	"go.mongodb.org/mongo-driver/bson"
)

func RunScan(c *gin.Context) {
	id := c.PostForm("programID")

	program, err := model.GetProgramById(id)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	scan, err := model.CreateBindScan(c, program)

	_, err = scan.Run()
	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was a problem running the scan",
		})
		return
	}
	updateResult, err := program.UpdateScans(scan)

	if err != nil {
		c.JSON(400, gin.H{
			"error": "There was a problem Saving the results",
		})
		return
	}

	c.JSON(200, gin.H{
		"message": "Work Done",
		"result":  updateResult,
	})
	return
}

func ViewTerminal(c *gin.Context) {
	id := c.Param("id")
	// scanName := c.Param("scanName")

	var program model.Program

	filter := bson.D{{"name", id}}
	err := CollectionProgram.FindOne(context.TODO(), filter).Decode(&program)

	if err != nil {
		log.Fatal(err)
	}
	c.HTML(http.StatusOK, "terminal.tmpl", gin.H{
		"raw": program.Scans[0].Raw,
	})
}
