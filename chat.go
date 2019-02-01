package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func chatHandler(c *gin.Context) {

	token := c.PostForm("token")
	if token != knownToken {
		log.Printf("invalid token. Got:%s Expected:%s", token, knownToken)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid access token",
			"status":  http.StatusBadRequest,
		})
		return
	}

	msg := Message{}
	err := c.ShouldBind(&msg)

	if err == nil {
		log.Printf("invalid body content: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid body content",
			"status":  http.StatusBadRequest,
		})
		return
	}

	log.Printf("Message: %v", msg)

	// result, err := search(c.Request.Context(), query)
	// if err != nil {
	// 	log.Printf("error on search: %v", err)
	// 	c.JSON(http.StatusBadRequest, gin.H{
	// 		"message": err,
	// 		"status":  http.StatusBadRequest,
	// 	})
	// 	return
	// }

	hd := &Header{
		Title: fmt.Sprintf("Hello: %s", msg.Sender.DisplayName),
	}

	sn := &Section{
		Widgets: []*Widget{
			&Widget{
				TextParagraph: &TextParagraph{
					Text: "query sentiment here",
				},
				Image: &Image{
					ImageURL: msg.Sender.AvatarURL,
				},
			},
		},
	}

	rez := []*Card{
		&Card{
			Header:   hd,
			Sections: []*Section{sn},
		},
	}

	c.JSON(http.StatusOK, rez)

}
