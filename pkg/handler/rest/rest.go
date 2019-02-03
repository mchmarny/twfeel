package rest

import (
	"log"
	"net/http"

	"github.com/mchmarny/twfeel/pkg/processor"
	"github.com/mchmarny/twfeel/pkg/common"

	"github.com/gin-gonic/gin"
)

// Handler handles queries from the REST interface
func Handler(c *gin.Context) {

	// token
	if !common.IsValidAccessToken(c.Query("token")) {
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
