package service

import "testing"

func TestCanTransitionPayment(t *testing.T) {
	tests := []struct {
		name     string
		from     string
		to       string
		expected bool
	}{
		{name: "initiated to success", from: "initiated", to: "success", expected: true},
		{name: "success to failed invalid", from: "success", to: "failed", expected: false},
		{name: "partial refunded to refunded", from: "partial_refunded", to: "refunded", expected: true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if got := canTransitionPayment(tc.from, tc.to); got != tc.expected {
				t.Fatalf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}

func TestIsDelayedOrDuplicatePaymentCallback(t *testing.T) {
	if !isDelayedOrDuplicatePaymentCallback("success", "failed") {
		t.Fatal("expected delayed failed callback after success to be ignored")
	}

	if !isDelayedOrDuplicatePaymentCallback("initiated", "initiated") {
		t.Fatal("expected same-state callback to be treated as duplicate")
	}

	if isDelayedOrDuplicatePaymentCallback("initiated", "success") {
		t.Fatal("expected valid success callback after initiated not to be ignored")
	}
}
