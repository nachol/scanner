package handlers

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nachol/scanner/model"
)

func FetchH1Image(c *gin.Context) {
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

}
