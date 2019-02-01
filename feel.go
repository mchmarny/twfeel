package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func feelHandler(c *gin.Context) {

	start := time.Now()

	query := c.Param("query")
	if query == "" {
		log.Println("nil id")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Required parameter missing: query",
			"status": http.StatusBadRequest,
		})
		return
	}

	result, err := search(c.Request.Context() , query)
	if err != nil {
		log.Printf("error on search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
			"status": http.StatusBadRequest,
		})
		return
	}

	elapsed := time.Since(start)
	log.Printf("Query duration: %s", elapsed)

	c.JSON(http.StatusOK, result)


}
