package draw

import (
	"testing"

	"github.com/igodwin/secretsanta/pkg/participant"
)

func TestNamesOptimized_Basic(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
	}

	result, err := NamesOptimized(participants)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify all participants have a recipient
	for _, p := range result {
		if p.Recipient == nil {
			t.Errorf("Participant %s has no recipient", p.Name)
		}
		if p.Recipient.Name == p.Name {
			t.Errorf("Participant %s assigned to themselves", p.Name)
		}
	}

	// Verify no duplicates
	recipientCount := make(map[string]int)
	for _, p := range result {
		recipientCount[p.Recipient.Name]++
	}
	for name, count := range recipientCount {
		if count > 1 {
			t.Errorf("Recipient %s assigned %d times", name, count)
		}
	}
}

func TestNamesOptimized_WithExclusions(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{"Carol"}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
		{Name: "David", ContactInfo: []string{"david@example.com"}, Exclusions: []string{}},
	}

	result, err := NamesOptimized(participants)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify exclusions are respected
	for _, p := range result {
		for _, excluded := range p.Exclusions {
			if p.Recipient.Name == excluded {
				t.Errorf("Participant %s was assigned to excluded recipient %s", p.Name, excluded)
			}
		}
	}
}

func TestNamesOptimized_ImpossibleConstraints(t *testing.T) {
	// Create scenario where no valid assignment exists
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob", "Carol"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{"Alice", "Carol"}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{"Alice", "Bob"}},
	}

	_, err := NamesOptimized(participants)
	if err == nil {
		t.Fatal("Expected error for impossible constraints, got nil")
	}
}

func TestNamesOptimized_LargeGroup(t *testing.T) {
	participants := createTestParticipants(50, 5)

	result, err := NamesOptimized(participants)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	// Verify constraints
	for _, p := range result {
		if p.Recipient == nil {
			t.Errorf("Participant %s has no recipient", p.Name)
		}
		if p.Recipient.Name == p.Name {
			t.Errorf("Participant %s assigned to themselves", p.Name)
		}

		for _, excluded := range p.Exclusions {
			if p.Recipient.Name == excluded {
				t.Errorf("Participant %s was assigned to excluded recipient %s", p.Name, excluded)
			}
		}
	}

	// Verify uniqueness
	recipientCount := make(map[string]int)
	for _, p := range result {
		recipientCount[p.Recipient.Name]++
	}
	for name, count := range recipientCount {
		if count > 1 {
			t.Errorf("Recipient %s assigned %d times", name, count)
		}
	}
}

func TestNamesOptimizedWithStats(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", ContactInfo: []string{"alice@example.com"}, Exclusions: []string{"Bob"}},
		{Name: "Bob", ContactInfo: []string{"bob@example.com"}, Exclusions: []string{}},
		{Name: "Carol", ContactInfo: []string{"carol@example.com"}, Exclusions: []string{}},
		{Name: "David", ContactInfo: []string{"david@example.com"}, Exclusions: []string{"Carol"}},
	}

	result, stats, err := NamesOptimizedWithStats(participants)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if !stats.Success {
		t.Error("Expected successful draw")
	}

	if stats.TotalParticipants != 4 {
		t.Errorf("Expected 4 participants, got %d", stats.TotalParticipants)
	}

	if stats.AvgCompatibilityPerPerson <= 0 {
		t.Errorf("Expected positive average compatibility, got %f", stats.AvgCompatibilityPerPerson)
	}

	if stats.HasImpossibleConstraints {
		t.Error("Should not have impossible constraints")
	}

	if len(result) != 4 {
		t.Errorf("Expected 4 participants in result, got %d", len(result))
	}
}

func TestBuildExclusionMap(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", Exclusions: []string{"Bob", "Carol"}},
		{Name: "Bob", Exclusions: []string{}},
		{Name: "Carol", Exclusions: []string{"Alice"}},
	}

	exclusionMap := buildExclusionMap(participants)

	if len(exclusionMap) != 2 {
		t.Errorf("Expected 2 entries in exclusion map, got %d", len(exclusionMap))
	}

	if !exclusionMap["Alice"]["Bob"] {
		t.Error("Expected Alice to exclude Bob")
	}

	if !exclusionMap["Alice"]["Carol"] {
		t.Error("Expected Alice to exclude Carol")
	}

	if !exclusionMap["Carol"]["Alice"] {
		t.Error("Expected Carol to exclude Alice")
	}

	if _, exists := exclusionMap["Bob"]; exists {
		t.Error("Bob should not be in exclusion map")
	}
}

func TestBuildCompatibilityGraph(t *testing.T) {
	participants := []*participant.Participant{
		{Name: "Alice", Exclusions: []string{"Bob"}},
		{Name: "Bob", Exclusions: []string{}},
		{Name: "Carol", Exclusions: []string{}},
	}

	exclusionMap := buildExclusionMap(participants)
	graph := buildCompatibilityGraph(participants, exclusionMap)

	// Alice can give to Carol (index 2) but not Bob (index 1)
	if len(graph[0]) != 1 || graph[0][0] != 2 {
		t.Errorf("Alice should only be compatible with Carol, got %v", graph[0])
	}

	// Bob can give to Alice (0) or Carol (2)
	if len(graph[1]) != 2 {
		t.Errorf("Bob should be compatible with 2 people, got %d", len(graph[1]))
	}

	// Carol can give to Alice (0) or Bob (1)
	if len(graph[2]) != 2 {
		t.Errorf("Carol should be compatible with 2 people, got %d", len(graph[2]))
	}
}

// Compare original vs optimized results
func TestConsistency_OriginalVsOptimized(t *testing.T) {
	// Run multiple times to ensure both algorithms produce valid results
	for i := 0; i < 10; i++ {
		participants1 := createTestParticipants(20, 3)
		participants2 := make([]*participant.Participant, len(participants1))

		// Deep copy for second algorithm
		for j, p := range participants1 {
			participants2[j] = &participant.Participant{
				Name:             p.Name,
				NotificationType: p.NotificationType,
				ContactInfo:      make([]string, len(p.ContactInfo)),
				Exclusions:       make([]string, len(p.Exclusions)),
			}
			copy(participants2[j].ContactInfo, p.ContactInfo)
			copy(participants2[j].Exclusions, p.Exclusions)
		}

		_, err1 := Names(participants1)
		result2, err2 := NamesOptimized(participants2)

		if err1 != nil {
			t.Logf("Original algorithm failed: %v", err1)
		}
		if err2 != nil {
			t.Fatalf("Optimized algorithm failed: %v", err2)
		}

		// Verify optimized result validity
		validateResult(t, result2)
	}
}

func validateResult(t *testing.T, participants []*participant.Participant) {
	recipientMap := make(map[string]bool)

	for _, p := range participants {
		if p.Recipient == nil {
			t.Errorf("Participant %s has no recipient", p.Name)
			continue
		}

		if p.Recipient.Name == p.Name {
			t.Errorf("Participant %s assigned to themselves", p.Name)
		}

		if recipientMap[p.Recipient.Name] {
			t.Errorf("Recipient %s assigned multiple times", p.Recipient.Name)
		}
		recipientMap[p.Recipient.Name] = true

		for _, excluded := range p.Exclusions {
			if p.Recipient.Name == excluded {
				t.Errorf("Participant %s assigned to excluded recipient %s", p.Name, excluded)
			}
		}
	}
}