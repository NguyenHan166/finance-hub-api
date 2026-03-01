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

// BudgetRepository handles budget data operations
type BudgetRepository struct {
	collection *mongo.Collection
}

// NewBudgetRepository creates a new budget repository
func NewBudgetRepository(db *mongo.Database) *BudgetRepository {
	return &BudgetRepository{
		collection: db.Collection("budgets"),
	}
}

// Create creates a new budget
func (r *BudgetRepository) Create(userID string, req models.CreateBudgetRequest) (*models.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	budget := &models.Budget{
		ID:             uuid.New().String(),
		UserID:         userID,
		Month:          req.Month,
		Scope:          req.Scope,
		CategoryID:     req.CategoryID,
		Limit:          req.Limit,
		Spent:          0, // Will be calculated
		AlertEnabled:   req.AlertEnabled,
		AlertThreshold: req.AlertThreshold,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, budget)
	if err != nil {
		return nil, err
	}

	return budget, nil
}

// GetByID retrieves a budget by ID
func (r *BudgetRepository) GetByID(id, userID string) (*models.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var budget models.Budget
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&budget)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

// GetByMonth retrieves all budgets for a specific month
func (r *BudgetRepository) GetByMonth(userID, month string) ([]models.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"month":   month,
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []models.Budget
	if err = cursor.All(ctx, &budgets); err != nil {
		return nil, err
	}

	if budgets == nil {
		budgets = []models.Budget{}
	}

	return budgets, nil
}

// GetByMonthAndScope retrieves budget by month and scope
func (r *BudgetRepository) GetByMonthAndScope(userID, month, scope string, categoryID *string) (*models.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"month":   month,
		"scope":   scope,
	}

	if scope == "category" && categoryID != nil {
		filter["category_id"] = *categoryID
	} else if scope == "total" {
		filter["$or"] = []bson.M{
			{"category_id": nil},
			{"category_id": bson.M{"$exists": false}},
		}
	}

	var budget models.Budget
	err := r.collection.FindOne(ctx, filter).Decode(&budget)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &budget, nil
}

// Update updates a budget
func (r *BudgetRepository) Update(id, userID string, req models.UpdateBudgetRequest) (*models.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	setFields := update["$set"].(bson.M)

	if req.Limit != nil {
		setFields["limit"] = *req.Limit
	}
	if req.AlertEnabled != nil {
		setFields["alert_enabled"] = *req.AlertEnabled
	}
	if req.AlertThreshold != nil {
		setFields["alert_threshold"] = *req.AlertThreshold
	}

	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Budget

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &updated, nil
}

// UpdateSpent updates the spent amount for a budget
func (r *BudgetRepository) UpdateSpent(id, userID string, spent float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	update := bson.M{
		"$set": bson.M{
			"spent":      spent,
			"updated_at": time.Now(),
		},
	}

	_, err := r.collection.UpdateOne(ctx, filter, update)
	return err
}

// Delete deletes a budget
func (r *BudgetRepository) Delete(id, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return mongo.ErrNoDocuments
	}

	return nil
}

// GetAll retrieves all budgets for a user
func (r *BudgetRepository) GetAll(userID string) ([]models.Budget, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "month", Value: -1}, {Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var budgets []models.Budget
	if err = cursor.All(ctx, &budgets); err != nil {
		return nil, err
	}

	if budgets == nil {
		budgets = []models.Budget{}
	}

	return budgets, nil
}
