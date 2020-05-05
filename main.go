package main

import (
	"context"

	model "./model"
	scan "./scan"
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	// "reflect"
	"strconv"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

	// Mongo connection
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	// Check the connection
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to MongoDB!")

	db := client.Database("scanner")
	collectionProgram := db.Collection("program")

	/*
		Index Showing Programs
	*/
	r.GET("/index", func(c *gin.Context) {
		var results []*model.Program

		cur, err := collectionProgram.Find(context.TODO(), bson.D{{}})
		if err != nil {
			log.Fatal(err)
		}

		for cur.Next(context.TODO()) {

			// create a value into which the single document can be decoded
			var elem model.Program
			err := cur.Decode(&elem)
			if err != nil {
				log.Fatal(err)
			}

			results = append(results, &elem)
		}

		if err := cur.Err(); err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title":    "Index",
			"programs": results,
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
			w.Targets = append(w.Targets, tmp)

		}
		w.CreatedAt = time.Now()

		// b, err := json.Marshal(&w)
		if err != nil {
			return
		}

		insertResult, err := collectionProgram.InsertOne(context.TODO(), w)
		if err != nil {
			log.Fatal(err)
		}

		log.Printf("Success: %s\n", insertResult)

		c.Redirect(302, "/index")
	})

	/*
		Deletes an existing program (Soft Delete)
	*/
	r.GET("/delete/:id", func(c *gin.Context) {

		id := c.Param("id")
		filter := bson.D{{"name", id}}
		_, err := collectionProgram.DeleteOne(context.TODO(), filter)
		if err != nil {
			log.Fatal(err)
		}

		c.Redirect(302, "/index")
	})

	/*
		Shows the info of a program
	*/
	r.GET("/view/:id", func(c *gin.Context) {
		id := c.Param("id")
		var program model.Program

		filter := bson.D{{"name", id}}
		err = collectionProgram.FindOne(context.TODO(), filter).Decode(&program)
		if err != nil {
			log.Fatal(err)
		}

		c.HTML(http.StatusOK, "view.tmpl", gin.H{
			"program": program,
		})
	})

	/*
		Shows the info of a program
	*/
	r.GET("/terminal/:id/:scanName", func(c *gin.Context) {
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
		err := collectionProgram.FindOne(context.TODO(), filter).Decode(&program)

		if err != nil {
			log.Fatal(err)
		}
		c.HTML(http.StatusOK, "terminal.tmpl", gin.H{
			"raw": program.Scans[0].Raw,
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

	r.POST("/runScan", func(c *gin.Context) {
		id := c.PostForm("programID")
		var program model.Program
		filter := bson.D{{"name", id}}
		err = collectionProgram.FindOne(context.TODO(), filter).Decode(&program)
		if err != nil {
			log.Fatal(err)
		}

		var scan scan.Scan
		c.Bind(&scan)
		options := c.PostFormMap("options")
		scan.Options = options
		scan.Threads = program.Threads
		scan.LoadModes()
		_, err := scan.Run()
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
		updateResult, err := collectionProgram.UpdateOne(context.TODO(), filter, update)

		if err != nil {
			log.Fatal(err)
		}

		c.JSON(200, gin.H{
			"message": "Work Done",
			"result":  updateResult,
		})
	})

	log.Printf("Listening on port %d", port)
	r.Run(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:8080
}
