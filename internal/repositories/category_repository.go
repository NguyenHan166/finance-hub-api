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

// CategoryRepository handles category data operations
type CategoryRepository struct {
	collection *mongo.Collection
}

// NewCategoryRepository creates a new category repository
func NewCategoryRepository(db *mongo.Database) *CategoryRepository {
	return &CategoryRepository{
		collection: db.Collection("categories"),
	}
}

// Create creates a new category
func (r *CategoryRepository) Create(userID string, req models.CreateCategoryRequest) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	category := &models.Category{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      req.Name,
		Type:      req.Type,
		Icon:      req.Icon,
		Color:     req.Color,
		IsDefault: false,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := r.collection.InsertOne(ctx, category)
	if err != nil {
		return nil, err
	}

	return category, nil
}

// GetByID retrieves a category by ID
func (r *CategoryRepository) GetByID(id, userID string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var category models.Category
	filter := bson.M{"_id": id, "user_id": userID}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}

// GetAll retrieves all categories for a user
func (r *CategoryRepository) GetAll(userID string) ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID}
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	if categories == nil {
		categories = []models.Category{}
	}

	return categories, nil
}

// GetByType retrieves categories by type
func (r *CategoryRepository) GetByType(userID string, categoryType string) ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"user_id": userID, "type": categoryType}
	opts := options.Find().SetSort(bson.D{{Key: "name", Value: 1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var categories []models.Category
	if err = cursor.All(ctx, &categories); err != nil {
		return nil, err
	}

	if categories == nil {
		categories = []models.Category{}
	}

	return categories, nil
}

// Delete deletes a category
func (r *CategoryRepository) Delete(id, userID string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"_id":        id,
		"user_id":    userID,
		"is_default": false,
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
