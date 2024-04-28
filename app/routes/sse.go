package routes

import (
	"context"
	"fmt"
	"net/http"
)

var msgChan chan string

func sseEndpoint(w http.ResponseWriter, r *http.Request) {
	fmt.Println("sseEndpoint... setting up")
	// get some id
	//

	// set up sse
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	msgChan = make(chan string)
	fmt.Println("sseEndpoint: msgChan = make(chan string)")

	// Create a context for handling client disconnection
	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	defer func() {
		if msgChan != nil {
			close(msgChan)
			msgChan = nil
		}
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	for {
		select {
		case message := <-msgChan:
			fmt.Fprintf(w, "%s", message)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func formatSSEMessage(msg string) string {
	var message string
	message = "event: message\n"
	// might need a loop here... and do data: %s\n for each line
	// then a single extra \n at the end
	message += fmt.Sprintf("data: %s\n\n", msg)

	return message
}

func sendSSEMessage(msg string) {
	if msgChan != nil {
		message := formatSSEMessage(msg)
		msgChan <- message
	}
}
