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
	body1 := []byte(`token=DGAZAc5lF2gwnomWDXjsSwK9&team_id=T93ELUK42&team_domain=knative&channel_id=D93E7DHT6&channel_name=directmessage&user_id=U94GEEP9V&user_name=mchmarny&command=%2Fkfeel&text=knative&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT93ELUK42%2F540718281363%2FI4JCetnXwIE5Ryna7IMLX77Q&trigger_id=540718281763.309496971138.252b8586d3df8f40f1ff381ded8ec20d`)
	submitSlackHandler(t, body1)

	body2 := []byte(`token=DGAZAc5lF2gwnomWDXjsSwK9&team_id=T93ELUK42&team_domain=knative&channel_id=D93E7DHT6&channel_name=directmessage&user_id=U94GEEP9V&user_name=mchmarny&command=%2Fkfeel&text=knative&response_url=https%3A%2F%2Fhooks.slack.com%2Fcommands%2FT93ELUK42%2F540718281363%2FI4JCetnXwIE5Ryna7IMLX77Q&trigger_id=540718281763.309496971138.252b8586d3df8f40f1ff381ded8ec20d`)
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
