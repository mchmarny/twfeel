package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"net/http/httputil"

	"github.com/gin-gonic/gin"
)

const (
	defaultPort      = "8080"
	portVariableName = "PORT"
)

func main() {

	log.Print("Initializing twfeel API server...")

	// router
	r := gin.Default()

	// api
	v1 := r.Group("/v1")
	{
		v1.GET("/feel/:query", feelHandler)
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

func withLog(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Header().Set("Expires", "0")

		reqDump, err := httputil.DumpRequest(r, true)
		if err != nil {
			log.Println(err)
		} else {
			log.Println(string(reqDump))
		}

		next.ServeHTTP(w, r)
	}
}
