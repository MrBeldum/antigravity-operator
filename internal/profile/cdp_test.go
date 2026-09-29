package profile_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/tiagoboas/antigravity-operator/internal/profile"
)

func TestCDP_ListTabs_Mock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/json" {
			tabs := []profile.Tab{
				{
					ID:    "tab1",
					Title: "Google AI",
					Type:  "page",
					URL:   "https://discuss.ai.google.dev",
				},
				{
					ID:    "sw1",
					Title: "Service Worker",
					Type:  "service_worker",
					URL:   "https://discuss.ai.google.dev/sw.js",
				},
			}
			_ = json.NewEncoder(w).Encode(tabs)
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	parts := strings.Split(server.URL, ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	tabs, err := profile.ListTabs(port)
	if err != nil {
		t.Fatalf("ListTabs failed: %v", err)
	}

	if len(tabs) != 1 {
		t.Fatalf("expected 1 page tab (filtering out service_worker), got %d", len(tabs))
	}
	if tabs[0].Title != "Google AI" {
		t.Errorf("expected title 'Google AI', got '%s'", tabs[0].Title)
	}
}

func TestCDP_OpenAndCloseTab_Mock(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/json/new") {
			_ = json.NewEncoder(w).Encode(profile.Tab{
				ID:    "new-tab-id",
				Title: "New Tab",
				Type:  "page",
				URL:   "https://github.com",
			})
			return
		}
		if strings.HasPrefix(r.URL.Path, "/json/close/") {
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("Target is closing"))
			return
		}
		http.NotFound(w, r)
	}))
	defer server.Close()

	parts := strings.Split(server.URL, ":")
	port, _ := strconv.Atoi(parts[len(parts)-1])

	tab, err := profile.OpenTab(port, "https://github.com")
	if err != nil {
		t.Fatalf("OpenTab failed: %v", err)
	}
	if tab.ID != "new-tab-id" {
		t.Errorf("expected new-tab-id, got %s", tab.ID)
	}

	if err := profile.CloseTab(port, "new-tab-id"); err != nil {
		t.Fatalf("CloseTab failed: %v", err)
	}
}

func TestCDP_Eval_NoTabs(t *testing.T) {
	// Inactive port
	_, err := profile.Eval(59997, "1+1")
	if err == nil {
		t.Error("expected error for inactive port on Eval, got nil")
	}
}

func TestCDP_Screenshot_NoTabs(t *testing.T) {
	// Inactive port
	err := profile.Screenshot(59997, "/tmp/should_fail.png")
	if err == nil {
		t.Error("expected error for inactive port on Screenshot, got nil")
	}
}
