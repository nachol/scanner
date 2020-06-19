package main

import (
	"log"

	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/nachol/scanner/handlers"
	"github.com/nachol/scanner/model"
)

func main() {
	// GIN SETUP
	port := 8000
	r := gin.Default()
	r.LoadHTMLGlob("./assets/html/*")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")
	r.StaticFile("/js/ansi_up.js", "./assets/js/ansi_up.js")
	r.StaticFile("/js/ansi_up.js.map", "./assets/js/ansi_up.js.map")
	r.StaticFile("/js/termynal.js", "./assets/js/termynal.js")

	c, err := model.NewConnection("mongodb://localhost:27017")

	if err != nil {
		log.Fatalln("Error connecting to DB")
	}

	collectionProgram := c.GetCollection()
	handlers.CollectionProgram = collectionProgram //----> Migrandolo a Model
	model.CollectionProgram = collectionProgram

	/*
		Index Showing Programs
	*/
	r.GET("/index", handlers.Index)

	/*
		Show the form to create a program
	*/
	r.GET("/new", handlers.NewProgram)

	/*
		Creates a new Program
	*/
	r.POST("/create-program", handlers.CreateProgram)

	/*
		Deletes an existing program (Soft Delete)
	*/
	r.GET("/delete/:id", handlers.DeleteProgram)

	/*
		Shows the info of a program
	*/
	r.GET("/view/:id", handlers.ViewProgram)

	r.GET("/api/fetchimg", handlers.FetchH1Image)

	r.POST("/runScan", handlers.RunScan)

	r.GET("/terminal/:id/:scanName", handlers.ViewTerminal)

	log.Printf("Listening on port %d", port)
	r.Run(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:8080
}
