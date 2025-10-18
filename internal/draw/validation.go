package draw

import (
	"fmt"

	"github.com/igodwin/secretsanta/pkg/participant"
)

// ValidationResult provides information about constraint feasibility
type ValidationResult struct {
	IsValid                   bool
	Errors                    []string
	Warnings                  []string
	ParticipantsWithNoOptions []string
	MinCompatibility          int
	AvgCompatibility          float64
	TotalParticipants         int
}

// ValidateParticipants performs fast O(N²) validation of participant constraints
// This should be called BEFORE attempting a draw to catch impossible configurations
// Returns detailed validation results including any errors or warnings
func ValidateParticipants(participants []*participant.Participant) *ValidationResult {
	result := &ValidationResult{
		IsValid:           true,
		Errors:            make([]string, 0),
		Warnings:          make([]string, 0),
		TotalParticipants: len(participants),
	}

	n := len(participants)

	// Check for empty list
	if n == 0 {
		result.IsValid = false
		result.Errors = append(result.Errors, "no participants provided")
		return result
	}

	// Check for minimum participants
	if n < 2 {
		result.IsValid = false
		result.Errors = append(result.Errors, "need at least 2 participants for Secret Santa")
		return result
	}

	// Build exclusion map for O(1) lookups
	exclusionMap := buildExclusionMap(participants)

	// Check for duplicate names
	nameMap := make(map[string]bool)
	for _, p := range participants {
		if nameMap[p.Name] {
			result.IsValid = false
			result.Errors = append(result.Errors, fmt.Sprintf("duplicate participant name: %s", p.Name))
		}
		nameMap[p.Name] = true
	}

	totalCompatible := 0
	minCompatible := n

	// Check each participant's compatibility
	for i, giver := range participants {
		compatibleCount := 0

		// Validate contact info
		if len(giver.ContactInfo) == 0 {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("participant %s has no contact info", giver.Name))
		}

		// Check for invalid exclusions (non-existent participants)
		for _, excluded := range giver.Exclusions {
			if !nameMap[excluded] {
				result.Warnings = append(result.Warnings,
					fmt.Sprintf("participant %s excludes non-existent participant: %s", giver.Name, excluded))
			}
		}

		// Count compatible recipients
		for j, recipient := range participants {
			if i == j {
				continue // Can't give to self
			}

			// Check exclusions
			if exclusions, exists := exclusionMap[giver.Name]; exists {
				if exclusions[recipient.Name] {
					continue
				}
			}

			compatibleCount++
		}

		totalCompatible += compatibleCount
		if compatibleCount < minCompatible {
			minCompatible = compatibleCount
		}

		// If any participant has no valid recipients, assignment is impossible
		if compatibleCount == 0 {
			result.IsValid = false
			result.ParticipantsWithNoOptions = append(result.ParticipantsWithNoOptions, giver.Name)
			result.Errors = append(result.Errors,
				fmt.Sprintf("participant %s has no valid recipients (excluded everyone or too many exclusions)", giver.Name))
		}
	}

	result.MinCompatibility = minCompatible
	if n > 0 {
		result.AvgCompatibility = float64(totalCompatible) / float64(n)
	}

	// Warning if compatibility is very low (might be hard to find valid assignment)
	if result.IsValid && minCompatible < 2 && n > 3 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("low compatibility detected: some participants only have %d valid recipient(s)", minCompatible))
	}

	// Warning if average compatibility is below 50%
	avgCompat := result.AvgCompatibility
	if result.IsValid && avgCompat < float64(n)/2 {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("low average compatibility: %.1f out of %d possible recipients", avgCompat, n-1))
	}

	// Advanced validation: Check Hall's Marriage Theorem for small groups
	// This catches impossible configurations like two people both needing the same recipient
	// For N≤10: O(2^N) exhaustive check is acceptable (~400µs)
	// For N>10: Use O(N³) heuristics to avoid exponential slowdown
	if result.IsValid && n <= 10 {
		graph := buildCompatibilityGraph(participants, exclusionMap)
		if !checkHallsTheorem(graph, n) {
			result.IsValid = false
			result.Errors = append(result.Errors,
				"impossible configuration detected: constraints are too restrictive (Hall's Marriage Theorem violation)")
		}
	} else if result.IsValid && n > 10 {
		// For larger groups, use heuristic check
		if !checkHeuristicFeasibility(participants, exclusionMap) {
			result.IsValid = false
			result.Errors = append(result.Errors,
				"impossible configuration detected: constraints appear too restrictive")
		}
	}

	return result
}

// ValidateParticipantsQuick performs a quick validation check
// Returns true if valid, false otherwise (no detailed error messages)
func ValidateParticipantsQuick(participants []*participant.Participant) bool {
	n := len(participants)
	if n < 2 {
		return false
	}

	exclusionMap := buildExclusionMap(participants)

	// Quick check: ensure every participant has at least one valid recipient
	for i, giver := range participants {
		hasOption := false
		for j, recipient := range participants {
			if i == j {
				continue
			}

			// Check exclusions
			if exclusions, exists := exclusionMap[giver.Name]; exists {
				if exclusions[recipient.Name] {
					continue
				}
			}

			hasOption = true
			break // Found at least one valid recipient
		}

		if !hasOption {
			return false
		}
	}

	return true
}

// checkHallsTheorem verifies Hall's Marriage Theorem for the compatibility graph
// For every subset S of givers, the number of potential recipients must be >= |S|
// This is O(2^N) so only use for small N (<=20)
func checkHallsTheorem(graph [][]int, n int) bool {
	// Check all possible subsets of givers (2^n subsets)
	for mask := 1; mask < (1 << n); mask++ {
		subsetSize := 0
		recipients := make(map[int]bool)

		// For each giver in this subset
		for i := 0; i < n; i++ {
			if mask&(1<<i) != 0 {
				subsetSize++
				// Add all their potential recipients to the set
				for _, recipientIdx := range graph[i] {
					recipients[recipientIdx] = true
				}
			}
		}

		// Hall's condition: |recipients| must be >= |subset|
		if len(recipients) < subsetSize {
			return false
		}
	}

	return true
}

// checkHeuristicFeasibility performs heuristic feasibility checks for larger groups
// This is not as thorough as Hall's Theorem but runs in polynomial time
func checkHeuristicFeasibility(participants []*participant.Participant, exclusionMap map[string]map[string]bool) bool {
	n := len(participants)
	graph := buildCompatibilityGraph(participants, exclusionMap)

	// Heuristic 1: Total edges must be at least N
	totalEdges := 0
	for _, compatible := range graph {
		totalEdges += len(compatible)
	}
	if totalEdges < n {
		return false
	}

	// Heuristic 2: Check pairs of participants with very low compatibility
	// If two people can only receive from overlapping small sets, flag it
	for i := 0; i < n; i++ {
		for j := i + 1; j < n; j++ {
			// Find who can give to person i
			giversForI := make(map[int]bool)
			for k := 0; k < n; k++ {
				for _, recip := range graph[k] {
					if recip == i {
						giversForI[k] = true
						break
					}
				}
			}

			// Find who can give to person j
			giversForJ := make(map[int]bool)
			for k := 0; k < n; k++ {
				for _, recip := range graph[k] {
					if recip == j {
						giversForJ[k] = true
						break
					}
				}
			}

			// If both have very few givers and they overlap completely
			if len(giversForI) <= 2 && len(giversForJ) <= 2 {
				allOverlap := true
				for giver := range giversForI {
					if !giversForJ[giver] {
						allOverlap = false
						break
					}
				}
				if allOverlap && len(giversForI) == 1 {
					return false
				}
			}
		}
	}

	return true
}