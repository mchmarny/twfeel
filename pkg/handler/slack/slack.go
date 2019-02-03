package slack

import (
	"log"
	"net/http"
	"strings"
	"fmt"
	"time"

	"github.com/mchmarny/twfeel/pkg/processor"
	"github.com/mchmarny/twfeel/pkg/common"

	"github.com/gin-gonic/gin"
)

const (
	tokenKey = "slack"
)

// Handler handles queries from chat service
func Handler(c *gin.Context) {

	// request
	post := Request{}
	err := c.ShouldBind(&post)
	if err != nil {
		log.Printf("invalid query string content: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid query string content",
			"status":  http.StatusBadRequest,
		})
		return
	}

	log.Printf("Request: %v", post)

	if post.Token == "" || post.Text == "" {
		log.Println("invalid request")
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request",
			"status":  http.StatusBadRequest,
		})
		return
	}

	log.Printf("Token: %s", post.Token)
	log.Printf("Domain: %s", post.Domain)
	log.Printf("ChannelName: %s", post.ChannelName)
	log.Printf("UserName: %s", post.UserName)

	// token
	if !common.IsValidAccessToken(tokenKey, post.Token) {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token",
			"status":  http.StatusBadRequest,
		})
		return
	}

	queryText := strings.Trim(post.Text, " ")
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
			sentiment = ":disappointed:"
		case common.PositiveSentiment:
			sentiment = ":smile:"
		default:
			sentiment = ":neutral_face:"
	}

	a1 := Attachment{}
	a1.Fields = []*Field{
		&Field{
			Title: "Sentiment",
			Value: fmt.Sprintf(
				"Based on analyses of last *%d* tweets the general sentiment about `%s` is %s",
				result.Tweets, queryText, sentiment),
		},

		&Field{
			Title: "Context",
			Value: fmt.Sprintf("Score *%.2f*, Magnitude *%.2f*, Sentiment *%s*",
			result.Score, result.Magnitude, result.Sentiment),
		},
	}

	a1.Actions = []*Action{
		&Action{
			Type: "button",
			Text: "Tweets",
			URL: fmt.Sprintf(
				"https://twitter.com/search?q=%s+-filter:retweets+until:%s",
				queryText, time.Now().Format("2006-01-02")),
			Style: "primary",
		},
	}

    pl := &Payload{
      Text: fmt.Sprintf("%s, here is the info you requested from <https://github.com/mchmarny/twfeel|twfeel>...", post.UserName),
      Username: "robot",
      Channel: "#general",
      IconEmoji: ":knative:",
      Attachments: []Attachment{a1},
    }

	c.JSON(http.StatusOK, pl)

}
