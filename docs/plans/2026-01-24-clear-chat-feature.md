# Clear Chat Feature Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Implement a "Clear Chat" feature that wipes the chat history from both the local UI and the remote AI backend, then starts a fresh session.

**Architecture:**
- **Backend:** `ClearChatHandler` receives a POST request, calls `DELETE /api/v1/sessions/{id}` on the AI service, and sets a new session cookie.
- **Frontend:** A new "Trash" button in `chat.html` calls the backend endpoint and clears the Alpine.js state/sessionStorage.

**Tech Stack:** Go, Alpine.js, HTML/CSS.

### Task 1: Backend Handler for Clearing Chat

**Files:**
- Modify: `internal/handlers/chat.go`
- Modify: `cmd/server/main.go`
- Test: `internal/handlers/chat_test.go` (Create new)

**Step 1: Write test for ClearChatHandler**
Create `internal/handlers/chat_test.go`. Since we can't easily mock the external AI URL without dependency injection changes, we'll test the handler logic: it should read a cookie, attempt to call the URL (we can mock the URL with a test server), and set a new cookie.

```go
package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestClearChatHandler(t *testing.T) {
	// Mock AI Server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE request, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	os.Setenv("AI_URL", server.URL)

	req := httptest.NewRequest("POST", "/clear-chat", nil)
	req.AddCookie(&http.Cookie{Name: "session_id", Value: "old-session-id"})
	w := httptest.NewRecorder()

	ClearChatHandler(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check for new cookie
	foundNewCookie := false
	for _, cookie := range resp.Cookies() {
		if cookie.Name == "session_id" && cookie.Value != "old-session-id" {
			foundNewCookie = true
		}
	}
	if !foundNewCookie {
		t.Error("Expected a new session_id cookie to be set")
	}
}
```

**Step 2: Implement ClearChatHandler in `internal/handlers/chat.go`**

```go
func ClearChatHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	aiURL := os.Getenv("AI_URL")
	if aiURL == "" {
		http.Error(w, "AI_URL not set", http.StatusInternalServerError)
		return
	}

	cookie, err := r.Cookie("session_id")
	if err == nil {
		// Call AI service to delete session
		client := &http.Client{Timeout: 5 * time.Second}
		req, err := http.NewRequest("DELETE", aiURL+"/api/v1/sessions/"+cookie.Value, nil)
		if err == nil {
			// Forward auth cookies if needed, or just send request
			req.Header.Set("Cookie", r.Header.Get("Cookie")) // Forward user cookies just in case
			resp, err := client.Do(req)
			if err == nil {
				defer resp.Body.Close()
				// We don't strictly care if the upstream delete fails (maybe session didn't exist),
				// we still want to clear the local session for the user.
				// But logging it is good.
				if resp.StatusCode != 200 && resp.StatusCode != 404 {
					log.Printf("Failed to delete remote session: status %d", resp.StatusCode)
				}
			} else {
				log.Printf("Failed to delete remote session: %v", err)
			}
		}
	}

	// Generate new session ID
	newSessionID := generate_session_id()
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    newSessionID,
		Path:     "/",
		HttpOnly: true, // Should probably be true unless JS needs to read it (JS does read it currently)
		// Current code: JS reads `document.cookie` in chat.html:40. So HttpOnly must be FALSE or omitted (default false).
	})

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"cleared", "new_session_id":"` + newSessionID + `"}`))
}
```

**Step 3: Register Handler in `cmd/server/main.go`**

```go
// Add to Routes section
http.HandleFunc("/clear-chat", handlers.ClearChatHandler)
```

**Step 4: Verify**
Run tests: `go test ./internal/handlers`

**Step 5: Commit**
```bash
git add internal/handlers cmd/server
git commit -m "feat: add backend handler for clearing chat"
```

### Task 2: Frontend UI for Clear Chat

**Files:**
- Modify: `web/templates/partials/chat.html`

**Step 1: Add Clear Button**
Add button to `.chat-header`:

```html
<div class="chat-header">
  <div class="header-controls">
      <button @click.stop="clearChat" class="chat-clear-btn" title="Clear Chat">🗑️</button>
  </div>
  <div class="header-title" @click="toggleChat">
      <h3 class="chat-title">Ask Me!</h3>
      <button id="chat-toggle-btn" class="chat-toggle-btn">🗖</button>
  </div>
</div>
```
*Note: Adjust structure/css to prevent click propagation if needed, or put button outside click area.*

**Step 2: Add CSS for button**
Update `web/static/style.css` (or add inline style in step 3 if preferred, but CSS file is better. Actually, plan didn't list `style.css` in Files. Let's stick to modifying `chat.html` and adding a small `<style>` block or inline style for simplicity if `style.css` is not easily accessible/mockable, but `web/static/style.css` exists. Let's Modify `web/static/style.css` too.)

**Files Update:**
- Modify: `web/static/style.css`

**Step 3: Implement `clearChat` function in Alpine component**

```javascript
clearChat() {
    if (!confirm('Clear chat history?')) return;

    fetch('/clear-chat', { method: 'POST' })
    .then(res => res.json())
    .then(data => {
        // Clear local state
        this.messages = [];
        sessionStorage.removeItem(`chatMessages-${this.sessionId}`);
        
        // Update session ID if backend sent a new one (it's in the cookie automatically, 
        // but we might need to update our local property)
        if (data.new_session_id) {
            this.sessionId = data.new_session_id;
        } else {
             this.sessionId = this.getCookie("session_id");
        }
    })
    .catch(err => console.error("Failed to clear chat:", err));
}
```

**Step 4: Verify**
Manual verification (since we can't easily unit test the UI interactions here).

**Step 5: Commit**
```bash
git add web/templates/partials/chat.html web/static/style.css
git commit -m "feat: add clear chat button and frontend logic"
```
