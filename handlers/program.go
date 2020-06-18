package handlers

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/nachol/scanner/model"
	"github.com/nachol/scanner/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var CollectionProgram *mongo.Collection

func Index(c *gin.Context) {
	results := model.GetPrograms()
	c.HTML(http.StatusOK, "index.tmpl", gin.H{
		"title":    "Index",
		"programs": results,
	})
}

func NewProgram(c *gin.Context) {
	c.HTML(http.StatusOK, "new.tmpl", gin.H{
		"title": "New Work",
	})

}

func CreateProgram(c *gin.Context) {
	name := c.PostForm("name")
	scope := c.PostForm("scope")
	logo := c.PostForm("logo")

	threads, _ := strconv.Atoi(c.PostForm("threads"))
	url := c.PostForm("url")

	targets := strings.Split(scope, "\n")
	targets = utils.DeleteEmpty(targets)
	w := model.Program{Name: name, Threads: threads, URL: url, Logo: logo}
	for _, target := range targets {
		tmp := model.Target{Name: string(target)}
		w.Targets = append(w.Targets, tmp)

	}
	w.CreatedAt = time.Now()

	insertResult, err := model.CreateProgram(&w)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Success: %s\n", insertResult)

	c.Redirect(302, "/index")
}

func DeleteProgram(c *gin.Context) {

	id := c.Param("id")
	filter := bson.D{{"name", id}}
	_, err := CollectionProgram.DeleteOne(context.TODO(), filter)
	if err != nil {
		log.Fatal(err)
	}

	c.Redirect(302, "/index")
}

func ViewProgram(c *gin.Context) {
	id := c.Param("id")
	var program model.Program

	filter := bson.D{{"name", id}}
	err := CollectionProgram.FindOne(context.TODO(), filter).Decode(&program)
	if err != nil {
		log.Fatal(err)
	}

	c.HTML(http.StatusOK, "view.tmpl", gin.H{
		"program": program,
	})
}
