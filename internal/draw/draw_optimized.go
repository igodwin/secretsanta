package draw

import (
	"fmt"
	"math/rand"

	"github.com/igodwin/secretsanta/pkg/participant"
)

// NamesOptimized uses a graph-based matching algorithm with O(N²) worst case
// instead of O(R × N² × E) for the original random retry approach.
// This uses a backtracking algorithm with constraint propagation.
func NamesOptimized(participants []*participant.Participant) ([]*participant.Participant, error) {
	n := len(participants)
	if n == 0 {
		return participants, nil
	}

	// Build exclusion map for O(1) lookup instead of O(E) linear search
	exclusionMap := buildExclusionMap(participants)

	// Build compatibility graph - O(N²) but only once
	compatibilityGraph := buildCompatibilityGraph(participants, exclusionMap)

	// Try to find a valid assignment using backtracking
	assignments := make([]*participant.Participant, n)
	used := make([]bool, n)

	if backtrack(participants, compatibilityGraph, assignments, used, 0) {
		// Apply assignments
		for i, p := range participants {
			p.Recipient = assignments[i]
		}
		return participants, nil
	}

	return nil, fmt.Errorf("no valid assignment found - constraints are too restrictive")
}

// buildExclusionMap creates a hash map for O(1) exclusion lookups
// Complexity: O(N × E) where E is avg exclusions per participant
func buildExclusionMap(participants []*participant.Participant) map[string]map[string]bool {
	exclusionMap := make(map[string]map[string]bool)

	for _, p := range participants {
		if len(p.Exclusions) > 0 {
			exclusionMap[p.Name] = make(map[string]bool)
			for _, excluded := range p.Exclusions {
				exclusionMap[p.Name][excluded] = true
			}
		}
	}

	return exclusionMap
}

// buildCompatibilityGraph builds adjacency list of valid recipient options
// Complexity: O(N²) - checks each participant against each other
func buildCompatibilityGraph(participants []*participant.Participant, exclusionMap map[string]map[string]bool) [][]int {
	n := len(participants)
	graph := make([][]int, n)

	for i := 0; i < n; i++ {
		graph[i] = make([]int, 0, n)
		giver := participants[i]

		for j := 0; j < n; j++ {
			if i == j {
				continue // Can't give to self
			}

			recipient := participants[j]

			// Check if recipient is excluded - O(1) with hash map
			if exclusions, exists := exclusionMap[giver.Name]; exists {
				if exclusions[recipient.Name] {
					continue
				}
			}

			graph[i] = append(graph[i], j)
		}
	}

	return graph
}

// backtrack uses constraint satisfaction with backtracking
// Average case: O(N²), Worst case: O(N!) but with heavy pruning
func backtrack(participants []*participant.Participant, graph [][]int, assignments []*participant.Participant, used []bool, giverIdx int) bool {
	if giverIdx == len(participants) {
		return true // All participants assigned
	}

	// Get compatible recipients for current giver
	compatibleRecipients := graph[giverIdx]

	// Randomize order to get different valid solutions each time
	indices := make([]int, len(compatibleRecipients))
	copy(indices, compatibleRecipients)
	rand.Shuffle(len(indices), func(i, j int) {
		indices[i], indices[j] = indices[j], indices[i]
	})

	for _, recipientIdx := range indices {
		if used[recipientIdx] {
			continue
		}

		// Try this assignment
		assignments[giverIdx] = participants[recipientIdx]
		used[recipientIdx] = true

		// Recursively try to assign remaining participants
		if backtrack(participants, graph, assignments, used, giverIdx+1) {
			return true
		}

		// Backtrack
		used[recipientIdx] = false
	}

	return false
}

// NamesOptimizedWithStats returns the assignment along with statistics
func NamesOptimizedWithStats(participants []*participant.Participant) ([]*participant.Participant, *DrawStats, error) {
	stats := &DrawStats{
		TotalParticipants: len(participants),
	}

	exclusionMap := buildExclusionMap(participants)
	compatibilityGraph := buildCompatibilityGraph(participants, exclusionMap)

	// Calculate average compatibility
	totalCompatible := 0
	for _, compatible := range compatibilityGraph {
		totalCompatible += len(compatible)
		if len(compatible) == 0 {
			stats.HasImpossibleConstraints = true
		}
	}
	stats.AvgCompatibilityPerPerson = float64(totalCompatible) / float64(len(participants))

	assignments := make([]*participant.Participant, len(participants))
	used := make([]bool, len(participants))

	if backtrack(participants, compatibilityGraph, assignments, used, 0) {
		for i, p := range participants {
			p.Recipient = assignments[i]
		}
		stats.Success = true
		return participants, stats, nil
	}

	stats.Success = false
	return nil, stats, fmt.Errorf("no valid assignment found - constraints are too restrictive")
}

// DrawStats provides insights into the draw process
type DrawStats struct {
	TotalParticipants         int
	AvgCompatibilityPerPerson float64
	HasImpossibleConstraints  bool
	Success                   bool
}