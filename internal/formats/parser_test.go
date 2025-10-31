package formats

import (
	"testing"

	"github.com/igodwin/secretsanta/pkg/participant"
)

func TestParseJSON(t *testing.T) {
	data := []byte(`[
  {
    "name": "Alice",
    "notification_type": "email",
    "contact_info": ["alice@example.com"],
    "exclusions": ["Bob"]
  },
  {
    "name": "Bob",
    "notification_type": "email",
    "contact_info": ["bob@example.com"],
    "exclusions": ["Alice"]
  }
]`)

	participants, err := Parse(data, FormatJSON)
	if err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	if participants[0].Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", participants[0].Name)
	}
}

func TestParseYAML(t *testing.T) {
	data := []byte(`- name: Alice
  notification_type: email
  contact_info:
    - alice@example.com
  exclusions:
    - Bob
- name: Bob
  notification_type: email
  contact_info:
    - bob@example.com
  exclusions:
    - Alice`)

	participants, err := Parse(data, FormatYAML)
	if err != nil {
		t.Fatalf("Failed to parse YAML: %v", err)
	}

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	if participants[0].Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", participants[0].Name)
	}
}

func TestParseCSV(t *testing.T) {
	data := []byte(`name,notification_type,contact_info,exclusions
Alice,email,alice@example.com,Bob
Bob,email,bob@example.com,Alice
Carol,slack,@carol,`)

	participants, err := Parse(data, FormatCSV)
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	if len(participants) != 3 {
		t.Errorf("Expected 3 participants, got %d", len(participants))
	}

	if participants[0].Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", participants[0].Name)
	}

	if len(participants[0].ContactInfo) != 1 || participants[0].ContactInfo[0] != "alice@example.com" {
		t.Errorf("Expected contact info 'alice@example.com', got %v", participants[0].ContactInfo)
	}

	if len(participants[2].Exclusions) != 0 {
		t.Errorf("Expected Carol to have no exclusions, got %v", participants[2].Exclusions)
	}
}

func TestParseTSV(t *testing.T) {
	data := []byte("name\tnotification_type\tcontact_info\texclusions\nAlice\temail\talice@example.com\tBob\nBob\temail\tbob@example.com\tAlice")

	participants, err := Parse(data, FormatTSV)
	if err != nil {
		t.Fatalf("Failed to parse TSV: %v", err)
	}

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	if participants[0].Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", participants[0].Name)
	}
}

func TestParseTOML(t *testing.T) {
	data := []byte(`[[participants]]
name = "Alice"
notification_type = "email"
contact_info = ["alice@example.com"]
exclusions = ["Bob"]

[[participants]]
name = "Bob"
notification_type = "email"
contact_info = ["bob@example.com"]
exclusions = ["Alice"]`)

	participants, err := Parse(data, FormatTOML)
	if err != nil {
		t.Fatalf("Failed to parse TOML: %v", err)
	}

	if len(participants) != 2 {
		t.Errorf("Expected 2 participants, got %d", len(participants))
	}

	if participants[0].Name != "Alice" {
		t.Errorf("Expected name 'Alice', got '%s'", participants[0].Name)
	}
}

func TestDetectFormat(t *testing.T) {
	tests := []struct {
		filename string
		expected FileFormat
		wantErr  bool
	}{
		{"test.json", FormatJSON, false},
		{"test.yaml", FormatYAML, false},
		{"test.yml", FormatYAML, false},
		{"test.toml", FormatTOML, false},
		{"test.csv", FormatCSV, false},
		{"test.tsv", FormatTSV, false},
		{"test.txt", "", true},
		{"test", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			format, err := DetectFormat(tt.filename)
			if (err != nil) != tt.wantErr {
				t.Errorf("DetectFormat() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if format != tt.expected {
				t.Errorf("DetectFormat() = %v, want %v", format, tt.expected)
			}
		})
	}
}

func TestCSVMultipleValues(t *testing.T) {
	data := []byte(`name,notification_type,contact_info,exclusions
Alice,email,alice@example.com; alice.alt@example.com,Bob; Carol`)

	participants, err := Parse(data, FormatCSV)
	if err != nil {
		t.Fatalf("Failed to parse CSV: %v", err)
	}

	if len(participants) != 1 {
		t.Errorf("Expected 1 participant, got %d", len(participants))
	}

	if len(participants[0].ContactInfo) != 2 {
		t.Errorf("Expected 2 contact info entries, got %d", len(participants[0].ContactInfo))
	}

	if len(participants[0].Exclusions) != 2 {
		t.Errorf("Expected 2 exclusions, got %d", len(participants[0].Exclusions))
	}
}

