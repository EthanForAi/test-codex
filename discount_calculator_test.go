package main

import "testing"

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
