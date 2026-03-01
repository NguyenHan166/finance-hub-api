package repositories

import (
	"context"
	"finance-hub-api/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// AccountRepository handles account data operations
type AccountRepository struct {
	collection *mongo.Collection
}

// NewAccountRepository creates a new account repository
func NewAccountRepository(db *mongo.Database) *AccountRepository {
	return &AccountRepository{
		collection: db.Collection("accounts"),
	}
}

// Create creates a new account
func (r *AccountRepository) Create(userID string, req models.CreateAccountRequest) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get max display order
	maxOrder := 0
	filter := bson.M{"user_id": userID}
	opts := options.FindOne().SetSort(bson.D{{Key: "display_order", Value: -1}})
	var lastAccount models.Account
	err := r.collection.FindOne(ctx, filter, opts).Decode(&lastAccount)
	if err == nil {
		maxOrder = lastAccount.DisplayOrder
	}

	// Set display order
	displayOrder := maxOrder + 1
	if req.DisplayOrder != nil {
		displayOrder = *req.DisplayOrder
	}

	account := &models.Account{
		ID:                  uuid.New().String(),
		UserID:              userID,
		Name:                req.Name,
		Type:                req.Type,
		Balance:             req.Balance,
		Currency:            req.Currency,
		Icon:                req.Icon,
		Color:               req.Color,
		BankBIN:             req.BankBIN,
		BankCode:            req.BankCode,
		BankName:            req.BankName,
		AccountNumber:       req.AccountNumber,
		CardNumber:          req.CardNumber,
		CreditLimit:         req.CreditLimit,
		StatementDate:       req.StatementDate,
		DueDate:             req.DueDate,
		IsActive:            true,
		IsExcludedFromTotal: req.IsExcludedFromTotal,
		DisplayOrder:        displayOrder,
		CreatedAt:           time.Now(),
		UpdatedAt:           time.Now(),
	}

	_, err = r.collection.InsertOne(ctx, account)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// GetByID retrieves an account by ID
func (r *AccountRepository) GetByID(id, userID string) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var account models.Account
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&account)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetAll retrieves all accounts for a user
func (r *AccountRepository) GetAll(userID string, pagination models.PaginationQuery) ([]models.Account, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	// Get total count
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get accounts with pagination
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "created_at", Value: -1}})
	opts.SetLimit(int64(pagination.Limit))
	opts.SetSkip(int64(pagination.GetOffset()))

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var accounts []models.Account
	if err = cursor.All(ctx, &accounts); err != nil {
		return nil, 0, err
	}

	if accounts == nil {
		accounts = []models.Account{}
	}

	return accounts, int(totalCount), nil
}

// Update updates an account
func (r *AccountRepository) Update(id, userID string, req models.UpdateAccountRequest) (*models.Account, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	// Add fields to update if provided
	if req.Name != nil {
		update["$set"].(bson.M)["name"] = *req.Name
	}
	if req.Balance != nil {
		update["$set"].(bson.M)["balance"] = *req.Balance
	}
	if req.Icon != nil {
		update["$set"].(bson.M)["icon"] = *req.Icon
	}
	if req.Color != nil {
		update["$set"].(bson.M)["color"] = *req.Color
	}
	if req.BankBIN != nil {
		update["$set"].(bson.M)["bank_bin"] = *req.BankBIN
	}
	if req.BankCode != nil {
		update["$set"].(bson.M)["bank_code"] = *req.BankCode
	}
	if req.BankName != nil {
		update["$set"].(bson.M)["bank_name"] = *req.BankName
	}
	if req.AccountNumber != nil {
		update["$set"].(bson.M)["account_number"] = *req.AccountNumber
	}
	if req.CardNumber != nil {
		update["$set"].(bson.M)["card_number"] = *req.CardNumber
	}
	if req.CreditLimit != nil {
		update["$set"].(bson.M)["credit_limit"] = *req.CreditLimit
	}
	if req.StatementDate != nil {
		update["$set"].(bson.M)["statement_date"] = *req.StatementDate
	}
	if req.DueDate != nil {
		update["$set"].(bson.M)["due_date"] = *req.DueDate
	}
	if req.IsActive != nil {
		update["$set"].(bson.M)["is_active"] = *req.IsActive
	}
	if req.IsExcludedFromTotal != nil {
		update["$set"].(bson.M)["is_excluded_from_total"] = *req.IsExcludedFromTotal
	}
	if req.DisplayOrder != nil {
		update["$set"].(bson.M)["display_order"] = *req.DisplayOrder
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var account models.Account
	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&account)

	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &account, nil
}

// GetTotalBalance calculates total balance across all active accounts
func (r *AccountRepository) GetTotalBalance(userID string) (float64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":                userID,
		"is_active":              true,
		"is_excluded_from_total": false,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return 0, err
	}
	defer cursor.Close(ctx)

	var total float64
	for cursor.Next(ctx) {
		var account models.Account
		if err := cursor.Decode(&account); err != nil {
			continue
		}
		total += account.Balance
	}

	return total, nil
}

// GetSummary returns account summary statistics
func (r *AccountRepository) GetSummary(userID string) (*models.AccountSummary, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "is_active": true}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	summary := &models.AccountSummary{
		AccountsByType: make(map[string]int),
	}

	var totalBalance float64
	var accountCount int

	for cursor.Next(ctx) {
		var account models.Account
		if err := cursor.Decode(&account); err != nil {
			continue
		}

		accountCount++
		summary.AccountsByType[account.Type]++

		if !account.IsExcludedFromTotal {
			totalBalance += account.Balance
		}
	}

	summary.TotalAccounts = accountCount
	summary.TotalBalance = totalBalance
	summary.NetWorth = totalBalance

	return summary, nil
}

// CountByUser counts accounts for a user
func (r *AccountRepository) CountByUser(userID string) (int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "is_active": true}
	return r.collection.CountDocuments(ctx, filter)
}

// Delete deletes an account
func (r *AccountRepository) Delete(id, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "user_id": userID}
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// UpdateBalance updates account balance
func (r *AccountRepository) UpdateBalance(id, userID string, amount float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"_id": id, "user_id": userID}
	update := bson.M{
		"$inc": bson.M{"balance": amount},
		"$set": bson.M{"updated_at": time.Now()},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}
