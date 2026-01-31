package petstore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserModel handles database operations for User entities
type UserModel struct {
	C *mongo.Collection
}

// Returns all users in the collection
func (m *UserModel) All(ctx context.Context) ([]User, error) {
	// Define variables
	b := []User{}

	// Find all users
	userCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = userCursor.All(ctx, &b)
	if err != nil {
		return nil, err
	}

	return b, err
}

// Finds a user by their MongoDB ObjectID
func (m *UserModel) FindByID(ctx context.Context, id string) (*User, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find user by id
	var user = User{}
	err = m.C.FindOne(ctx, bson.M{"_id": p}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Inserts a new user into the collection
func (m *UserModel) Insert(ctx context.Context, user UserEntity) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(ctx, user)
}

// Finds a user by their username
func (m *UserModel) FindByName(ctx context.Context, username string) (*UserEntity, error) {
	var user UserEntity
	err := m.C.FindOne(ctx, bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Updates an existing user's data
func (m *UserModel) Update(ctx context.Context, user UserEntity) (*mongo.UpdateResult, error) {
	// filter by username, as the API uses username for updates
	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": user}
	return m.C.UpdateOne(ctx, filter, update)
}

// Deletes a user by their username
func (m *UserModel) Delete(ctx context.Context, username string) (*mongo.DeleteResult, error) {
	return m.C.DeleteOne(ctx, bson.M{"username": username})
}
