package petstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// UserModel represent a mgo database session with a user data model
type UserModel struct {
	C *redis.Client
}

// All method will be used to get all records from users table
func (m *UserModel) All() ([]UserEntity, error) {
	// Define variables
	ctx := context.TODO()
	users := []UserEntity{}

	// Find all item
	iter := m.C.Scan(ctx, 0, "users:*", 0).Iterator()

	for iter.Next(ctx) {
		// Find item by id
		data, err := m.C.Get(ctx, iter.Val()).Result()
		if err != nil {
			// Checks if the user was not found
			return nil, err
		}
		var user UserEntity
		err = json.Unmarshal([]byte(data), &user)
		// Checks if the user was not found
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	} // for
	return users, nil
}

// FindByID will be used to find a user registry by id
func (m *UserModel) FindByID(id string) (*UserEntity, error) {
	ctx := context.Background()

	// Find order by id
	data, err := m.C.Get(ctx, id).Result()
	if err != nil {
		// Checks if the order was not found
		return nil, err
	}
	var user UserEntity
	json.Unmarshal([]byte(data), &user)
	return &user, nil
}

// Insert will be used to insert a new user registry
func (m *UserModel) Insert(user UserEntity) (*UserEntity, error) {
	// Add user
	data, err := json.Marshal(user)
	ctx := context.Background()

	// Find user by id
	err = m.C.Set(ctx, fmt.Sprintf("user:%v", user.Id), data, 0).Err()
	if err != nil {
		// Checks if the user was not found
		return nil, err
	}
	return &user, nil
}

// Delete will be used to delete a user registry
func (m *UserModel) Delete(id string) error {
	ctx := context.Background()
	// Delete pet by id
	err := m.C.Del(ctx, fmt.Sprintf("user:%v", id)).Err()
	return err
}
