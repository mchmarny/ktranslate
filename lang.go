package main

import (
	"context"
	"log"

	lang "cloud.google.com/go/language/apiv1"
	langpb "google.golang.org/genproto/googleapis/cloud/language/v1"
)

func scoreSentiment(ctx context.Context, s string) (magnitude float32, score float32, err error) {

	client, err := lang.NewClient(ctx)
	if err != nil {
		log.Printf("Error while creating client: %v", err)
		return 0, 0, err
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
		return 0, 0, err
	}

	return result.DocumentSentiment.Magnitude,
		result.DocumentSentiment.Score, nil

}
