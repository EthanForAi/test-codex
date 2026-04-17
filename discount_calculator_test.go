package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestCalculateDiscount(t *testing.T) {
	tests := []struct {
		name        string
		price       float64
		coupon      string
		isVIP       bool
		want        float64
		wantErrText string
	}{
		{
			name:   "returns original price when coupon is empty",
			price:  100,
			coupon: "   ",
			want:   100,
		},
		{
			name:   "applies percent coupon",
			price:  200,
			coupon: "percent10",
			want:   180,
		},
		{
			name:   "rounds percent coupon result to two decimals",
			price:  99.99,
			coupon: "PERCENT15",
			want:   84.99,
		},
		{
			name:   "applies fixed amount coupon",
			price:  120,
			coupon: "FIXED30",
			want:   90,
		},
		{
			name:   "fixed amount coupon cannot reduce below zero",
			price:  20,
			coupon: "FIXED50",
			want:   0,
		},
		{
			name:   "applies full reduction coupon when threshold is met",
			price:  300,
			coupon: "FULL200-40",
			want:   260,
		},
		{
			name:        "returns error when full reduction threshold is not met",
			price:       199.99,
			coupon:      "FULL200-40",
			wantErrText: "coupon requires minimum spend 200.00",
		},
		{
			name:   "applies VIP coupon for VIP user",
			price:  500,
			coupon: "VIP20",
			isVIP:  true,
			want:   400,
		},
		{
			name:        "rejects VIP coupon for non VIP user",
			price:       500,
			coupon:      "VIP20",
			isVIP:       false,
			wantErrText: "coupon is only valid for VIP users",
		},
		{
			name:        "rejects negative price",
			price:       -1,
			coupon:      "PERCENT10",
			wantErrText: "price cannot be negative",
		},
		{
			name:        "rejects invalid percent coupon value",
			price:       100,
			coupon:      "PERCENT100",
			wantErrText: "percent coupon must be between 0 and 100",
		},
		{
			name:        "rejects unsupported coupon",
			price:       100,
			coupon:      "SPRINGSALE",
			wantErrText: "unsupported coupon: SPRINGSALE",
		},
		{
			name:        "rejects malformed full reduction coupon",
			price:       100,
			coupon:      "FULL200",
			wantErrText: "full reduction coupon must be in format FULL<price>-<discount>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := CalculateDiscount(tt.price, tt.coupon, tt.isVIP)

			if tt.wantErrText != "" {
				if err == nil {
					t.Fatalf("CalculateDiscount(%v, %q, %t) error = nil, want %q", tt.price, tt.coupon, tt.isVIP, tt.wantErrText)
				}
				if err.Error() != tt.wantErrText {
					t.Fatalf("CalculateDiscount(%v, %q, %t) error = %q, want %q", tt.price, tt.coupon, tt.isVIP, err.Error(), tt.wantErrText)
				}
				return
			}

			if err != nil {
				t.Fatalf("CalculateDiscount(%v, %q, %t) unexpected error = %v", tt.price, tt.coupon, tt.isVIP, err)
			}
			if got != tt.want {
				t.Fatalf("CalculateDiscount(%v, %q, %t) = %v, want %v", tt.price, tt.coupon, tt.isVIP, got, tt.want)
			}
		})
	}
}

func TestCalculateDiscountLogsSuccess(t *testing.T) {
	var logBuffer bytes.Buffer
	restoreLogger(t, &logBuffer)

	got, err := CalculateDiscount(200, "percent10", false)
	if err != nil {
		t.Fatalf("CalculateDiscount returned unexpected error: %v", err)
	}
	if got != 180 {
		t.Fatalf("CalculateDiscount returned %v, want 180", got)
	}

	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, `calculate discount started: price=200.00 coupon="percent10" normalized_coupon="PERCENT10" vip=false`) {
		t.Fatalf("expected start log, got %q", logOutput)
	}
	if !strings.Contains(logOutput, `calculate discount completed: rule="PERCENT10" original=200.00 final=180.00`) {
		t.Fatalf("expected completion log, got %q", logOutput)
	}
}

func TestCalculateDiscountLogsFailure(t *testing.T) {
	var logBuffer bytes.Buffer
	restoreLogger(t, &logBuffer)

	_, err := CalculateDiscount(100, "VIP20", false)
	if err == nil {
		t.Fatal("CalculateDiscount error = nil, want non-nil")
	}

	logOutput := logBuffer.String()
	if !strings.Contains(logOutput, `calculate discount failed: price=100.00 coupon="VIP20" vip=false error=coupon is only valid for VIP users`) {
		t.Fatalf("expected failure log, got %q", logOutput)
	}
}

func restoreLogger(t *testing.T, output *bytes.Buffer) {
	t.Helper()

	previousWriter := log.Writer()
	previousFlags := log.Flags()
	previousPrefix := log.Prefix()

	log.SetOutput(output)
	log.SetFlags(0)
	log.SetPrefix("")

	t.Cleanup(func() {
		log.SetOutput(previousWriter)
		log.SetFlags(previousFlags)
		log.SetPrefix(previousPrefix)
	})
}
