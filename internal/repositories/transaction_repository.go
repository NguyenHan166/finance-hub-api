package repositories

import (
	"context"
	"finance-hub-api/internal/models"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// TransactionRepository handles transaction data operations
type TransactionRepository struct {
	collection *mongo.Collection
}

// NewTransactionRepository creates a new transaction repository
func NewTransactionRepository(db *mongo.Database) *TransactionRepository {
	return &TransactionRepository{
		collection: db.Collection("transactions"),
	}
}

// Create creates a new transaction
func (r *TransactionRepository) Create(userID string, req models.CreateTransactionRequest) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	transaction := &models.Transaction{
		ID:              uuid.New().String(),
		UserID:          userID,
		AccountID:       req.AccountID,
		ToAccountID:     req.ToAccountID,
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Merchant:        req.Merchant,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		Notes:           req.Notes,
		Tags:            req.Tags,
		AttachmentURL:   req.AttachmentURL,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	// Initialize empty tags array if nil
	if transaction.Tags == nil {
		transaction.Tags = []string{}
	}

	_, err := r.collection.InsertOne(ctx, transaction)
	if err != nil {
		return nil, err
	}

	return transaction, nil
}

// GetByID retrieves a transaction by ID
func (r *TransactionRepository) GetByID(id, userID string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var transaction models.Transaction
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&transaction)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// GetAll retrieves all transactions for a user with filters
func (r *TransactionRepository) GetAll(userID string, filters models.TransactionFilterQuery) ([]models.Transaction, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter query
	filter := bson.M{"user_id": userID}

	// Account filter
	if filters.AccountID != "" {
		// Include transactions where AccountID or ToAccountID matches (for transfers)
		filter["$or"] = []bson.M{
			{"account_id": filters.AccountID},
			{"to_account_id": filters.AccountID},
		}
	}

	// Category filter
	if filters.CategoryID != "" {
		filter["category_id"] = filters.CategoryID
	}

	// Type filter
	if filters.Type != "" {
		filter["type"] = filters.Type
	}

	// Date range filters
	if filters.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", filters.StartDate)
		if err == nil {
			if filter["transaction_date"] == nil {
				filter["transaction_date"] = bson.M{}
			}
			filter["transaction_date"].(bson.M)["$gte"] = startDate
		}
	}

	if filters.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", filters.EndDate)
		if err == nil {
			// Add one day and use $lt to include the entire end date
			endDate = endDate.AddDate(0, 0, 1)
			if filter["transaction_date"] == nil {
				filter["transaction_date"] = bson.M{}
			}
			filter["transaction_date"].(bson.M)["$lt"] = endDate
		}
	}

	// Month filter (YYYY-MM format)
	if filters.Month != "" {
		// Parse month string
		parts := strings.Split(filters.Month, "-")
		if len(parts) == 2 {
			year, _ := strconv.Atoi(parts[0])
			month, _ := strconv.Atoi(parts[1])
			
			startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endOfMonth := startOfMonth.AddDate(0, 1, 0)
			
			filter["transaction_date"] = bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			}
		}
	}

	// Amount range filters
	if filters.MinAmount != "" {
		minAmount, err := strconv.ParseFloat(filters.MinAmount, 64)
		if err == nil {
			if filter["amount"] == nil {
				filter["amount"] = bson.M{}
			}
			filter["amount"].(bson.M)["$gte"] = minAmount
		}
	}

	if filters.MaxAmount != "" {
		maxAmount, err := strconv.ParseFloat(filters.MaxAmount, 64)
		if err == nil {
			if filter["amount"] == nil {
				filter["amount"] = bson.M{}
			}
			filter["amount"].(bson.M)["$lte"] = maxAmount
		}
	}

	// Search filter (search in merchant, description, notes)
	if filters.Search != "" {
		searchRegex := primitive.Regex{Pattern: filters.Search, Options: "i"}
		filter["$or"] = []bson.M{
			{"merchant": searchRegex},
			{"description": searchRegex},
			{"notes": searchRegex},
		}
	}

	// Tags filter
	if filters.Tags != "" {
		tags := strings.Split(filters.Tags, ",")
		filter["tags"] = bson.M{"$in": tags}
	}

	// Get total count
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Build sort options
	sortField := "transaction_date"
	if filters.SortBy == "amount" {
		sortField = "amount"
	}
	
	sortOrder := -1 // Default descending
	if filters.SortOrder == "asc" {
		sortOrder = 1
	}

	// Get transactions with pagination
	opts := options.Find()
	opts.SetSort(bson.D{{Key: sortField, Value: sortOrder}, {Key: "created_at", Value: -1}})
	
	// Apply pagination
	filters.SetDefaults()
	opts.SetLimit(int64(filters.Limit))
	opts.SetSkip(int64(filters.GetOffset()))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, 0, err
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	return transactions, int(totalCount), nil
}

