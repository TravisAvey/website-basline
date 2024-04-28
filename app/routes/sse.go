package routes

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
)

var msgChan chan string

func sseEndpoint(w http.ResponseWriter, r *http.Request) {
	// get some id

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	msgChan = make(chan string)

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

	scanner := bufio.NewScanner(strings.NewReader(msg))
	for scanner.Scan() {
		message += fmt.Sprintf("data: %s\n", scanner.Text())
	}
	message += "\n"

	return message
}

func sendSSEMessage(msg string) {
	if msgChan != nil {
		message := formatSSEMessage(msg)
		msgChan <- message
	}
}
