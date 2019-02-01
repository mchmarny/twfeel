package main

import (
	"fmt"
	"log"
	"net/http"

	// "strings"

	"github.com/gin-gonic/gin"
)

func chatHandler(c *gin.Context) {

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
	if post.Token != knownToken {
		log.Printf("invalid token. Expected:%s Got:%s", knownToken, post.Token)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token",
			"status":  http.StatusBadRequest,
		})
		return
	}

	if post.Message == nil || post.Message.Sender == nil {
		log.Println("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"status":  http.StatusBadRequest,
		})
		return
	}

	senderName := post.Message.Sender.DisplayName
	queryText := post.Message.Text
	log.Printf("Query: %s", queryText)


	// run query and score
	result, err := search(c.Request.Context(), queryText)
	if err != nil {
		log.Printf("error on search: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err,
			"status":  http.StatusBadRequest,
		})
		return
	}


	// format results
	sentiment := "positive"
	if result.Score < 0 {
		sentiment = "negative"
	}

	txt := "Hi %s, I ran analyses on last *%d* tweets related to `%s` and the general sentiment is *%s* (magnitude: *%f*)%s"

	rateLimitInfo := ""
	if result.Score == 0 {
		rateLimitInfo = " - Twitter rate-limited :("
	}

	rez := &Message{
		Text: fmt.Sprintf(txt, senderName, result.Tweets, queryText, sentiment, result.Magnitude, rateLimitInfo),
	}

	c.JSON(http.StatusOK, rez)

}
