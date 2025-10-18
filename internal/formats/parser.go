package formats

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"strings"

	"github.com/igodwin/secretsanta/pkg/participant"
	"github.com/pelletier/go-toml/v2"
	"gopkg.in/yaml.v3"
)

// FileFormat represents the type of file being parsed
type FileFormat string

const (
	FormatJSON FileFormat = "json"
	FormatYAML FileFormat = "yaml"
	FormatTOML FileFormat = "toml"
	FormatCSV  FileFormat = "csv"
	FormatTSV  FileFormat = "tsv"
)

// DetectFormat determines the file format based on extension
func DetectFormat(filename string) (FileFormat, error) {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".json":
		return FormatJSON, nil
	case ".yaml", ".yml":
		return FormatYAML, nil
	case ".toml":
		return FormatTOML, nil
	case ".csv":
		return FormatCSV, nil
	case ".tsv":
		return FormatTSV, nil
	default:
		return "", fmt.Errorf("unsupported file format: %s", ext)
	}
}

// Parse parses participant data from various file formats
func Parse(data []byte, format FileFormat) ([]*participant.Participant, error) {
	switch format {
	case FormatJSON:
		return parseJSON(data)
	case FormatYAML:
		return parseYAML(data)
	case FormatTOML:
		return parseTOML(data)
	case FormatCSV:
		return parseCSV(data, ',')
	case FormatTSV:
		return parseCSV(data, '\t')
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

// parseJSON parses JSON format
func parseJSON(data []byte) ([]*participant.Participant, error) {
	var participants []*participant.Participant
	if err := json.Unmarshal(data, &participants); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return participants, nil
}

// parseYAML parses YAML format
func parseYAML(data []byte) ([]*participant.Participant, error) {
	var participants []*participant.Participant
	if err := yaml.Unmarshal(data, &participants); err != nil {
		return nil, fmt.Errorf("invalid YAML: %w", err)
	}
	return participants, nil
}

// parseTOML parses TOML format
func parseTOML(data []byte) ([]*participant.Participant, error) {
	var wrapper struct {
		Participants []*participant.Participant `toml:"participants"`
	}
	if err := toml.Unmarshal(data, &wrapper); err != nil {
		return nil, fmt.Errorf("invalid TOML: %w", err)
	}
	return wrapper.Participants, nil
}

// parseCSV parses CSV/TSV format
// Expected columns: Name, NotificationType, ContactInfo, Exclusions
func parseCSV(data []byte, delimiter rune) ([]*participant.Participant, error) {
	reader := csv.NewReader(strings.NewReader(string(data)))
	reader.Comma = delimiter
	reader.TrimLeadingSpace = true

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("failed to read header: %w", err)
	}

	// Validate header
	if len(header) < 3 {
		return nil, fmt.Errorf("invalid CSV format: expected at least 3 columns (name, notification_type, contact_info)")
	}

	// Normalize headers
	for i := range header {
		header[i] = strings.ToLower(strings.TrimSpace(header[i]))
	}

	var participants []*participant.Participant
	lineNum := 1 // Start from 1 since we already read the header

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading line %d: %w", lineNum+1, err)
		}
		lineNum++

		if len(record) < 3 {
			return nil, fmt.Errorf("line %d: insufficient columns", lineNum)
		}

		name := strings.TrimSpace(record[0])
		if name == "" {
			continue // Skip empty rows
		}

		notificationType := strings.TrimSpace(record[1])
		if notificationType == "" {
			notificationType = "email" // Default
		}

		// Parse contact info (comma-separated within the cell)
		contactInfoStr := strings.TrimSpace(record[2])
		contactInfo := parseListField(contactInfoStr)

		// Parse exclusions if present
		var exclusions []string
		if len(record) > 3 && strings.TrimSpace(record[3]) != "" {
			exclusions = parseListField(strings.TrimSpace(record[3]))
		}

		participants = append(participants, &participant.Participant{
			Name:             name,
			NotificationType: notificationType,
			ContactInfo:      contactInfo,
			Exclusions:       exclusions,
		})
	}

	return participants, nil
}

// parseListField parses a comma or semicolon-separated list
func parseListField(s string) []string {
	if s == "" {
		return nil
	}

	// Support both comma and semicolon as separators
	var items []string
	separator := ","
	if strings.Contains(s, ";") && !strings.Contains(s, ",") {
		separator = ";"
	}

	parts := strings.Split(s, separator)
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			items = append(items, trimmed)
		}
	}
	return items
}