// Update updates a transaction
func (r *TransactionRepository) Update(id, userID string, req models.UpdateTransactionRequest) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Only update provided fields
	if req.AccountID != nil {
		update["$set"].(bson.M)["account_id"] = *req.AccountID
	}
	if req.ToAccountID != nil {
		update["$set"].(bson.M)["to_account_id"] = *req.ToAccountID
	}
	if req.CategoryID != nil {
		update["$set"].(bson.M)["category_id"] = *req.CategoryID
	}
	if req.Type != nil {
		update["$set"].(bson.M)["type"] = *req.Type
	}
	if req.Amount != nil {
		update["$set"].(bson.M)["amount"] = *req.Amount
	}
	if req.Merchant != nil {
		update["$set"].(bson.M)["merchant"] = *req.Merchant
	}
	if req.Description != nil {
		update["$set"].(bson.M)["description"] = *req.Description
	}
	if req.TransactionDate != nil {
		update["$set"].(bson.M)["transaction_date"] = *req.TransactionDate
	}
	if req.Notes != nil {
		update["$set"].(bson.M)["notes"] = *req.Notes
	}
	if req.Tags != nil {
		update["$set"].(bson.M)["tags"] = req.Tags
	}
	if req.AttachmentURL != nil {
		update["$set"].(bson.M)["attachment_url"] = *req.AttachmentURL
	}

	filter := bson.M{"_id": id, "user_id": userID}

	var transaction models.Transaction
	err := r.collection.FindOneAndUpdate(
		ctx,
		filter,
		update,
		options.FindOneAndUpdate().SetReturnDocument(options.After),
	).Decode(&transaction)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// Delete deletes a transaction
func (r *TransactionRepository) Delete(id, userID string) (*models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var transaction models.Transaction
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOneAndDelete(ctx, filter).Decode(&transaction)
	if err == mongo.ErrNoDocuments {
		return nil, mongo.ErrNoDocuments
	}
	if err != nil {
		return nil, err
	}

	return &transaction, nil
}

// BulkUpdateCategory updates category for multiple transactions
func (r *TransactionRepository) BulkUpdateCategory(userID string, transactionIDs []string, categoryID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     bson.M{"$in": transactionIDs},
		"user_id": userID,
		"type":    bson.M{"$ne": "transfer"}, // Don't update transfers
	}

	update := bson.M{
		"$set": bson.M{
			"category_id": categoryID,
			"updated_at":  time.Now(),
		},
	}

	result, err := r.collection.UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}

	return result.ModifiedCount, nil
}

// BulkDelete deletes multiple transactions
func (r *TransactionRepository) BulkDelete(userID string, transactionIDs []string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     bson.M{"$in": transactionIDs},
		"user_id": userID,
	}

	result, err := r.collection.DeleteMany(ctx, filter)
	if err != nil {
		return 0, err
	}

	return result.DeletedCount, nil
}

// GetRecentTransactions retrieves the most recent transactions
func (r *TransactionRepository) GetRecentTransactions(userID string, limit int) ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	opts := options.Find()
	opts.SetSort(bson.D{{Key: "transaction_date", Value: -1}, {Key: "created_at", Value: -1}})
	opts.SetLimit(int64(limit))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	return transactions, nil
}

