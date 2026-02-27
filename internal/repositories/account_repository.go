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

	account := &models.Account{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Balance:   req.Balance,
		Currency:  req.Currency,
		BankName:  req.BankName,
		IsActive:  true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, account)
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
	if req.BankName != nil {
		update["$set"].(bson.M)["bank_name"] = *req.BankName
	}
	if req.IsActive != nil {
		update["$set"].(bson.M)["is_active"] = *req.IsActive
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
