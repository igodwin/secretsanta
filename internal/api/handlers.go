package api

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/igodwin/secretsanta/internal/draw"
	"github.com/igodwin/secretsanta/internal/formats"
	"github.com/igodwin/secretsanta/pkg/config"
	"github.com/igodwin/secretsanta/pkg/participant"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/igodwin/secretsanta/api/grpc/pb"
)

type Server struct {
	addr string
}

func NewServer(addr string) *Server {
	return &Server{addr: addr}
}

// Participant input/output models
type ParticipantRequest struct {
	Name             string   `json:"name"`
	NotificationType string   `json:"notification_type"`
	ContactInfo      []string `json:"contact_info"`
	Exclusions       []string `json:"exclusions"`
}

type ParticipantResponse struct {
	Name             string   `json:"name"`
	NotificationType string   `json:"notification_type"`
	ContactInfo      []string `json:"contact_info"`
	Exclusions       []string `json:"exclusions"`
	Recipient        *string  `json:"recipient,omitempty"`
}

type ValidationResponse struct {
	Valid                     bool     `json:"valid"`
	Errors                    []string `json:"errors,omitempty"`
	Warnings                  []string `json:"warnings,omitempty"`
	ParticipantsWithNoOptions []string `json:"participants_with_no_options,omitempty"`
	MinCompatibility          int      `json:"min_compatibility"`
	AvgCompatibility          float64  `json:"avg_compatibility"`
	TotalParticipants         int      `json:"total_participants"`
}

type DrawRequest struct {
	Participants []participant.Participant `json:"participants"`
	ArchiveEmail string                    `json:"archive_email,omitempty"`
}

type DrawResponse struct {
	Success      bool                   `json:"success"`
	Participants []*ParticipantResponse `json:"participants,omitempty"`
	Error        string                 `json:"error,omitempty"`
}

