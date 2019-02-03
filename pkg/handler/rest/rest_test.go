package rest

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.GET("/feel/:query", Handler)
	return r
}

func TestRestHandler(t *testing.T) {
	url := "/feel/knative?token=1234"
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", url, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	log.Printf(w.Body.String())
}
