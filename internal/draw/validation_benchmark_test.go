package draw

import (
	"testing"
)

func BenchmarkValidateParticipants_10People(b *testing.B) {
	participants := createTestParticipants(10, 1)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipants(participants)
	}
}

func BenchmarkValidateParticipants_15People(b *testing.B) {
	participants := createTestParticipants(15, 2)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipants(participants)
	}
}

func BenchmarkValidateParticipants_20People(b *testing.B) {
	participants := createTestParticipants(20, 2)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipants(participants)
	}
}

func BenchmarkValidateParticipants_50People(b *testing.B) {
	participants := createTestParticipants(50, 3)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipants(participants)
	}
}

func BenchmarkValidateParticipants_100People(b *testing.B) {
	participants := createTestParticipants(100, 5)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipants(participants)
	}
}

func BenchmarkValidateParticipants_500People(b *testing.B) {
	participants := createTestParticipants(500, 15)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipants(participants)
	}
}

func BenchmarkValidateParticipantsQuick_10People(b *testing.B) {
	participants := createTestParticipants(10, 1)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipantsQuick(participants)
	}
}

func BenchmarkValidateParticipantsQuick_100People(b *testing.B) {
	participants := createTestParticipants(100, 5)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipantsQuick(participants)
	}
}

func BenchmarkValidateParticipantsQuick_500People(b *testing.B) {
	participants := createTestParticipants(500, 15)
	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		_ = ValidateParticipantsQuick(participants)
	}
}

// Benchmark validation + draw to measure total overhead
func BenchmarkValidateAndDraw_100People(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(100, 5)
		b.StartTimer()

		result := ValidateParticipants(participants)
		if !result.IsValid {
			b.Fatalf("Validation failed: %v", result.Errors)
		}

		_, err := Names(participants)
		if err != nil {
			b.Fatalf("Draw failed: %v", err)
		}
	}
}

func BenchmarkDrawOnly_100People(b *testing.B) {
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		b.StopTimer()
		participants := createTestParticipants(100, 5)
		b.StartTimer()

		_, err := Names(participants)
		if err != nil {
			b.Fatalf("Draw failed: %v", err)
		}
	}
}