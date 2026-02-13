# Chat UI Overhaul Implementation Plan

> **For Claude:** REQUIRED SUB-SKILL: Use superpowers:executing-plans to implement this plan task-by-task.

**Goal:** Redesign the chat UI for a cleaner, more minimal look and replace the browser "Clear Chat" alert with a non-intrusive two-step confirmation.

**Architecture:**
- **UI:** Reorder header elements, remove boxy borders, and add hover-only opacity.
- **Logic:** Add `confirmingClear` state to the Alpine.js component to handle two-step deletion without browser popups.

**Tech Stack:** Alpine.js, CSS, HTML.

### Task 1: Redesign Chat Header and State

**Files:**
- Modify: `web/templates/partials/chat.html`

**Step 1: Update Alpine state and `clearChat` logic**
Add `confirmingClear: false` to the returned object.
Update `clearChat()`:

```javascript
clearChat() {
    if (!this.confirmingClear) {
        this.confirmingClear = true;
        setTimeout(() => { this.confirmingClear = false; }, 3000);
        return;
    }
    
    fetch("/clear-chat", { method: "POST" })
      .then((res) => {
        if (res.ok) {
          this.messages = [];
          sessionStorage.removeItem(`chatMessages-${this.sessionId}`);
          this.confirmingClear = false;
          return res.json().catch(() => ({}));
        }
      })
      .then((data) => {
        if (data && data.new_session_id) {
          this.sessionId = data.new_session_id;
        }
      })
      .catch((err) => {
        console.error("Failed to clear chat:", err);
        this.confirmingClear = false;
      });
}
```

**Step 2: Update Header HTML Structure**
Rearrange the header to put the title on the left and controls on the right.
Update the trash button to show "Clear?" when in confirmation state.

```html
<div class="chat-header">
    <div class="header-title" @click="toggleChat">
        <h3 class="chat-title">Ask Me!</h3>
    </div>
    <div class="header-controls">
        <button @click.stop="clearChat" 
                class="chat-clear-btn" 
                :class="{ 'confirming': confirmingClear }"
                :title="confirmingClear ? 'Click again to confirm' : 'Clear Chat'">
            <template x-if="!confirmingClear"><span>🗑️</span></template>
            <template x-if="confirmingClear"><span class="confirm-text">Clear?</span></template>
        </button>
        <button id="chat-toggle-btn" class="chat-toggle-btn" @click="toggleChat">🗖</button>
    </div>
</div>
```

**Step 3: Commit**
```bash
git add web/templates/partials/chat.html
git commit -m "feat: redesign chat header and implement two-step clear"
```

### Task 2: Polish CSS Aesthetics

**Files:**
- Modify: `web/static/style.css`

**Step 1: Overhaul Chat Styles**
- Fix header flex alignment.
- Style the `confirm-text` and `confirming` state.
- Improve bubble spacing and minimal look.

```css
.chat-header {
  padding: 0.5rem 1rem;
  border-bottom: 1px solid var(--table-border);
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.header-title {
  cursor: pointer;
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.chat-clear-btn {
  font-size: 0.9rem;
  padding: 2px 6px;
  border-radius: 4px;
}

.chat-clear-btn.confirming {
  color: #cc0000;
  font-weight: bold;
}

.confirm-text {
  font-size: 0.8rem;
}

/* Minimal bubbles */
.bubble {
  padding: 0.4rem 0.8rem;
  border-radius: 8px; /* Slightly more boxy to match site */
  font-size: 0.9rem;
}
```

**Step 2: Verify**
Check file contents.

**Step 3: Commit**
```bash
git add web/static/style.css
git commit -m "style: polish chat UI aesthetics"
```
