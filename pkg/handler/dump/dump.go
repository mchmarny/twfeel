package dump

import (
	"log"
	"net/http"
	"io/ioutil"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	tokenKey = "rest"
)

// Handler handles queries from the REST interface
func Handler(c *gin.Context) {


	// body first
	body, _ := ioutil.ReadAll(c.Request.Body)
	log.Printf("Body: %s", string(body))

	// query
	log.Println("Query:")
	for k, v := range c.Request.URL.Query() {
		log.Printf("\t %s: %s", k, strings.Join(v, ","))
	}

	//headers
	log.Println("Headers:")
	for k, vs := range c.Request.Header {
		for _, v := range vs {
			log.Printf("\t %s: %s", k, v)
		}
	}

	c.JSON(http.StatusOK, nil)

}
