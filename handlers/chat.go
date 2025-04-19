package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
	"github.com/google/uuid"
)

type llmPayload struct {
	Message string `json:"message"`
	Stream  bool   `json:"stream"`
}

// Function to generate a new session ID (UUID)
func generate_session_id() string {
	return uuid.New().String()
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	aiURL := os.Getenv("AI_URL")
	llmModel := os.Getenv("LLM_MODEL")

	if aiURL == "" {
		http.Error(w, "AI_URL environment variable not set", http.StatusInternalServerError)
		return
	}

	// Define session_id variable
	var session_id string

	// Get session cookie
	cookie, err := r.Cookie("session_id")
	if err != nil {
		log.Println("No session_id cookie found or error:", err)
		session_id = generate_session_id() // Generate new session_id if not found
	} else {
		session_id = cookie.Value // Otherwise, use the existing cookie value
	}
	log.Printf("[ChatHandler] Using session_id: %s", session_id)

	message := r.FormValue("message")
	log.Printf("[ChatHandler] Received message: %s", message)

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	payload := llmPayload{Message: message, Stream: true}
	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", aiURL+"/api/v1/chat", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[ChatHandler] Failed to create request: %v", err)
		writeStreamError(w)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	if llmModel != "" {
		req.Header.Set("LLM_MODEL", llmModel)
	}
	// Forward client cookies
	req.Header.Set("Cookie", r.Header.Get("Cookie"))

	client := &http.Client{Timeout: 0}
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("[ChatHandler] Request failed: %v", err)
		writeStreamError(w)
		return
	}
	defer resp.Body.Close()

	// Forward Set-Cookie headers from FastAPI to client
	for _, cookie := range resp.Cookies() {
		http.SetCookie(w, cookie)
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	fmt.Fprint(w, ": Connected\n\n")
	flusher.Flush()

	buf := make([]byte, 1024)
	for {
		n, err := resp.Body.Read(buf)
		if n > 0 {
			chunk := buf[:n]
			w.Write(chunk)
			flusher.Flush()
			time.Sleep(10 * time.Millisecond)
		}

		if err == io.EOF {
			fmt.Fprintf(w, "event: complete\ndata: \n\n")
			flusher.Flush()
			break
		}
		if err != nil {
			log.Printf("[ChatHandler] Stream read error: %v", err)
			fmt.Fprintf(w, "event: error\ndata: Stream reading error\n\n")
			flusher.Flush()
			break
		}
	}
}

func writeStreamError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "event: error\ndata: Something went wrong while getting a response.\n\n")
}

func SessionHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("[SessionHandler] Received /session request")

	existingCookies := r.Header.Get("Cookie") // Get the cookies from the client

	apiURL := os.Getenv("AI_URL") + "/api/v1/auth/get-session"
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		http.Error(w, "Failed to create session request", http.StatusInternalServerError)
		return
	}
	req.Header.Set("Cookie", existingCookies) // Forward the cookies from the client

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		http.Error(w, "Failed to fetch session", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Forward Set-Cookie headers from FastAPI to the client
	for _, cookie := range resp.Cookies() {
		http.SetCookie(w, cookie)
	}

	for k, vv := range resp.Header {
		if k == "Set-Cookie" {
			continue
		}
		for _, v := range vv {
			w.Header().Add(k, v)
		}
	}

	w.WriteHeader(resp.StatusCode)
	io.Copy(w, resp.Body)
}
