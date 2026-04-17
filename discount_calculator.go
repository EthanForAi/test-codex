package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"
)

func CalculateDiscount(price float64, coupon string, isVIP bool) (float64, error) {
	if price < 0 {
		err := errors.New("price cannot be negative")
		log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, coupon, isVIP, err)
		return 0, err
	}

	normalizedCoupon := strings.ToUpper(strings.TrimSpace(coupon))
	log.Printf("calculate discount started: price=%.2f coupon=%q normalized_coupon=%q vip=%t", price, coupon, normalizedCoupon, isVIP)

	if normalizedCoupon == "" {
		finalPrice := roundToCurrency(price)
		log.Printf("calculate discount completed: rule=%q original=%.2f final=%.2f", "NONE", price, finalPrice)
		return finalPrice, nil
	}

	finalPrice := price

	switch {
	case strings.HasPrefix(normalizedCoupon, "PERCENT"):
		percent, err := parsePositiveNumber(strings.TrimPrefix(normalizedCoupon, "PERCENT"))
		if err != nil {
			wrappedErr := fmt.Errorf("invalid percent coupon: %w", err)
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, wrappedErr)
			return 0, wrappedErr
		}
		if percent <= 0 || percent >= 100 {
			err := errors.New("percent coupon must be between 0 and 100")
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
			return 0, err
		}
		finalPrice = price * (1 - percent/100)
	case strings.HasPrefix(normalizedCoupon, "FIXED"):
		discount, err := parsePositiveNumber(strings.TrimPrefix(normalizedCoupon, "FIXED"))
		if err != nil {
			wrappedErr := fmt.Errorf("invalid fixed coupon: %w", err)
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, wrappedErr)
			return 0, wrappedErr
		}
		if discount <= 0 {
			err := errors.New("fixed coupon discount must be greater than 0")
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
			return 0, err
		}
		finalPrice = price - discount
	case strings.HasPrefix(normalizedCoupon, "FULL"):
		thresholdText := strings.TrimPrefix(normalizedCoupon, "FULL")
		parts := strings.Split(thresholdText, "-")
		if len(parts) != 2 {
			err := errors.New("full reduction coupon must be in format FULL<price>-<discount>")
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
			return 0, err
		}

		threshold, err := parsePositiveNumber(parts[0])
		if err != nil {
			wrappedErr := fmt.Errorf("invalid full reduction threshold: %w", err)
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, wrappedErr)
			return 0, wrappedErr
		}
		discount, err := parsePositiveNumber(parts[1])
		if err != nil {
			wrappedErr := fmt.Errorf("invalid full reduction discount: %w", err)
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, wrappedErr)
			return 0, wrappedErr
		}
		if threshold <= 0 || discount <= 0 {
			err := errors.New("full reduction values must be greater than 0")
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
			return 0, err
		}
		if price < threshold {
			wrappedErr := fmt.Errorf("coupon requires minimum spend %.2f", threshold)
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, wrappedErr)
			return 0, wrappedErr
		}
		finalPrice = price - discount
	case strings.HasPrefix(normalizedCoupon, "VIP"):
		if !isVIP {
			err := errors.New("coupon is only valid for VIP users")
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
			return 0, err
		}
		percent, err := parsePositiveNumber(strings.TrimPrefix(normalizedCoupon, "VIP"))
		if err != nil {
			wrappedErr := fmt.Errorf("invalid VIP coupon: %w", err)
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, wrappedErr)
			return 0, wrappedErr
		}
		if percent <= 0 || percent >= 100 {
			err := errors.New("VIP coupon must be between 0 and 100")
			log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
			return 0, err
		}
		finalPrice = price * (1 - percent/100)
	default:
		err := fmt.Errorf("unsupported coupon: %s", normalizedCoupon)
		log.Printf("calculate discount failed: price=%.2f coupon=%q vip=%t error=%v", price, normalizedCoupon, isVIP, err)
		return 0, err
	}

	if finalPrice < 0 {
		finalPrice = 0
	}

	finalPrice = roundToCurrency(finalPrice)
	log.Printf("calculate discount completed: rule=%q original=%.2f final=%.2f", normalizedCoupon, price, finalPrice)
	return finalPrice, nil
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
