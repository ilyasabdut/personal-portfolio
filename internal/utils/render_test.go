package utils

import (
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestRenderTemplate(t *testing.T) {
	// Setup a recorder
	w := httptest.NewRecorder()

	// Ensure we are in the project root for relative paths to work
	// This is a common issue with Go tests depending on where they run
	// but standard `go test ./...` usually runs in the package dir.
	// However, RenderTemplate uses "web/templates/..." relative paths.
	// So we need to ensure the CWD is the project root.

	wd, _ := os.Getwd()
	// If we are in internal/utils, we need to go up two levels
	if strings.HasSuffix(wd, "internal/utils") {
		if err := os.Chdir("../../"); err != nil {
			t.Fatalf("Failed to chdir to project root: %v", err)
		}
		defer os.Chdir(wd) // Restore after test
	}

	data := map[string]interface{}{
		"AppEnv": "local",
	}

	// Act
	RenderTemplate(w, "index", data)

	// Assert
	if w.Code != 200 {
		t.Errorf("Expected status 200, got %d", w.Code)
		t.Logf("Response body: %s", w.Body.String())
	}

	body := w.Body.String()
	if !strings.Contains(body, "Ilyas Abduttawab") {
		t.Errorf("Expected body to contain 'Ilyas Abduttawab'")
	}
}