// HandleValidate validates participant data without performing draw
func (s *Server) HandleValidate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var participants []*participant.Participant
	if err := json.NewDecoder(r.Body).Decode(&participants); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	result := draw.ValidateParticipants(participants)

	response := ValidationResponse{
		Valid:                     result.IsValid,
		Errors:                    result.Errors,
		Warnings:                  result.Warnings,
		ParticipantsWithNoOptions: result.ParticipantsWithNoOptions,
		MinCompatibility:          result.MinCompatibility,
		AvgCompatibility:          result.AvgCompatibility,
		TotalParticipants:         result.TotalParticipants,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleDraw performs the Secret Santa draw
func (s *Server) HandleDraw(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var drawRequest DrawRequest
	if err := json.NewDecoder(r.Body).Decode(&drawRequest); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	// Convert to pointers for internal use
	participants := make([]*participant.Participant, len(drawRequest.Participants))
	for i := range drawRequest.Participants {
		participants[i] = &drawRequest.Participants[i]
	}

	// Validate first
	validation := draw.ValidateParticipants(participants)
	if !validation.IsValid {
		response := DrawResponse{
			Success: false,
			Error:   fmt.Sprintf("Validation failed: %v", validation.Errors),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(response)
		return
	}

	// Perform draw
	result, err := draw.Names(participants)
	if err != nil {
		response := DrawResponse{
			Success: false,
			Error:   err.Error(),
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(response)
		return
	}

	// If archive email is provided and notifier integration is available, use it
	if drawRequest.ArchiveEmail != "" {
		log.Printf("Draw completed with archive email: %s", drawRequest.ArchiveEmail)
		// Archive email will be used by notification system
		// Store it for future notification sending
	}

	// Convert to response format
	participantResponses := make([]*ParticipantResponse, len(result))
	for i, p := range result {
		var recipientName *string
		if p.Recipient != nil {
			recipientName = &p.Recipient.Name
		}

		participantResponses[i] = &ParticipantResponse{
			Name:             p.Name,
			NotificationType: p.NotificationType,
			ContactInfo:      p.ContactInfo,
			Exclusions:       p.Exclusions,
			Recipient:        recipientName,
		}
	}

	response := DrawResponse{
		Success:      true,
		Participants: participantResponses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleUpload handles file upload for participant data
// Supports JSON, YAML, TOML, CSV, and TSV formats
func (s *Server) HandleUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse multipart form (10MB max)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "File too large or invalid", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Detect file format from extension
	format, err := formats.DetectFormat(header.Filename)
	if err != nil {
		http.Error(w, fmt.Sprintf("Unsupported file format: %v", err), http.StatusBadRequest)
		return
	}

	// Read file content
	data, err := io.ReadAll(file)
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Parse file using appropriate parser
	participants, err := formats.Parse(data, format)
	if err != nil {
		http.Error(w, fmt.Sprintf("Invalid file format: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Uploaded file: %s (format: %s, %d participants)", header.Filename, format, len(participants))

	// Validate uploaded data
	validation := draw.ValidateParticipants(participants)

	response := struct {
		Success      bool                       `json:"success"`
		Participants []*participant.Participant `json:"participants"`
		Validation   ValidationResponse         `json:"validation"`
		Format       string                     `json:"format"`
	}{
		Success:      true,
		Participants: participants,
		Format:       string(format),
		Validation: ValidationResponse{
			Valid:                     validation.IsValid,
			Errors:                    validation.Errors,
			Warnings:                  validation.Warnings,
			ParticipantsWithNoOptions: validation.ParticipantsWithNoOptions,
			MinCompatibility:          validation.MinCompatibility,
			AvgCompatibility:          validation.AvgCompatibility,
			TotalParticipants:         validation.TotalParticipants,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// HandleExport exports participant data with assignments as JSON
func (s *Server) HandleExport(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var participants []*ParticipantResponse
	if err := json.NewDecoder(r.Body).Decode(&participants); err != nil {
		http.Error(w, fmt.Sprintf("Invalid JSON: %v", err), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", "attachment; filename=secretsanta-results.json")
	json.NewEncoder(w).Encode(participants)
}

// HandleTemplate generates and downloads a template file for the specified format
func (s *Server) HandleTemplate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get format from query parameter
	formatStr := r.URL.Query().Get("format")
	if formatStr == "" {
		http.Error(w, "Missing format parameter", http.StatusBadRequest)
		return
	}

	format := formats.FileFormat(formatStr)

	// Generate template
	data, mimeType, err := formats.GenerateTemplate(format)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate template: %v", err), http.StatusInternalServerError)
		return
	}

	// Set headers
	filename := formats.GetFilename(format)
	w.Header().Set("Content-Type", mimeType)
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
	w.Write(data)
}

// NotificationStatusResponse contains information about available notification types
type NotificationStatusResponse struct {
	Available        []string          `json:"available"`
	UsingNotifier    bool              `json:"using_notifier"`
	NotifierHealthy  bool              `json:"notifier_healthy,omitempty"`
	NotifierStatus   string            `json:"notifier_status,omitempty"`
	NotifierDetails  map[string]string `json:"notifier_details,omitempty"`
	SMTPConfigured   bool              `json:"smtp_configured"`
}

// HandleStatus returns the current notification configuration status
func (s *Server) HandleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cfg := config.GetConfig()
	response := NotificationStatusResponse{
		Available: []string{"stdout"}, // stdout always available
	}

	// Check if SMTP is configured
	if cfg.SMTP.Host != "" && cfg.SMTP.FromAddress != "" {
		response.Available = append(response.Available, "email")
		response.SMTPConfigured = true
	}

	// Check if external notifier service is configured
	if cfg.Notifier.ServiceAddr != "" {
		response.UsingNotifier = true

		// Try to query the notifier service for available types
		notifierTypes, healthy, status, details := checkNotifierHealth(cfg.Notifier.ServiceAddr)
		response.NotifierHealthy = healthy
		response.NotifierStatus = status
		response.NotifierDetails = details

		if healthy && len(notifierTypes) > 0 {
			// Use types from notifier service
			response.Available = notifierTypes
		} else if healthy {
			// Notifier is healthy but didn't report types, assume standard types
			response.Available = []string{"email", "slack", "stdout"}
		}
		// If notifier is not healthy, keep the basic types we detected
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// checkNotifierHealth queries the notifier service health endpoint
func checkNotifierHealth(serviceAddr string) ([]string, bool, string, map[string]string) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, serviceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock())
	if err != nil {
		return nil, false, "unreachable", nil
	}
	defer conn.Close()

	client := pb.NewNotifierServiceClient(conn)

	healthResp, err := client.HealthCheck(ctx, &pb.HealthCheckRequest{})
	if err != nil {
		return nil, false, "error", nil
	}

	// Extract available notification types from components
	var types []string
	if healthResp.Components != nil {
		// Check for configured providers in the components map
		// Common keys might be "email", "slack", "ntfy", etc.
		for key := range healthResp.Components {
			switch key {
			case "email", "smtp":
				types = appendUnique(types, "email")
			case "slack":
				types = appendUnique(types, "slack")
			case "ntfy":
				types = appendUnique(types, "ntfy")
			}
		}
	}

	// Always add stdout as it's always available
	types = appendUnique(types, "stdout")

	return types, healthResp.Healthy, healthResp.Status, healthResp.Components
}

// appendUnique appends a string to a slice only if it's not already present
func appendUnique(slice []string, item string) []string {
	for _, existing := range slice {
		if existing == item {
			return slice
		}
	}
	return append(slice, item)
}

// Start starts the HTTP server
func (s *Server) Start() error {
	mux := http.NewServeMux()

	// API endpoints
	mux.HandleFunc("/api/validate", s.HandleValidate)
	mux.HandleFunc("/api/draw", s.HandleDraw)
	mux.HandleFunc("/api/upload", s.HandleUpload)
	mux.HandleFunc("/api/export", s.HandleExport)
	mux.HandleFunc("/api/template", s.HandleTemplate)
	mux.HandleFunc("/api/status", s.HandleStatus)

	// Static files
	fs := http.FileServer(http.Dir("internal/web/static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html for root
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		http.ServeFile(w, r, "internal/web/static/index.html")
	})

	log.Printf("Starting web server on %s", s.addr)
	return http.ListenAndServe(s.addr, corsMiddleware(mux))
}

// corsMiddleware adds CORS headers for development
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}
