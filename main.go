package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"

	ce "github.com/cloudevents/sdk-go"
	ev "github.com/mchmarny/gcputil/env"
)

func main() {

	port, err := strconv.Atoi(ev.MustGetEnvVar("PORT", "8080"))
	if err != nil {
		log.Fatalf("failed to parse port, %s", err.Error())
	}

	// Handler Mux
	mux := http.NewServeMux()

	// Ingres API Handler
	t, err := ce.NewHTTPTransport(
		ce.WithMethod("POST"),
		ce.WithPath("/"),
		ce.WithPort(port),
	)
	if err != nil {
		log.Fatalf("failed to create CloudEvents transport, %s", err.Error())
	}

	er, err := newEventReceiver(context.Background())
	if err != nil {
		log.Fatalf("failed to create event receiver, %s", err.Error())
	}

	// wire handler for CE
	t.SetReceiver(er)

	// Health Handler
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		fmt.Fprint(w, "ok")
	})

	// Events or UI Handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Method, %s", r.Method)
		if r.Method == "POST" {
			t.ServeHTTP(w, r)
			return
		}
		fmt.Fprint(w, "Nothing to see here. Use POST to send CloudEvents")
	})

	a := fmt.Sprintf(":%d", port)
	log.Fatal(http.ListenAndServe(a, mux))

}
