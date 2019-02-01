package main

import (
	"context"
	"log"

	lang "cloud.google.com/go/language/apiv1"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func scoreSentiment(ctx context.Context, s string) (sentiment *langpb.Sentiment, err error) {

	client, err := lang.NewClient(ctx)
	if err != nil {
		log.Panicf("Failed to create client: %v", err)
	}

	result, err := client.AnalyzeSentiment(ctx, &langpb.AnalyzeSentimentRequest{
		Document: &langpb.Document{
			Source: &langpb.Document_Content{
				Content: s,
			},
			Type: langpb.Document_PLAIN_TEXT,
		},
		EncodingType: langpb.EncodingType_UTF8,
	})
	if err != nil {
		log.Printf("Error while scoring: %v", err)
		return nil, err
	}

	return result.DocumentSentiment, nil

}
