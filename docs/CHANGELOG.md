# Changelog

## Refactor: Project Restructure to Standard Go Layout

### 🏗️ Architecture Changes
- **Standard Go Layout**: Reorganized project files to follow standard Go conventions.
  - `cmd/server/`: New location for the application entry point (`main.go`).
  - `internal/`: New location for private application code (`handlers`, `database`, `utils`).
  - `web/`: New location for static assets and templates (`static`, `templates`).
  - `scripts/`: Dedicated folder for shell scripts.

### ✨ Features
- **Makefile**: Added a `Makefile` with commands for common tasks:
  - `make run`: Run the server locally.
  - `make build`: Build the binary to `bin/server`.
  - `make docker-up`: Start the application with Docker Compose.
  - `make clean`: Clean build artifacts.

### 🛠️ Improvements & Fixes
- **Cross-Platform Paths**: Updated file path handling in `utils` and `main.go` to use `filepath.Join` for better Windows/Linux compatibility.
- **Docker Update**: Updated `Dockerfile` and `docker-compose.yml` to reflect the new directory structure.
- **Tests**: Added unit tests for template rendering (`internal/utils/render_test.go`).

### 📦 Files Changed
- Moved `main.go` → `cmd/server/main.go`
- Moved `handlers/` → `internal/handlers/`
- Moved `database/` → `internal/database/`
- Moved `utils/` → `internal/utils/`
- Moved `static/` & `templates/` → `web/`

## Feature: Clear Chat & UI Overhaul

### ✨ Features
- **Clear Chat**: Added a trash icon to the chat header to wipe conversation history from both the local UI and the AI backend.
- **Two-Step Confirmation**: Implemented a non-intrusive "Click again to confirm" mechanism for clearing chat, replacing browser alerts.

### 🎨 UI & UX Improvements
- **Minimal Redesign**: Overhauled the chat interface with cleaner alignment, boxy bubbles to match the site's aesthetic, and hover-only controls.
- **Alpine.js Refactor**: Improved component logic using idiomatic Alpine.js state management and safe timer handling.
- **DRY CSS**: Refactored `style.css` to unify input styles and remove redundant rules.

### 🛠️ Improvements & Fixes
- **Backend Handler**: Added a Go handler for session deletion with automatic new-session cookie generation.
- **Cookie Accessibility**: Set `HttpOnly: false` for the `session_id` cookie to allow frontend initialization.

