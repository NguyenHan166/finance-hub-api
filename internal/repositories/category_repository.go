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
		ParentID:  req.ParentID,
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

// Update updates a category
func (r *CategoryRepository) Update(id, userID string, req models.UpdateCategoryRequest) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Build update document
	update := bson.M{
		"$set": bson.M{
			"updated_at": time.Now(),
		},
	}

	setFields := update["$set"].(bson.M)

	if req.ParentID != nil {
		setFields["parent_id"] = req.ParentID
	}
	if req.Name != nil {
		setFields["name"] = *req.Name
	}
	if req.Type != nil {
		setFields["type"] = *req.Type
	}
	if req.Icon != nil {
		setFields["icon"] = req.Icon
	}
	if req.Color != nil {
		setFields["color"] = req.Color
	}

	filter := bson.M{
		"_id":     id,
		"user_id": userID,
	}

	opts := options.FindOneAndUpdate().SetReturnDocument(options.After)
	var updated models.Category

	err := r.collection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&updated)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, nil
		}
		return nil, err
	}

	return &updated, nil
}

// GetParentCategories retrieves all parent categories (no parent_id)
func (r *CategoryRepository) GetParentCategories(userID string) ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id": userID,
		"$or": []bson.M{
			{"parent_id": nil},
			{"parent_id": bson.M{"$exists": false}},
		},
	}
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

// GetChildCategories retrieves all child categories of a parent
func (r *CategoryRepository) GetChildCategories(userID, parentID string) ([]models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":   userID,
		"parent_id": parentID,
	}
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

// CountChildCategories counts the number of child categories
func (r *CategoryRepository) CountChildCategories(userID, categoryID string) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{
		"user_id":   userID,
		"parent_id": categoryID,
	}

	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return 0, err
	}

	return int(count), nil
}

// GetByIDWithoutUserCheck retrieves a category by ID without checking user_id (for validation)
func (r *CategoryRepository) GetByIDWithoutUserCheck(id string) (*models.Category, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var category models.Category
	filter := bson.M{"_id": id}

	err := r.collection.FindOne(ctx, filter).Decode(&category)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &category, nil
}
