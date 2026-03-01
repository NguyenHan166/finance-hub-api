package services

import (
	"errors"
	"finance-hub-api/internal/models"
	"finance-hub-api/internal/repositories"
	"time"
)

// ReportService handles report business logic
type ReportService struct {
	transactionRepo *repositories.TransactionRepository
	categoryRepo    *repositories.CategoryRepository
}

// NewReportService creates a new report service
func NewReportService(
	transactionRepo *repositories.TransactionRepository,
	categoryRepo *repositories.CategoryRepository,
) *ReportService {
	return &ReportService{
		transactionRepo: transactionRepo,
		categoryRepo:    categoryRepo,
	}
}

// GetOverview generates overview report for a date range
func (s *ReportService) GetOverview(userID string, startDate, endDate time.Time) (*models.OverviewReport, error) {
	// Get transactions for the period (exclude transfers)
	transactions, err := s.transactionRepo.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	var totalIncome, totalExpense float64
	transactionCount := 0

	for _, tx := range transactions {
		if tx.Type == "transfer" {
			continue // Exclude transfers
		}
		transactionCount++

		if tx.Type == "income" {
			totalIncome += tx.Amount
		} else if tx.Type == "expense" {
			totalExpense += tx.Amount
		}
	}

	netSaving := totalIncome - totalExpense
	savingRate := 0.0
	if totalIncome > 0 {
		savingRate = (netSaving / totalIncome) * 100
	}

	// Calculate days in range
	days := endDate.Sub(startDate).Hours() / 24
	if days == 0 {
		days = 1
	}
	avgDailyExpense := totalExpense / days

	// Calculate comparison with previous month
	prevMonthStart := startDate.AddDate(0, -1, 0)
	prevMonthEnd := endDate.AddDate(0, -1, 0)
	comparison, _ := s.calculateComparison(userID, prevMonthStart, prevMonthEnd, totalIncome, totalExpense, netSaving)

	return &models.OverviewReport{
		TotalIncome:         totalIncome,
		TotalExpense:        totalExpense,
		NetSaving:           netSaving,
		SavingRate:          savingRate,
		TransactionCount:    transactionCount,
		AvgDailyExpense:     avgDailyExpense,
		ComparedToPrevMonth: comparison,
	}, nil
}

// calculateComparison calculates percentage change compared to previous period
func (s *ReportService) calculateComparison(
	userID string,
	prevStart, prevEnd time.Time,
	currentIncome, currentExpense, currentSaving float64,
) (models.ComparisonMetrics, error) {
	prevTransactions, err := s.transactionRepo.GetByDateRange(userID, prevStart, prevEnd)
	if err != nil {
		return models.ComparisonMetrics{}, nil // Return zero on error
	}

	var prevIncome, prevExpense float64
	for _, tx := range prevTransactions {
		if tx.Type == "transfer" {
			continue
		}
		if tx.Type == "income" {
			prevIncome += tx.Amount
		} else if tx.Type == "expense" {
			prevExpense += tx.Amount
		}
	}
	prevSaving := prevIncome - prevExpense

	// Calculate percentage changes
	incomeChange := calculatePercentageChange(prevIncome, currentIncome)
	expenseChange := calculatePercentageChange(prevExpense, currentExpense)
	savingChange := calculatePercentageChange(prevSaving, currentSaving)

	return models.ComparisonMetrics{
		Income:  incomeChange,
		Expense: expenseChange,
		Saving:  savingChange,
	}, nil
}

// calculatePercentageChange calculates percentage change between old and new values
func calculatePercentageChange(oldValue, newValue float64) float64 {
	if oldValue == 0 {
		if newValue > 0 {
			return 100
		}
		return 0
	}
	return ((newValue - oldValue) / oldValue) * 100
}

