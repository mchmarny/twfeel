package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/twfeel/pkg/handler"

	"github.com/gin-gonic/gin"

)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"
)

func main() {

	log.Print("Initializing API server...")

	// router
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// api
	v1 := r.Group("/v1")
	{
		v1.GET("/feel/:query", handler.RestHandler)
		v1.POST("/chat", handler.ChatHandler)
	}

	// root
	r.GET("/", defaultHandler)

	// port
	port := os.Getenv(portVariableName)
	if port == "" {
		port = defaultPort
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Server starting: %s \n", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}

}

func defaultHandler(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}
