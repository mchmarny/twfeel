package handler

import (
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/twfeel/pkg/processor"

	"github.com/gin-gonic/gin"
)

var (
	knownToken = os.Getenv("ACCESS_TOKEN")
)

// RestHandler handles queries from the REST interface
func RestHandler(c *gin.Context) {

	token := c.Query("token")
	if token != knownToken {
		log.Printf("invalid token. Got:%s Expected:%s", token, knownToken)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token",
			"status":  http.StatusBadRequest,
		})
		return
	}

	query := c.Param("query")
	if query == "" {
		log.Println("nil id")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Required parameter missing: query",
			"status":  http.StatusBadRequest,
		})
		return
	}

	result, err := processor.Search(c.Request.Context(), query)
	if err != nil {
		log.Printf("error on search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
			"status":  http.StatusBadRequest,
		})
		return
	}

	c.JSON(http.StatusOK, result)

}