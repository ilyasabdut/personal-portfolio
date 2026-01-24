package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestClearChatHandler(t *testing.T) {
	// Mock the AI Backend
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Expect DELETE /api/v1/sessions/{session_id}
		expectedPathPrefix := "/api/v1/sessions/"
		if r.Method != http.MethodDelete {
			t.Errorf("Backend expected DELETE, got %s", r.Method)
		}
		if len(r.URL.Path) <= len(expectedPathPrefix) {
			t.Errorf("Backend expected session ID in path")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status":"success"}`))
	}))
	defer backendServer.Close()

	// Set AI_URL env var to mock server
	os.Setenv("AI_URL", backendServer.URL)
	defer os.Unsetenv("AI_URL")

	// Create request with existing session cookie
	req := httptest.NewRequest("POST", "/clear-chat", nil)
	existingSessionID := "old-session-id"
	cookie := &http.Cookie{
		Name:  "session_id",
		Value: existingSessionID,
	}
	req.AddCookie(cookie)

	// Recorder for response
	rr := httptest.NewRecorder()

	// Call Handler
	handler := http.HandlerFunc(ClearChatHandler)
	handler.ServeHTTP(rr, req)

	// Verify Status Code
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Verify Response Body
	var response map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &response); err != nil {
		t.Errorf("handler returned invalid JSON: %v", err)
	}

	if response["status"] != "cleared" {
		t.Errorf("handler returned wrong status: got %v want %v",
			response["status"], "cleared")
	}

	if response["new_session_id"] == "" {
		t.Errorf("handler did not return new_session_id")
	}

	if response["new_session_id"] == existingSessionID {
		t.Errorf("handler did not generate a NEW session ID")
	}

	// Verify Set-Cookie header
	cookies := rr.Result().Cookies()
	foundSessionCookie := false
	for _, c := range cookies {
		if c.Name == "session_id" {
			foundSessionCookie = true
			if c.Value != response["new_session_id"] {
				t.Errorf("cookie value does not match response session ID")
			}
			// Basic check for expiration (should be in future)
			if c.Expires.Before(time.Now()) {
				t.Errorf("cookie expiration should be in the future")
			}
		}
	}
	if !foundSessionCookie {
		t.Errorf("handler did not set session_id cookie")
	}
}
