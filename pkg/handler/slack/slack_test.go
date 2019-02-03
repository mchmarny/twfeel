package slack

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter() *gin.Engine {
	r := gin.Default()
	r.POST("/slack", Handler)
	return r
}

func TestSlackHandler(t *testing.T) {
	body1 := []byte(`token=test123&team_id=T1&team_domain=knative&channel_id=C1&channel_name=directmessage&user_id=U1&user_name=u1&command=%2Fkfeel&text=knative&response_url=https%3A%2F%2Ftest.domain.com&trigger_id=1.2.3`)
	submitSlackHandler(t, body1)

	body2 := []byte(`token=test123&team_id=T1&team_domain=knative&channel_id=C1&channel_name=directmessage&user_id=U1&user_name=u1&command=%2Fkfeel&text=knative&response_url=https%3A%2F%2Ftest.domain.com&trigger_id=1.2.3`)
	submitSlackHandler(t, body2)
}

func submitSlackHandler(t *testing.T, body []byte) {
	router := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/slack", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	log.Printf(w.Body.String())
}
