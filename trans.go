package main

import (
	"context"

	"cloud.google.com/go/translate"
	"golang.org/x/text/language"
)

func newTransaltionService(ctx context.Context) (srv *transaltionService, err error) {

	c, err := translate.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	return &transaltionService{
		client: c,
	}, nil

}

type transaltionService struct {
	client *translate.Client
}

func (t *transaltionService) tranlate(ctx context.Context, fromText, toLang string) (result string, err error) {

	// Sets the target language.
	target, err := language.Parse(toLang)
	if err != nil {
		return "", err
	}

	// Translates the text
	r, err := t.client.Translate(ctx, []string{fromText}, target, nil)
	if err != nil {
		return "", err
	}

	return r[0].Text, nil
}
