package formats

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"

	"github.com/igodwin/secretsanta/pkg/participant"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// GenerateTemplate creates a template file for the given format
func GenerateTemplate(format FileFormat) ([]byte, string, error) {
	// Sample participants for template
	participants := getSampleParticipants()

	var (
		data     []byte
		mimeType string
		err      error
	)

	switch format {
	case FormatJSON:
		data, err = generateJSON(participants)
		mimeType = "application/json"
	case FormatYAML:
		data, err = generateYAML(participants)
		mimeType = "application/x-yaml"
	case FormatTOML:
		data, err = generateTOML(participants)
		mimeType = "application/toml"
	case FormatCSV:
		data, err = generateCSV(participants, ',')
		mimeType = "text/csv"
	case FormatTSV:
		data, err = generateCSV(participants, '\t')
		mimeType = "text/tab-separated-values"
	default:
		return nil, "", fmt.Errorf("unsupported format: %s", format)
	}

	if err != nil {
		return nil, "", err
	}

	return data, mimeType, nil
}

// getSampleParticipants returns sample participants for templates
func getSampleParticipants() []*participant.Participant {
	return []*participant.Participant{
		{
			Name:             "Alice Johnson",
			NotificationType: "email",
			ContactInfo:      []string{"alice@example.com"},
			Exclusions:       []string{"Bob Smith"},
		},
		{
			Name:             "Bob Smith",
			NotificationType: "email",
			ContactInfo:      []string{"bob@example.com"},
			Exclusions:       []string{"Alice Johnson"},
		},
		{
			Name:             "Carol Davis",
			NotificationType: "slack",
			ContactInfo:      []string{"@carol"},
			Exclusions:       []string{},
		},
		{
			Name:             "David Wilson",
			NotificationType: "email",
			ContactInfo:      []string{"david@example.com", "david.alt@example.com"},
			Exclusions:       []string{"Carol Davis"},
		},
	}
}

// generateJSON creates a JSON template
func generateJSON(participants []*participant.Participant) ([]byte, error) {
	data, err := json.MarshalIndent(participants, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to generate JSON: %w", err)
	}
	return data, nil
}

// generateYAML creates a YAML template
func generateYAML(participants []*participant.Participant) ([]byte, error) {
	data, err := yaml.Marshal(participants)
	if err != nil {
		return nil, fmt.Errorf("failed to generate YAML: %w", err)
	}
	return data, nil
}

// generateTOML creates a TOML template
func generateTOML(participants []*participant.Participant) ([]byte, error) {
	wrapper := struct {
		Participants []*participant.Participant `toml:"participants"`
	}{
		Participants: participants,
	}
	data, err := toml.Marshal(wrapper)
	if err != nil {
		return nil, fmt.Errorf("failed to generate TOML: %w", err)
	}
	return data, nil
}

// generateCSV creates a CSV/TSV template
func generateCSV(participants []*participant.Participant, delimiter rune) ([]byte, error) {
	buf := new(bytes.Buffer)
	writer := csv.NewWriter(buf)
	writer.Comma = delimiter

	// Write header
	header := []string{"name", "notification_type", "contact_info", "exclusions"}
	if err := writer.Write(header); err != nil {
		return nil, fmt.Errorf("failed to write CSV header: %w", err)
	}

	// Write data
	for _, p := range participants {
		record := []string{
			p.Name,
			p.NotificationType,
			joinList(p.ContactInfo),
			joinList(p.Exclusions),
		}
		if err := writer.Write(record); err != nil {
			return nil, fmt.Errorf("failed to write CSV record: %w", err)
		}
	}

	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, fmt.Errorf("CSV writer error: %w", err)
	}

	return buf.Bytes(), nil
}

// joinList joins a string slice with semicolons for CSV cells
func joinList(items []string) string {
	if len(items) == 0 {
		return ""
	}
	// Use semicolon to avoid conflict with CSV comma delimiter
	result := ""
	for i, item := range items {
		if i > 0 {
			result += "; "
		}
		result += item
	}
	return result
}

// GetFilename returns the appropriate filename for a format
func GetFilename(format FileFormat) string {
	switch format {
	case FormatJSON:
		return "secretsanta-template.json"
	case FormatYAML:
		return "secretsanta-template.yaml"
	case FormatTOML:
		return "secretsanta-template.toml"
	case FormatCSV:
		return "secretsanta-template.csv"
	case FormatTSV:
		return "secretsanta-template.tsv"
	default:
		return "secretsanta-template.txt"
	}
}
