package routes

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
	"strings"
)

var (
	msgChan   = make(map[chan string]struct{})
	loginChan = make(map[chan string]struct{})
)

func sseEndpoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	evtChan := make(chan string)
	msgChan[evtChan] = struct{}{}

	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	defer func() {
		delete(msgChan, evtChan)
		close(evtChan)
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	for {
		select {
		case message := <-evtChan:
			fmt.Fprintf(w, "%s", message)
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}

func sseLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	evtChan := make(chan string)
	loginChan[evtChan] = struct{}{}

	_, cancel := context.WithCancel(r.Context())
	defer cancel()

	defer func() {
		delete(loginChan, evtChan)
		close(evtChan)
	}()

	flusher, ok := w.(http.Flusher)
	if !ok {
		return
	}

	for {
		select {
		case message := <-evtChan:
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
		for ch := range msgChan {
			ch <- message
		}
	}
}

func sendLoginMessage(msg string) {
	if loginChan != nil {
		message := formatSSEMessage(msg)
		for ch := range loginChan {
			ch <- message
		}
	}
}
