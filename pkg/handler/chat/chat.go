package chat

import (
	"fmt"
	"log"
	"time"
	"net/http"
	"strings"

	"github.com/mchmarny/twfeel/pkg/processor"
	"github.com/mchmarny/twfeel/pkg/common"

	"github.com/gin-gonic/gin"
)

const (
	tokenKey = "chat"
)

// Handler handles queries from chat service
func Handler(c *gin.Context) {

	// request
	post := Request{}
	err := c.BindJSON(&post)
	if err != nil {
		log.Printf("invalid body content: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid body content",
			"status":  http.StatusBadRequest,
		})
		return
	}

	log.Printf("MessageType: %s", post.MessageType)
	log.Printf("EventTime: %s", post.EventTime)
	log.Printf("Token: %s", post.Token)

	// token
	if !common.IsValidAccessToken(tokenKey, post.Token) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token",
			"status":  http.StatusBadRequest,
		})
		return
	}

	log.Printf("Request: %v", post)

	if post.Message == nil || post.Message.Sender == nil {
		log.Println("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"status":  http.StatusBadRequest,
		})
		return
	}

	senderName := post.Message.Sender.Name
	queryText := strings.Trim(post.Message.ArgumentText, " ")
	log.Printf("Query: %s", queryText)

	// run query and score
	result, err := processor.Search(c.Request.Context(), queryText)
	if err != nil {
		log.Printf("error on search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
			"status":  http.StatusBadRequest,
		})
		return
	}

	if result == nil {
		log.Fatal("BUG, result should never be nil")
	}

	// format results
	sentiment := ""
	switch result.Sentiment {
		case common.NegativeSentiment:
			sentiment = "`:(` negative"
		case common.PositiveSentiment:
			sentiment = "`:)` *positive*"
		default:
			sentiment = "`:|` *meh*"
	}

	txt := "Hi <%s>, I ran analyses on last *%d* tweets related to `%s` and the general sentiment is %s  -- meta: score *%.2f*, magnitude *%.2f*, <https://twitter.com/search?q=%s+-filter:retweets+until:%s|tweets>"
	txt =  fmt.Sprintf(txt, senderName, result.Tweets, queryText, sentiment, result.Score, result.Magnitude, queryText, time.Now().Format("2006-01-02"))

	rez := &Message{Text: txt}

	c.JSON(http.StatusOK, rez)

}
