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
		CategoryID:      req.CategoryID,
		Type:            req.Type,
		Amount:          req.Amount,
		Description:     req.Description,
		TransactionDate: req.TransactionDate,
		Notes:           req.Notes,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
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

// GetAll retrieves all transactions for a user
func (r *TransactionRepository) GetAll(userID string, pagination models.PaginationQuery) ([]models.Transaction, int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}

	// Get total count
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Get transactions with pagination
	opts := options.Find()
	opts.SetSort(bson.D{{Key: "transaction_date", Value: -1}, {Key: "created_at", Value: -1}})
	opts.SetLimit(int64(pagination.Limit))
	opts.SetSkip(int64(pagination.GetOffset()))

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
