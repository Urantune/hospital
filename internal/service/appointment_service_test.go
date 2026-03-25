package service

import "testing"

func TestCanTransitionAppointment(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		expected bool
	}{
		{name: "created to pending payment", from: "CREATED", to: "PENDING_PAYMENT", expected: true},
		{name: "confirmed to completed invalid", from: "CONFIRMED", to: "COMPLETED", expected: false},
		{name: "in progress to completed", from: "IN_PROGRESS", to: "COMPLETED", expected: true},
		{name: "cancelled to confirmed invalid", from: "CANCELLED", to: "CONFIRMED", expected: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := canTransitionAppointment(tc.from, tc.to)
			if result != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, result)
			}
		})
	}
}
