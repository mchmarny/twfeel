package dump

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"log"
	"bytes"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/dump", Handler)
	return r
}

func TestDumpHandler(t *testing.T) {
	url := "/dump?token=1234&k1=v1&k2=v2"
	router := setupRouter()
	w := httptest.NewRecorder()
	var jsonStr = []byte(`{"text":"my test"}`)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Add("X-Custom-Header-1", "header1")
	req.Header.Add("X-Custom-Header-2", "header2")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	log.Printf(w.Body.String())
}
