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
	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token has expired")
	ErrTokenUsed     = errors.New("token already used")
)

// VerificationTokenRepository handles verification token database operations
type VerificationTokenRepository struct {
	collection *mongo.Collection
}

// NewVerificationTokenRepository creates a new verification token repository
func NewVerificationTokenRepository(db *mongo.Database) *VerificationTokenRepository {
	return &VerificationTokenRepository{
		collection: db.Collection("verification_tokens"),
	}
}

// Create creates a new verification token
func (r *VerificationTokenRepository) Create(ctx context.Context, token *models.VerificationToken) error {
	// Generate ID if not provided
	if token.ID == "" {
		token.ID = uuid.New().String()
	}

	// Set timestamps
	token.CreatedAt = time.Now()
	token.Used = false

	_, err := r.collection.InsertOne(ctx, token)
	if err != nil {
		return err
	}

	return nil
}

// FindByToken finds a verification token by its token string
func (r *VerificationTokenRepository) FindByToken(ctx context.Context, tokenStr string, tokenType string) (*models.VerificationToken, error) {
	var token models.VerificationToken

	filter := bson.M{
		"token": tokenStr,
		"type":  tokenType,
	}

	err := r.collection.FindOne(ctx, filter).Decode(&token)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, ErrTokenNotFound
		}
		return nil, err
	}

	// Check if expired
	if time.Now().After(token.ExpiresAt) {
		return nil, ErrTokenExpired
	}

	// Check if already used
	if token.Used {
		return nil, ErrTokenUsed
	}

	return &token, nil
}

// MarkAsUsed marks a token as used
func (r *VerificationTokenRepository) MarkAsUsed(ctx context.Context, tokenID string) error {
	filter := bson.M{"_id": tokenID}
	update := bson.M{"$set": bson.M{"used": true}}

	result, err := r.collection.UpdateOne(ctx, filter, update)
	if err != nil {
		return err
	}

	if result.MatchedCount == 0 {
		return ErrTokenNotFound
	}

	return nil
}

// DeleteByUserIDAndType deletes all tokens of a specific type for a user
func (r *VerificationTokenRepository) DeleteByUserIDAndType(ctx context.Context, userID string, tokenType string) error {
	filter := bson.M{
		"user_id": userID,
		"type":    tokenType,
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// DeleteExpired deletes all expired tokens
func (r *VerificationTokenRepository) DeleteExpired(ctx context.Context) error {
	filter := bson.M{
		"expires_at": bson.M{"$lt": time.Now()},
	}

	_, err := r.collection.DeleteMany(ctx, filter)
	return err
}

// FindByUserIDAndType finds tokens by user ID and type
func (r *VerificationTokenRepository) FindByUserIDAndType(ctx context.Context, userID string, tokenType string) ([]models.VerificationToken, error) {
	var tokens []models.VerificationToken

	filter := bson.M{
		"user_id": userID,
		"type":    tokenType,
		"used":    false,
		"expires_at": bson.M{"$gt": time.Now()},
	}

	cursor, err := r.collection.Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &tokens); err != nil {
		return nil, err
	}

	return tokens, nil
}
