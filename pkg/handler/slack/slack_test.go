package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupChatRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/slack", Handler)
	return r
}

func TestSlackHandler(t *testing.T) {
	url := "/slack?token=1234&team_domain=knative&channel_id=C2147483705&channel_name=test&user_id=U2147483697&user_name=Steve&command=/tfeel&text=knative"
	router := setupChatRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", url, nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	log.Printf(w.Body.String())
}
