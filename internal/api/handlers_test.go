package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/igodwin/secretsanta/pkg/participant"
)

func TestHandleValidate(t *testing.T) {
	server := NewServer(":8080")

	tests := []struct {
		name           string
		participants   []*participant.Participant
		expectedValid  bool
		expectedStatus int
	}{
		{
			name: "Valid configuration",
			participants: []*participant.Participant{
				{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
				{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
			},
			expectedValid:  true,
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid - participant with no options",
			participants: []*participant.Participant{
				{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob"}},
				{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{"Alice"}},
			},
			expectedValid:  false,
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body, _ := json.Marshal(tt.participants)
			req := httptest.NewRequest(http.MethodPost, "/api/validate", bytes.NewReader(body))
			w := httptest.NewRecorder()

			server.HandleValidate(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			var response ValidationResponse
			json.NewDecoder(w.Body).Decode(&response)

			if response.Valid != tt.expectedValid {
				t.Errorf("Expected valid=%v, got valid=%v", tt.expectedValid, response.Valid)
			}
		})
	}
}

func TestHandleDraw(t *testing.T) {
	server := NewServer(":8080")

	drawRequest := DrawRequest{
		Participants: []participant.Participant{
			{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
			{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
			{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
		},
	}

	body, _ := json.Marshal(drawRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/draw", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.HandleDraw(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}

	var response DrawResponse
	json.NewDecoder(w.Body).Decode(&response)

	if !response.Success {
		t.Errorf("Expected success=true, got false: %s", response.Error)
	}

	if len(response.Participants) != 3 {
		t.Errorf("Expected 3 participants, got %d", len(response.Participants))
	}

	// Verify all participants have recipients
	for _, p := range response.Participants {
		if p.Recipient == nil {
			t.Errorf("Participant %s has no recipient", p.Name)
		}
	}
}

func TestHandleDrawInvalid(t *testing.T) {
	server := NewServer(":8080")

	// Invalid: only 1 participant
	drawRequest := DrawRequest{
		Participants: []participant.Participant{
			{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
		},
	}

	body, _ := json.Marshal(drawRequest)
	req := httptest.NewRequest(http.MethodPost, "/api/draw", bytes.NewReader(body))
	w := httptest.NewRecorder()

	server.HandleDraw(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected status 400, got %d", w.Code)
	}

	var response DrawResponse
	json.NewDecoder(w.Body).Decode(&response)

	if response.Success {
		t.Error("Expected success=false for invalid configuration")
	}
}

func TestMethodNotAllowed(t *testing.T) {
	server := NewServer(":8080")

	endpoints := []string{"/api/validate", "/api/draw", "/api/upload"}

	for _, endpoint := range endpoints {
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		w := httptest.NewRecorder()

		switch endpoint {
		case "/api/validate":
			server.HandleValidate(w, req)
		case "/api/draw":
			server.HandleDraw(w, req)
		case "/api/upload":
			server.HandleUpload(w, req)
		}

		if w.Code != http.StatusMethodNotAllowed {
			t.Errorf("Expected 405 for GET %s, got %d", endpoint, w.Code)
		}
	}
}
