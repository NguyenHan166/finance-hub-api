package utils

import (
	"fmt"
	"math"
	"regexp"
	"strings"
	"time"
)

// FormatCurrency formats a number as currency
func FormatCurrency(amount float64, currency string) string {
	return fmt.Sprintf("%s %.2f", currency, amount)
}

// RoundToTwoDecimals rounds a float to 2 decimal places
func RoundToTwoDecimals(value float64) float64 {
	return math.Round(value*100) / 100
}

// ValidateEmail validates email format
func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

// Slugify creates a URL-friendly slug from a string
func Slugify(text string) string {
	text = strings.ToLower(text)
	text = strings.TrimSpace(text)
	text = regexp.MustCompile(`[^a-z0-9]+`).ReplaceAllString(text, "-")
	text = strings.Trim(text, "-")
	return text
}

// ParseDateRange parses a date range string (e.g., "2024-01-01,2024-12-31")
func ParseDateRange(dateRange string) (time.Time, time.Time, error) {
	parts := strings.Split(dateRange, ",")
	if len(parts) != 2 {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid date range format")
	}

	startDate, err := time.Parse("2006-01-02", strings.TrimSpace(parts[0]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid start date: %w", err)
	}

	endDate, err := time.Parse("2006-01-02", strings.TrimSpace(parts[1]))
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid end date: %w", err)
	}

	return startDate, endDate, nil
}

// StringPtr returns a pointer to a string
func StringPtr(s string) *string {
	return &s
}

// Float64Ptr returns a pointer to a float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr returns a pointer to a bool
func BoolPtr(b bool) *bool {
	return &b
}

// Contains checks if a slice contains an element
func Contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
