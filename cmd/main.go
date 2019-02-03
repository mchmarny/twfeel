package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mchmarny/twfeel/pkg/handler/chat"
	"github.com/mchmarny/twfeel/pkg/handler/slack"
	"github.com/mchmarny/twfeel/pkg/handler/rest"

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
		v1.GET("/feel/:query", rest.Handler)
		v1.POST("/chat", chat.Handler)
		v1.POST("/slack", slack.Handler)
	}

	// root & health
	r.GET("/", defaultHandler)
	r.GET("/health", defaultHandler)

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