// Test export functions
func TestExportJSON(t *testing.T) {
	testParticipants := []*participant.Participant{
		{
			Name:             "Alice",
			NotificationType: "email",
			ContactInfo:      []string{"alice@example.com"},
			Exclusions:       []string{"Bob"},
		},
		{
			Name:             "Bob",
			NotificationType: "email",
			ContactInfo:      []string{"bob@example.com"},
			Exclusions:       []string{"Alice"},
		},
	}

	data, mimeType, err := ExportParticipants(testParticipants, FormatJSON)
	if err != nil {
		t.Fatalf("Failed to export JSON: %v", err)
	}

	if mimeType != "application/json" {
		t.Errorf("Expected MIME type 'application/json', got '%s'", mimeType)
	}

	// Verify we can parse it back
	parsed, err := Parse(data, FormatJSON)
	if err != nil {
		t.Fatalf("Failed to parse exported JSON: %v", err)
	}

	if len(parsed) != 2 {
		t.Errorf("Expected 2 participants after round-trip, got %d", len(parsed))
	}
}

// Test export functions
func TestExportYAML(t *testing.T) {
	testParticipants := []*participant.Participant{
		{
			Name:             "Alice",
			NotificationType: "email",
			ContactInfo:      []string{"alice@example.com"},
			Exclusions:       []string{"Bob"},
		},
	}

	data, mimeType, err := ExportParticipants(testParticipants, FormatYAML)
	if err != nil {
		t.Fatalf("Failed to export YAML: %v", err)
	}

	if mimeType != "application/x-yaml" {
		t.Errorf("Expected MIME type 'application/x-yaml', got '%s'", mimeType)
	}

	// Verify we can parse it back
	parsed, err := Parse(data, FormatYAML)
	if err != nil {
		t.Fatalf("Failed to parse exported YAML: %v", err)
	}

	if len(parsed) != 1 || parsed[0].Name != "Alice" {
		t.Errorf("YAML round-trip failed")
	}
}

// Test export functions
func TestExportCSV(t *testing.T) {
	testParticipants := []*participant.Participant{
		{
			Name:             "Alice",
			NotificationType: "email",
			ContactInfo:      []string{"alice@example.com"},
			Exclusions:       []string{"Bob"},
		},
	}

	data, mimeType, err := ExportParticipants(testParticipants, FormatCSV)
	if err != nil {
		t.Fatalf("Failed to export CSV: %v", err)
	}

	if mimeType != "text/csv" {
		t.Errorf("Expected MIME type 'text/csv', got '%s'", mimeType)
	}

	// Verify we can parse it back
	parsed, err := Parse(data, FormatCSV)
	if err != nil {
		t.Fatalf("Failed to parse exported CSV: %v", err)
	}

	if len(parsed) != 1 || parsed[0].Name != "Alice" {
		t.Errorf("CSV round-trip failed")
	}
}

// Test export functions
func TestExportTSV(t *testing.T) {
	testParticipants := []*participant.Participant{
		{
			Name:             "Alice",
			NotificationType: "email",
			ContactInfo:      []string{"alice@example.com"},
			Exclusions:       []string{"Bob"},
		},
	}

	data, mimeType, err := ExportParticipants(testParticipants, FormatTSV)
	if err != nil {
		t.Fatalf("Failed to export TSV: %v", err)
	}

	if mimeType != "text/tab-separated-values" {
		t.Errorf("Expected MIME type 'text/tab-separated-values', got '%s'", mimeType)
	}

	// Verify we can parse it back
	parsed, err := Parse(data, FormatTSV)
	if err != nil {
		t.Fatalf("Failed to parse exported TSV: %v", err)
	}

	if len(parsed) != 1 || parsed[0].Name != "Alice" {
		t.Errorf("TSV round-trip failed")
	}
}

// Test export functions
func TestExportTOML(t *testing.T) {
	testParticipants := []*participant.Participant{
		{
			Name:             "Alice",
			NotificationType: "email",
			ContactInfo:      []string{"alice@example.com"},
			Exclusions:       []string{"Bob"},
		},
	}

	data, mimeType, err := ExportParticipants(testParticipants, FormatTOML)
	if err != nil {
		t.Fatalf("Failed to export TOML: %v", err)
	}

	if mimeType != "application/toml" {
		t.Errorf("Expected MIME type 'application/toml', got '%s'", mimeType)
	}

	// Verify we can parse it back
	parsed, err := Parse(data, FormatTOML)
	if err != nil {
		t.Fatalf("Failed to parse exported TOML: %v", err)
	}

	if len(parsed) != 1 || parsed[0].Name != "Alice" {
		t.Errorf("TOML round-trip failed")
	}
}
