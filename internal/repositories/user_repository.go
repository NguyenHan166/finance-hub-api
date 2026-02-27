package repositories

import (
	"context"
	"errors"
	"finance-hub-api/internal/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserAlreadyExists = errors.New("user already exists")
	ErrInvalidUserID     = errors.New("invalid user ID")
)

// UserRepository handles user database operations
type UserRepository struct {
	collection *mongo.Collection
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *mongo.Database) *UserRepository {
	return &UserRepository{
		collection: db.Collection("users"),
	}
}

// Create creates a new user
func (r *UserRepository) Create(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existing, _ := r.FindByEmail(ctx, user.Email)
	if existing != nil {
		return ErrUserAlreadyExists
	}

	// Generate ID if not provided
	if user.ID == "" {
		user.ID = uuid.New().String()
	}

	// Set timestamps
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now

	// Set defaults
	if user.AuthProvider == "" {
		user.AuthProvider = "email"
	}
	user.IsActive = true

	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}

	return nil
}

// FindByID finds a user by ID
func (r *UserRepository) FindByID(ctx context.Context, id string) (*models.User, error) {
	var user models.User
	
	filter := bson.M{"_id": id}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// FindByEmail finds a user by email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	
	filter := bson.M{"email": email}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// FindByGoogleID finds a user by Google ID
func (r *UserRepository) FindByGoogleID(ctx context.Context, googleID string) (*models.User, error) {
	var user models.User
	
	filter := bson.M{"google_id": googleID}
	err := r.collection.FindOne(ctx, filter).Decode(&user)
	
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

// Update updates a user
func (r *UserRepository) Update(ctx context.Context, id string, updates bson.M) error {
	// Add updated_at timestamp
	updates["updated_at"] = time.Now()

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdatePassword updates user password
func (r *UserRepository) UpdatePassword(ctx context.Context, id string, hashedPassword string) error {
	updates := bson.M{
		"password_hash": hashedPassword,
		"updated_at":    time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// UpdateLastLogin updates user's last login timestamp
func (r *UserRepository) UpdateLastLogin(ctx context.Context, id string) error {
	now := time.Now()
	updates := bson.M{
		"last_login_at": now,
		"updated_at":    now,
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// LinkGoogleAccount links a Google account to existing user
func (r *UserRepository) LinkGoogleAccount(ctx context.Context, userID string, googleID string, avatarURL *string) error {
	updates := bson.M{
		"google_id":     googleID,
		"auth_provider": "google",
		"updated_at":    time.Now(),
	}

	if avatarURL != nil && *avatarURL != "" {
		updates["avatar_url"] = *avatarURL
	}

	filter := bson.M{"_id": userID}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// Delete soft deletes a user (sets is_active to false)
func (r *UserRepository) Delete(ctx context.Context, id string) error {
	updates := bson.M{
		"is_active":  false,
		"updated_at": time.Now(),
	}

	filter := bson.M{"_id": id}
	update := bson.M{"$set": updates}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// HardDelete permanently deletes a user
func (r *UserRepository) HardDelete(ctx context.Context, id string) error {
	filter := bson.M{"_id": id}
	
	result, err := r.collection.DeleteOne(ctx, filter)
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return ErrUserNotFound
	}

	return nil
}

// ToUserProfile converts User to UserProfile (safe for client)
func ToUserProfile(user *models.User) models.UserProfile {
	return models.UserProfile{
		ID:            user.ID,
		Email:         user.Email,
		FullName:      user.FullName,
		AvatarURL:     user.AvatarURL,
		AuthProvider:  user.AuthProvider,
		EmailVerified: user.EmailVerified,
		LastLoginAt:   user.LastLoginAt,
		CreatedAt:     user.CreatedAt,
	}
}
