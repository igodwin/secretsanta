package draw

import (
	"fmt"
	"math/rand"
	"testing"

	"github.com/igodwin/secretsanta/pkg/participant"
)

// Helper function to create test participants
func createTestParticipants(count int, exclusionsPerPerson int) []*participant.Participant {
	participants := make([]*participant.Participant, count)

	for i := 0; i < count; i++ {
		participants[i] = &participant.Participant{
			Name:             fmt.Sprintf("Person_%d", i),
			NotificationType: "email",
			ContactInfo:      []string{fmt.Sprintf("person%d@example.com", i)},
			Exclusions:       []string{},
		}
	}

	// Add random exclusions
	for i := 0; i < count; i++ {
		exclusionCount := 0
		for exclusionCount < exclusionsPerPerson && exclusionCount < count-1 {
			excludedIdx := rand.Intn(count)
			if excludedIdx == i {
				continue // Can't exclude self
			}

			excludedName := participants[excludedIdx].Name
			alreadyExcluded := false
			for _, e := range participants[i].Exclusions {
				if e == excludedName {
					alreadyExcluded = true
					break
				}
			}

			if !alreadyExcluded {
				participants[i].Exclusions = append(participants[i].Exclusions, excludedName)
				exclusionCount++
			}
		}
	}

	return participants
}

// Benchmark original algorithm with small dataset
func BenchmarkNamesOriginal_10People_1Exclusion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(10, 1)
		b.StartTimer()
		_, _ = Names(participants)
	}
}

func BenchmarkNamesOriginal_20People_2Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(20, 2)
		b.StartTimer()
		_, _ = Names(participants)
	}
}

func BenchmarkNamesOriginal_50People_3Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(50, 3)
		b.StartTimer()
		_, _ = Names(participants)
	}
}

func BenchmarkNamesOriginal_100People_5Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(100, 5)
		b.StartTimer()
		_, _ = Names(participants)
	}
}

func BenchmarkNamesOriginal_200People_10Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(200, 10)
		b.StartTimer()
		_, _ = Names(participants)
	}
}

func BenchmarkNamesOriginal_500People_15Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(500, 15)
		b.StartTimer()
		_, _ = Names(participants)
	}
}

// Benchmark optimized algorithm with same datasets
func BenchmarkNamesOptimized_10People_1Exclusion(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(10, 1)
		b.StartTimer()
		_, _ = NamesOptimized(participants)
	}
}

func BenchmarkNamesOptimized_20People_2Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(20, 2)
		b.StartTimer()
		_, _ = NamesOptimized(participants)
	}
}

func BenchmarkNamesOptimized_50People_3Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(50, 3)
		b.StartTimer()
		_, _ = NamesOptimized(participants)
	}
}

func BenchmarkNamesOptimized_100People_5Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(100, 5)
		b.StartTimer()
		_, _ = NamesOptimized(participants)
	}
}

func BenchmarkNamesOptimized_200People_10Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(200, 10)
		b.StartTimer()
		_, _ = NamesOptimized(participants)
	}
}

func BenchmarkNamesOptimized_500People_15Exclusions(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(500, 15)
		b.StartTimer()
		_, _ = NamesOptimized(participants)
	}
}

// Memory allocation benchmarks
func BenchmarkMemoryAlloc_Original_100People(b *testing.B) {
	participants := createTestParticipants(100, 5)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = Names(participants)
	}
}

func BenchmarkMemoryAlloc_Optimized_100People(b *testing.B) {
	participants := createTestParticipants(100, 5)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_, _ = NamesOptimized(participants)
	}
}