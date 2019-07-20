package main

import (
	"context"
	"fmt"
	"log"

	ce "github.com/cloudevents/sdk-go"
	ev "github.com/mchmarny/gcputil/env"
	"github.com/tidwall/gjson"
)

var (
	textPath = ev.MustGetEnvVar("TEXT_PATH", "text")
)

func newEventReceiver(ctx context.Context) (er *eventReceiver, err error) {
	s, err := newTransaltionService(ctx)
	return &eventReceiver{
		ts: s,
	}, nil
}

type eventReceiver struct {
	ts *transaltionService
}

func (r *eventReceiver) Receive(ctx context.Context, event ce.Event, resp *ce.EventResponse) error {

	//log.Printf("Raw Event: %v", event)

	// get content
	var textValue string
	if event.DataContentType() == "text/plain" {
		if err := event.DataAs(textValue); err != nil {
			log.Printf("Failed to DataAs string: %s", err.Error())
			return err
		}
	} else if event.DataContentType() == "application/json" {
		content, err := event.DataBytes()
		if err != nil {
			log.Printf("Failed to DataAs bytes: %s", err.Error())
			return err
		}
		textValue = gjson.GetBytes(content, textPath).String()
	} else {
		return fmt.Errorf("Invalid Data Content Type: %s. Only application/json and text/plain supported",
			event.DataContentType())
	}

	log.Printf("Text to translate: %s", textValue)

	// empty handler
	if textValue == "" {
		return nil
	}

	translatedText, err := r.ts.tranlate(ctx, textValue, "en")
	if err != nil {
		log.Printf("Failed to translate text %s. Error: %s", textValue, err.Error())
		return err
	}

	// set the extension
	event.SetExtension("translation", map[string]string{
		"original": textValue,
		"result":   translatedText,
	})

	log.Printf("Translated event: %v", event.Context)
	resp.RespondWith(200, &event)

	return nil

}
