package server

import (
	"github.com/ashirko/imlogs/internal/db"
	"github.com/ashirko/imlogs/internal/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
)

func addHandlers(router *gin.Engine) {
	router.GET("/logs", getLogs)
	router.GET("/logs/count", getLogsCount)
	router.POST("/logs", postLogs)
	router.POST("/logs/batch", postLogsBatch)
}

func getLogs(c *gin.Context) {
	params := make(map[string]interface{})

	// Get 'source' from the query string.
	if source, exists := c.GetQuery("source"); exists {
		params["source"] = source
	}

	// Get 'count' from the query string and convert it to int.
	if countStr, exists := c.GetQuery("count"); exists {
		if countInt, err := strconv.Atoi(countStr); err != nil {
			c.String(http.StatusBadRequest, "Count should be integer")
			return
		} else {
			params["count"] = countInt
		}
	}

	logs, err := db.GetDB().GetLogs(params)
	if err != nil {
		log.Println("Cannot fetch logs from the database: ", err)
		c.String(http.StatusInternalServerError, "DB query failed")
		return
	}

	c.JSON(200, gin.H{
		"logs": logs,
	})
}

func getLogsCount(c *gin.Context) {
	count, err := db.GetDB().GetLogsCount()
	if err != nil {
		log.Println("Cannot fetch logs count: ", err)
		c.String(http.StatusInternalServerError, "DB query failed")
		return
	}
	c.JSON(200, gin.H{
		"count": count,
	})
}

func postLogs(c *gin.Context) {
	var json models.LogLine
	err := c.BindJSON(&json)
	if err != nil {
		log.Println("BindJSON failed: ", err)
		c.String(http.StatusBadRequest, "Incorrect JSON")
		return
	}

	err = db.GetDB().AddLogLine(json)
	if err != nil {
		log.Println("Cannot save data: ", err)
		c.String(http.StatusInternalServerError, "Cannot save data")
		return
	}
	c.String(http.StatusOK, "ok")
}

func postLogsBatch(c *gin.Context) {
	var json models.LogLines
	err := c.BindJSON(&json)
	if err != nil {
		log.Println("BindJSON failed: ", err)
		c.String(http.StatusBadRequest, "Incorrect JSON")
		return
	}

	err = db.GetDB().AddLogLines(json)
	if err != nil {
		log.Println("Commit failed: ", err)
		c.String(http.StatusInternalServerError, "Cannot save data")
	}

	c.String(http.StatusOK, "ok")
}
