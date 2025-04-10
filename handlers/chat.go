package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"os"
)

type llmPayload struct {
	Message string `json:"message"`
	Stream  bool   `json:"stream"`
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	aiURL := os.Getenv("AI_URL")
	if aiURL == "" {
		fmt.Println("AI_URL environment variable not set")
		return
	}

	message := r.FormValue("message")
	log.Printf("[ChatHandler] Received message: %s", message)

	// Ensure streaming is supported
	flusher, ok := w.(http.Flusher)
	if !ok {
		log.Println("[ChatHandler] Streaming unsupported")
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	// Prepare request to local LLM API
	payload := llmPayload{Message: message, Stream: true}
	body, _ := json.Marshal(payload)
	log.Printf("[ChatHandler] Payload: %s", string(body))

	req, err := http.NewRequest("POST", aiURL+"/api/v1/chat", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[ChatHandler] Failed to create request: %v", err)
		writeStreamError(w)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ChatHandler] Failed to send request: %v", err)
		writeStreamError(w)
		return
	}
	defer resp.Body.Close()

	log.Printf("[ChatHandler] LLM responded with status: %d", resp.StatusCode)

	// Start streaming raw chunks (escaped)
	w.Header().Set("Content-Type", "text/plain")
	buf := make([]byte, 1024)

	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			log.Printf("[ChatHandler] Read chunk: %q", chunk)

			lines := bytes.Split(chunk, []byte("\n"))
			for _, line := range lines {
				if bytes.HasPrefix(line, []byte("data: ")) {
					content := bytes.TrimPrefix(line, []byte("data: "))
					escaped := html.EscapeString(string(content))
					log.Printf("[ChatHandler] Streaming data: %s", escaped)
					fmt.Fprint(w, escaped)
					flusher.Flush()
				}
			}
		}

		if err == io.EOF {
			log.Println("[ChatHandler] Finished streaming (EOF)")
			break
		}
		if err != nil {
			log.Printf("[ChatHandler] Error while reading stream: %v", err)
			break
		}
	}
}

func writeStreamError(w http.ResponseWriter) {
	log.Println("[ChatHandler] Writing stream error fallback")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprint(w, "⚠️ Something went wrong while getting a response.")
}
