package database

import (
	"context"
	"log"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Config holds database configuration
type Config struct {
	URI      string
	Database string
}

// Database wraps mongo.Database with additional functionality
type Database struct {
	Client   *mongo.Client
	Database *mongo.Database
}

// NewConnection creates a new MongoDB connection
func NewConnection(cfg Config) (*Database, error) {
	// Set client options
	clientOptions := options.Client().ApplyURI(cfg.URI)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}

	// Ping the database to verify connection
	if err := client.Ping(ctx, nil); err != nil {
		return nil, err
	}

	log.Println("✅ MongoDB connected successfully")

	return &Database{
		Client:   client,
		Database: client.Database(cfg.Database),
	}, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	log.Println("🔌 Closing database connection...")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	return d.Client.Disconnect(ctx)
}

// Health checks if the database is healthy
func (d *Database) Health() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return d.Client.Ping(ctx, nil)
}

// Collection returns a MongoDB collection
func (d *Database) Collection(name string) *mongo.Collection {
	return d.Database.Collection(name)
}
