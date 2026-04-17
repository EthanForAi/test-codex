package main

import (
	"errors"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func CalculateDiscount(price float64, coupon string, isVIP bool) (float64, error) {
	if price < 0 {
		return 0, errors.New("price cannot be negative")
	}

	normalizedCoupon := strings.ToUpper(strings.TrimSpace(coupon))
	if normalizedCoupon == "" {
		return roundToCurrency(price), nil
	}

	finalPrice := price

	switch {
	case strings.HasPrefix(normalizedCoupon, "PERCENT"):
		percent, err := parsePositiveNumber(strings.TrimPrefix(normalizedCoupon, "PERCENT"))
		if err != nil {
			return 0, fmt.Errorf("invalid percent coupon: %w", err)
		}
		if percent <= 0 || percent >= 100 {
			return 0, errors.New("percent coupon must be between 0 and 100")
		}
		finalPrice = price * (1 - percent/100)
	case strings.HasPrefix(normalizedCoupon, "FIXED"):
		discount, err := parsePositiveNumber(strings.TrimPrefix(normalizedCoupon, "FIXED"))
		if err != nil {
			return 0, fmt.Errorf("invalid fixed coupon: %w", err)
		}
		if discount <= 0 {
			return 0, errors.New("fixed coupon discount must be greater than 0")
		}
		finalPrice = price - discount
	case strings.HasPrefix(normalizedCoupon, "FULL"):
		thresholdText := strings.TrimPrefix(normalizedCoupon, "FULL")
		parts := strings.Split(thresholdText, "-")
		if len(parts) != 2 {
			return 0, errors.New("full reduction coupon must be in format FULL<price>-<discount>")
		}

		threshold, err := parsePositiveNumber(parts[0])
		if err != nil {
			return 0, fmt.Errorf("invalid full reduction threshold: %w", err)
		}
		discount, err := parsePositiveNumber(parts[1])
		if err != nil {
			return 0, fmt.Errorf("invalid full reduction discount: %w", err)
		}
		if threshold <= 0 || discount <= 0 {
			return 0, errors.New("full reduction values must be greater than 0")
		}
		if price < threshold {
			return 0, fmt.Errorf("coupon requires minimum spend %.2f", threshold)
		}
		finalPrice = price - discount
	case strings.HasPrefix(normalizedCoupon, "VIP"):
		if !isVIP {
			return 0, errors.New("coupon is only valid for VIP users")
		}
		percent, err := parsePositiveNumber(strings.TrimPrefix(normalizedCoupon, "VIP"))
		if err != nil {
			return 0, fmt.Errorf("invalid VIP coupon: %w", err)
		}
		if percent <= 0 || percent >= 100 {
			return 0, errors.New("VIP coupon must be between 0 and 100")
		}
		finalPrice = price * (1 - percent/100)
	default:
		return 0, fmt.Errorf("unsupported coupon: %s", normalizedCoupon)
	}

	if finalPrice < 0 {
		finalPrice = 0
	}

	return roundToCurrency(finalPrice), nil
}

func parsePositiveNumber(value string) (float64, error) {
	parsed, err := strconv.ParseFloat(strings.TrimSpace(value), 64)
	if err != nil {
		return 0, err
	}
	return parsed, nil
}

func roundToCurrency(value float64) float64 {
	return math.Round(value*100) / 100
}
