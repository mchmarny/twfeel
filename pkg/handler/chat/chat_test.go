package chat

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"io/ioutil"
	"bytes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupChatRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/chat", Handler)
	return r
}

func TestChatHandler(t *testing.T) {

	raw, err := ioutil.ReadFile("../../../sample/chat.json")
    if err != nil {
		t.Error(raw)
		return
	}

	router := setupChatRouter()

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/chat", bytes.NewReader(raw))
	router.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)

	log.Printf(w.Body.String())
}