// GetByCategory generates category breakdown report
func (s *ReportService) GetByCategory(userID string, startDate, endDate time.Time) ([]models.CategoryReport, error) {
	// Get expense transactions
	transactions, err := s.transactionRepo.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Filter expenses and group by category
	categoryMap := make(map[string]*models.CategoryReport)
	var totalExpense float64

	for _, tx := range transactions {
		if tx.Type != "expense" || tx.CategoryID == nil {
			continue
		}

		totalExpense += tx.Amount
		categoryID := *tx.CategoryID

		if _, exists := categoryMap[categoryID]; !exists {
			categoryMap[categoryID] = &models.CategoryReport{
				CategoryID:       categoryID,
				CategoryName:     "Unknown", // Will be filled later
				Amount:           0,
				TransactionCount: 0,
				Trend:            "stable",
			}
		}

		categoryMap[categoryID].Amount += tx.Amount
		categoryMap[categoryID].TransactionCount++
	}

	// Get category names
	for categoryID, report := range categoryMap {
		category, err := s.categoryRepo.GetByID(categoryID, userID)
		if err == nil && category != nil {
			report.CategoryName = category.Name
		}
	}

	// Calculate percentages and determine trends
	result := make([]models.CategoryReport, 0, len(categoryMap))
	for _, report := range categoryMap {
		if totalExpense > 0 {
			report.Percentage = (report.Amount / totalExpense) * 100
		}
		
		// Simple trend calculation (could be improved with historical data)
		report.Trend = "stable"
		
		result = append(result, *report)
	}

	// Sort by amount descending
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Amount > result[i].Amount {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// GetByMerchant generates merchant breakdown report
func (s *ReportService) GetByMerchant(userID string, startDate, endDate time.Time) ([]models.MerchantReport, error) {
	// Get expense transactions
	transactions, err := s.transactionRepo.GetByDateRange(userID, startDate, endDate)
	if err != nil {
		return nil, err
	}

	// Filter expenses and group by merchant
	merchantMap := make(map[string]*models.MerchantReport)
	var totalExpense float64

	for _, tx := range transactions {
		if tx.Type != "expense" || tx.Merchant == nil {
			continue
		}

		totalExpense += tx.Amount
		merchant := *tx.Merchant

		if _, exists := merchantMap[merchant]; !exists {
			merchantMap[merchant] = &models.MerchantReport{
				Merchant:         merchant,
				Amount:           0,
				TransactionCount: 0,
			}
		}

		merchantMap[merchant].Amount += tx.Amount
		merchantMap[merchant].TransactionCount++
	}

	// Calculate percentages
	result := make([]models.MerchantReport, 0, len(merchantMap))
	for _, report := range merchantMap {
		if totalExpense > 0 {
			report.Percentage = (report.Amount / totalExpense) * 100
		}
		result = append(result, *report)
	}

	// Sort by amount descending
	for i := 0; i < len(result); i++ {
		for j := i + 1; j < len(result); j++ {
			if result[j].Amount > result[i].Amount {
				result[i], result[j] = result[j], result[i]
			}
		}
	}

	return result, nil
}

// GetWeeklySpending generates weekly spending report for a month
func (s *ReportService) GetWeeklySpending(userID, month string, categoryID *string) ([]models.WeeklySpending, error) {
	// Parse month (format: YYYY-MM)
	monthStart, monthEnd, err := parseMonthRangeForReport(month)
	if err != nil {
		return nil, err
	}

	// Get transactions for the month
	transactions, err := s.transactionRepo.GetByDateRange(userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	// Filter expenses and by category if specified
	var filteredTxs []models.Transaction
	for _, tx := range transactions {
		if tx.Type != "expense" {
			continue
		}
		if categoryID != nil && (tx.CategoryID == nil || *tx.CategoryID != *categoryID) {
			continue
		}
		filteredTxs = append(filteredTxs, tx)
	}

	// Group by weeks
	weeks := []models.WeeklySpending{}
	currentWeekStart := getWeekStart(monthStart)
	weekNum := 1

	for currentWeekStart.Before(monthEnd) || currentWeekStart.Equal(monthEnd) {
		weekEnd := currentWeekStart.AddDate(0, 0, 7).Add(-time.Second)
		
		// Adjust to month boundaries
		effectiveStart := currentWeekStart
		if effectiveStart.Before(monthStart) {
			effectiveStart = monthStart
		}
		effectiveEnd := weekEnd
		if effectiveEnd.After(monthEnd) {
			effectiveEnd = monthEnd
		}

		// Get transactions for this week
		var weekAmount float64
		weekTxCount := 0
		for _, tx := range filteredTxs {
			if (tx.TransactionDate.After(effectiveStart) || tx.TransactionDate.Equal(effectiveStart)) &&
				(tx.TransactionDate.Before(effectiveEnd) || tx.TransactionDate.Equal(effectiveEnd)) {
				weekAmount += tx.Amount
				weekTxCount++
			}
		}

		weeks = append(weeks, models.WeeklySpending{
			Week:             "week_" + string(rune(weekNum+'0')),
			Label:            "Tuần " + string(rune(weekNum+'0')),
			Amount:           weekAmount,
			TransactionCount: weekTxCount,
		})

		currentWeekStart = currentWeekStart.AddDate(0, 0, 7)
		weekNum++

		if weekNum > 5 {
			break // Max 5 weeks per month
		}
	}

	return weeks, nil
}

// GetWeeklyCashflow generates weekly cashflow report for a month
func (s *ReportService) GetWeeklyCashflow(userID, month string) ([]models.WeeklyCashflow, error) {
	// Parse month
	monthStart, monthEnd, err := parseMonthRangeForReport(month)
	if err != nil {
		return nil, err
	}

	// Get transactions for the month
	transactions, err := s.transactionRepo.GetByDateRange(userID, monthStart, monthEnd)
	if err != nil {
		return nil, err
	}

	// Filter out transfers
	var filteredTxs []models.Transaction
	for _, tx := range transactions {
		if tx.Type != "transfer" {
			filteredTxs = append(filteredTxs, tx)
		}
	}

	// Group by weeks
	weeks := []models.WeeklyCashflow{}
	currentWeekStart := getWeekStart(monthStart)
	weekNum := 1

	for currentWeekStart.Before(monthEnd) || currentWeekStart.Equal(monthEnd) {
		weekEnd := currentWeekStart.AddDate(0, 0, 7).Add(-time.Second)
		
		// Adjust to month boundaries
		effectiveStart := currentWeekStart
		if effectiveStart.Before(monthStart) {
			effectiveStart = monthStart
		}
		effectiveEnd := weekEnd
		if effectiveEnd.After(monthEnd) {
			effectiveEnd = monthEnd
		}

		// Get transactions for this week
		var weekIncome, weekExpense float64
		for _, tx := range filteredTxs {
			if (tx.TransactionDate.After(effectiveStart) || tx.TransactionDate.Equal(effectiveStart)) &&
				(tx.TransactionDate.Before(effectiveEnd) || tx.TransactionDate.Equal(effectiveEnd)) {
				if tx.Type == "income" {
					weekIncome += tx.Amount
				} else if tx.Type == "expense" {
					weekExpense += tx.Amount
				}
			}
		}

		weeks = append(weeks, models.WeeklyCashflow{
			Week:    "week_" + string(rune(weekNum+'0')),
			Label:   "Tuần " + string(rune(weekNum+'0')),
			Income:  weekIncome,
			Expense: weekExpense,
			Net:     weekIncome - weekExpense,
		})

		currentWeekStart = currentWeekStart.AddDate(0, 0, 7)
		weekNum++

		if weekNum > 5 {
			break
		}
	}

	return weeks, nil
}

// getWeekStart returns the Monday of the week containing the given date
func getWeekStart(date time.Time) time.Time {
	weekday := int(date.Weekday())
	if weekday == 0 {
		weekday = 7 // Sunday = 7
	}
	daysToMonday := weekday - 1
	monday := date.AddDate(0, 0, -daysToMonday)
	return time.Date(monday.Year(), monday.Month(), monday.Day(), 0, 0, 0, 0, monday.Location())
}

// parseMonthRangeForReport parses "YYYY-MM" format and returns start and end of month
func parseMonthRangeForReport(month string) (time.Time, time.Time, error) {
	if len(month) != 7 || month[4] != '-' {
		return time.Time{}, time.Time{}, errors.New("invalid month format, expected YYYY-MM")
	}

	start, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, err
	}

	// Start of month
	monthStart := time.Date(start.Year(), start.Month(), 1, 0, 0, 0, 0, time.UTC)
	
	// End of month (last second of last day)
	monthEnd := monthStart.AddDate(0, 1, 0).Add(-time.Second)

	return monthStart, monthEnd, nil
}