// GetSummary retrieves transaction summary statistics
func (r *TransactionRepository) GetSummary(userID string, filters models.TransactionFilterQuery) (*models.TransactionSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build filter (reuse logic from GetAll)
	filter := bson.M{"user_id": userID}

	// Add date filters if provided
	if filters.StartDate != "" {
		startDate, err := time.Parse("2006-01-02", filters.StartDate)
		if err == nil {
			if filter["transaction_date"] == nil {
				filter["transaction_date"] = bson.M{}
			}
			filter["transaction_date"].(bson.M)["$gte"] = startDate
		}
	}

	if filters.EndDate != "" {
		endDate, err := time.Parse("2006-01-02", filters.EndDate)
		if err == nil {
			endDate = endDate.AddDate(0, 0, 1)
			if filter["transaction_date"] == nil {
				filter["transaction_date"] = bson.M{}
			}
			filter["transaction_date"].(bson.M)["$lt"] = endDate
		}
	}

	if filters.Month != "" {
		parts := strings.Split(filters.Month, "-")
		if len(parts) == 2 {
			year, _ := strconv.Atoi(parts[0])
			month, _ := strconv.Atoi(parts[1])
			
			startOfMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
			endOfMonth := startOfMonth.AddDate(0, 1, 0)
			
			filter["transaction_date"] = bson.M{
				"$gte": startOfMonth,
				"$lt":  endOfMonth,
			}
		}
	}

	// Aggregation pipeline
	pipeline := []bson.M{
		{"$match": filter},
		{
			"$group": bson.M{
				"_id":   "$type",
				"count": bson.M{"$sum": 1},
				"total": bson.M{"$sum": "$amount"},
			},
		},
	}

	cursor, err := r.collection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	summary := &models.TransactionSummary{
		ByType: make(map[string]models.TransactionTypeStats),
	}

	for cursor.Next(ctx) {
		var result struct {
			ID    string  `bson:"_id"`
			Count int     `bson:"count"`
			Total float64 `bson:"total"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}

		summary.TotalTransactions += result.Count
		summary.ByType[result.ID] = models.TransactionTypeStats{
			Count:  result.Count,
			Amount: result.Total,
		}

		switch result.ID {
		case "income":
			summary.TotalIncome = result.Total
		case "expense":
			summary.TotalExpense = result.Total
		}
	}

	summary.NetAmount = summary.TotalIncome - summary.TotalExpense

	return summary, nil
}

// GetByAccountID retrieves all transactions for a specific account
func (r *TransactionRepository) GetByAccountID(userID, accountID string) ([]models.Transaction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"account_id": accountID},
			{"to_account_id": accountID},
		},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var transactions []models.Transaction
	if err = cursor.All(ctx, &transactions); err != nil {
		return nil, err
	}

	if transactions == nil {
		transactions = []models.Transaction{}
	}

	return transactions, nil
}

// CountByAccountID counts transactions for a specific account
func (r *TransactionRepository) CountByAccountID(userID, accountID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"account_id": accountID},
			{"to_account_id": accountID},
		},
	}

	return r.collection.CountDocuments(ctx, filter)
}

// GetTotalsByAccountID calculates income and expense totals for an account
func (r *TransactionRepository) GetTotalsByAccountID(userID, accountID string) (income, expense float64, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Income: transactions where this is the destination account
	incomeFilter := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"account_id": accountID, "type": "income"},
			{"to_account_id": accountID, "type": "transfer"},
		},
	}

	incomePipeline := []bson.M{
		{"$match": incomeFilter},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}

	incomeCursor, err := r.collection.Aggregate(ctx, incomePipeline)
	if err != nil {
		return 0, 0, err
	}
	defer incomeCursor.Close(ctx)

	if incomeCursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		if err := incomeCursor.Decode(&result); err == nil {
			income = result.Total
		}
	}

	// Expense: transactions where this is the source account
	expenseFilter := bson.M{
		"user_id":    userID,
		"account_id": accountID,
		"type":       bson.M{"$in": []string{"expense", "transfer"}},
	}

	expensePipeline := []bson.M{
		{"$match": expenseFilter},
		{"$group": bson.M{"_id": nil, "total": bson.M{"$sum": "$amount"}}},
	}

	expenseCursor, err := r.collection.Aggregate(ctx, expensePipeline)
	if err != nil {
		return income, 0, err
	}
	defer expenseCursor.Close(ctx)

	if expenseCursor.Next(ctx) {
		var result struct {
			Total float64 `bson:"total"`
		}
		if err := expenseCursor.Decode(&result); err == nil {
			expense = result.Total
		}
	}

	return income, expense, nil
}

