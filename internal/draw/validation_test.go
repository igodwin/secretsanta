package draw

import (
	"testing"

	"github.com/igodwin/secretsanta/pkg/participant"
)

func TestValidateParticipants_Valid(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	if !result.IsValid {
		t.Errorf("Expected valid, got invalid: %v", result.Errors)
	}

	if len(result.Errors) > 0 {
		t.Errorf("Expected no errors, got: %v", result.Errors)
	}

	if result.TotalParticipants != 3 {
		t.Errorf("Expected 3 participants, got %d", result.TotalParticipants)
	}

	if result.MinCompatibility != 2 {
		t.Errorf("Expected min compatibility 2, got %d", result.MinCompatibility)
	}
}

func TestValidateParticipants_NoParticipants(t *testing.T) {
	participants := []*participant.Participant{}

	result := ValidateParticipants(participants)

	if result.IsValid {
		t.Error("Expected invalid for empty participant list")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors for empty participant list")
	}
}

func TestValidateParticipants_TooFewParticipants(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	if result.IsValid {
		t.Error("Expected invalid for single participant")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors for single participant")
	}
}

func TestValidateParticipants_DuplicateNames(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Alice", ContactInfo: []string{"alice2@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	if result.IsValid {
		t.Error("Expected invalid for duplicate names")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors for duplicate names")
	}
}

func TestValidateParticipants_NoValidRecipients(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob", "Carol"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	if result.IsValid {
		t.Error("Expected invalid when participant has no valid recipients")
	}

	if len(result.ParticipantsWithNoOptions) != 1 {
		t.Errorf("Expected 1 participant with no options, got %d", len(result.ParticipantsWithNoOptions))
	}

	if result.ParticipantsWithNoOptions[0] != "Alice" {
		t.Errorf("Expected Alice to have no options, got %s", result.ParticipantsWithNoOptions[0])
	}
}

func TestValidateParticipants_MissingContactInfo(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{}, Exclusions: []string{}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	// Should still be valid (contact info is a warning, not an error)
	if !result.IsValid {
		t.Errorf("Expected valid despite missing contact info: %v", result.Errors)
	}

	// Should have a warning
	if len(result.Warnings) == 0 {
		t.Error("Expected warning for missing contact info")
	}
}

func TestValidateParticipants_InvalidExclusion(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"NonExistent"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	// Should still be valid (invalid exclusion is a warning)
	if !result.IsValid {
		t.Errorf("Expected valid despite invalid exclusion: %v", result.Errors)
	}

	// Should have a warning
	if len(result.Warnings) == 0 {
		t.Error("Expected warning for invalid exclusion")
	}
}

func TestValidateParticipants_LowCompatibility(t *testing.T) {
	// This configuration has low compatibility but IS valid:
	// Alice can give to David or Eve
	// Bob can give to Carol, David, or Eve
	// Carol can give to Alice, David, or Eve
	// David can give to anyone (Alice, Bob, Carol, or Eve)
	// Eve can give to anyone (Alice, Bob, Carol, or David)
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob", "Carol"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{"Alice"}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{"Bob"}},
		{Name: "David", ContactInfo: []string{"david@example.com"}, Exclusions: []string{}},
		{Name: "Eve", ContactInfo: []string{"eve@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	if !result.IsValid {
		t.Errorf("Expected valid: %v", result.Errors)
	}

	if result.MinCompatibility < 2 {
		t.Errorf("Expected min compatibility >= 2, got %d", result.MinCompatibility)
	}

	// This configuration should pass validation
	t.Logf("Configuration valid with min_compatibility=%d, avg_compatibility=%.1f",
		result.MinCompatibility, result.AvgCompatibility)
}

func TestValidateParticipantsQuick_Valid(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
	}

	if !ValidateParticipantsQuick(participants) {
		t.Error("Expected quick validation to pass")
	}
}

func TestValidateParticipantsQuick_Invalid(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob", "Carol"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
	}

	if ValidateParticipantsQuick(participants) {
		t.Error("Expected quick validation to fail")
	}
}

func TestValidateParticipants_ComplexScenario(t *testing.T) {
	// Real-world scenario: couple exclusions
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob"}}, // Alice's partner
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{"Alice"}},   // Bob's partner
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{"David"}},
		{Name: "David", ContactInfo: []string{"david@example.com"}, Exclusions: []string{"Carol"}},
		{Name: "Eve", ContactInfo: []string{"eve@example.com"}, Exclusions: []string{}},
		{Name: "Frank", ContactInfo: []string{"frank@example.com"}, Exclusions: []string{}},
	}

	result := ValidateParticipants(participants)

	if !result.IsValid {
		t.Errorf("Expected valid for couples scenario: %v", result.Errors)
	}

	if result.TotalParticipants != 6 {
		t.Errorf("Expected 6 participants, got %d", result.TotalParticipants)
	}

	// Each person in a couple has 4 valid recipients (everyone except self and partner)
	if result.MinCompatibility != 4 {
		t.Errorf("Expected min compatibility 4, got %d", result.MinCompatibility)
	}
}

func TestValidateParticipants_HallsTheoremViolation(t *testing.T) {
	// This is the bug case: Emily and Ivan both exclude each other,
	// so they both can only give to Eli, but Eli can only receive from one
	participants := []*participant.Participant{
		{Name: "Emily", ContactInfo: []string{"emily@example.com"}, Exclusions: []string{"Ivan"}},
		{Name: "Eli", ContactInfo: []string{"eli@example.com"}, Exclusions: []string{}},
		{Name: "Ivan", ContactInfo: []string{"ivan@example.com"}, Exclusions: []string{"Emily"}},
	}

	result := ValidateParticipants(participants)

	if result.IsValid {
		t.Error("Expected invalid configuration (Hall's Marriage Theorem violation)")
	}

	if len(result.Errors) == 0 {
		t.Error("Expected errors explaining why configuration is invalid")
	}

	t.Logf("Validation correctly caught the issue: %v", result.Errors)
}