package main

import (
	model "./model"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"
)

func main() {
	port := 8000
	r := gin.Default()
	r.LoadHTMLGlob("./assets/html/*")
	r.Static("/assets", "./assets")
	r.StaticFile("/favicon.ico", "./assets/favicon.ico")

	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// db.AutoMigrate(&model.Program{})
	// db.AutoMigrate(&model.Target{})

	for _, models := range []interface{}{
		model.Program{},
		model.Target{},
	} {
		if err := db.AutoMigrate(models).Error; err != nil {
			log.Println(err)
		} else {
			log.Println("Auto migrating", reflect.TypeOf(models).Name(), "...")
		}
	}

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	/*
		Index Showing Programs
	*/
	r.GET("/index", func(c *gin.Context) {
		programs := []model.Program{}
		db.Preload("Targets").Find(&programs)
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Index",
			"works": programs,
		})
	})

	/*
		Show the form to create a program
	*/
	r.GET("/new", func(c *gin.Context) {
		c.HTML(http.StatusOK, "new.tmpl", gin.H{
			"title": "New Work",
		})
	})

	/*
		Creates a new Program
	*/
	r.POST("/create-program", func(c *gin.Context) {
		name := c.PostForm("name")
		scope := c.PostForm("scope")
		logo := c.PostForm("logo")

		threads, _ := strconv.Atoi(c.PostForm("threads"))
		url := c.PostForm("url")

		targets := strings.Split(scope, "\n")
		targets = DeleteEmpty(targets)
		w := model.Program{Name: name, Threads: threads, URL: url, Logo: logo}
		for _, target := range targets {
			tmp := model.Target{Name: string(target)}
			db.Create(&tmp)
			db.Model(&tmp).Update("CreatedAt", time.Now())

			w.Targets = append(w.Targets, tmp)

		}

		db.Create(&w)
		db.Model(&w).Update("CreatedAt", time.Now())
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Error saving Work",
			})
		}

		c.Redirect(302, "/index")
	})

	/*
		Deletes an existing program (Soft Delete)
	*/
	r.GET("/delete/:id", func(c *gin.Context) {
		id := c.Param("id")

		program := model.Program{}
		db.Preload("Targets").First(&program, "id = ?", id)

		for _, target := range program.Targets {
			db.Delete(&target)
		}
		db.Delete(&program)
		c.Redirect(302, "/index")
	})

	/*
		Shows the info of a program
	*/
	r.GET("/view/:id", func(c *gin.Context) {
		id := c.Param("id")

		program := model.Program{}
		db.Preload("Targets").First(&program, "id = ?", id)

		c.HTML(http.StatusOK, "view.tmpl", gin.H{
			"program": program,
		})
	})

	r.GET("/api/fetchimg", func(c *gin.Context) {
		program := c.Query("program")

		query := `{"operationName":"LayoutDispatcher","variables":{"url":"` + program + `"},"query":"query LayoutDispatcher($url: URI!) {\n    resource(url: $url) {\n      ... on ResourceInterface {\n        ... on Team {\n          ...TeamProfileHeaderTeam\n        }\n      }\n    }\n  }\n  \n  fragment TeamProfileHeaderTeam on Team {\n    profile_picture(size: large)\n  }"}`

		r, err := http.Post("https://hackerone.com/graphql", "application/json", bytes.NewBuffer([]byte(query)))

		if err != nil {
			c.JSON(400, gin.H{
				"error": "Error fetching img",
			})
		}
		data, _ := ioutil.ReadAll(r.Body)

		response := model.H1Pictureresponse{}
		err = json.Unmarshal(data, &response)
		if err != nil {
			c.JSON(400, gin.H{
				"error": "Recovering img from response",
			})
		}
		c.JSON(200, gin.H{
			"pic": response.Data.Resource.Profile_picture,
		})

	})

	log.Printf("Listening on port %d", port)
	r.Run(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:8080
}
