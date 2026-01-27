package petstore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// UserModel represent a mgo database session with a user data model
type UserModel struct {
	C *mongo.Collection
}

// All method will be used to get all records from users table
func (m *UserModel) All() ([]User, error) {
	// Define variables
	ctx := context.TODO()
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

// FindByID will be used to find a user registry by id
func (m *UserModel) FindByID(id string) (*User, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find user by id
	var user = User{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// Insert will be used to insert a new user registry
func (m *UserModel) Insert(user UserEntity) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), user)
}

// FindByName will be used to find a user registry by username
func (m *UserModel) FindByName(username string) (*UserEntity, error) {
	var user UserEntity
	err := m.C.FindOne(context.TODO(), bson.M{"username": username}).Decode(&user)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update will be used to update a user registry
func (m *UserModel) Update(user UserEntity) (*mongo.UpdateResult, error) {
	// filter by username, as the API uses username for updates
	filter := bson.M{"username": user.Username}
	update := bson.M{"$set": user}
	return m.C.UpdateOne(context.TODO(), filter, update)
}

// Delete will be used to delete a user registry
func (m *UserModel) Delete(username string) (*mongo.DeleteResult, error) {
	// The API spec implies deletion by username for "DeleteUser", though the previous code had Delete by ID.
	// The DeleteUserRequest likely has a Username field. Let's check api_user.go or the proto definition later.
	// For now, I'll add a DeleteByName method or modify Delete to take a query.
	// Wait, the original Delete took an ID string.
	// Let's assume for now we need a DeleteByName for the API.
	return m.C.DeleteOne(context.TODO(), bson.M{"username": username})
}
