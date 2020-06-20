package handlers

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nachol/scanner/model"
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
	scanName := c.Param("scanName")

	// scanName := c.Param("scanName")
	program, err := model.GetProgramById(id)

	var raw string
	for key, val := range program.Scans {
		if val.Name == scanName {
			raw = program.Scans[key].Raw
			break
		}
	}

	if err != nil {
		log.Fatal(err)
	}
	c.HTML(http.StatusOK, "terminal.tmpl", gin.H{
		"raw": raw,
	})
}
