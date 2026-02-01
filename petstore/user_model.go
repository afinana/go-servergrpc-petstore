package petstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// UserModel handles database operations for User entities
type UserModel struct {
	Rdb *redis.Client
}

// Returns all users in the collection
func (m *UserModel) All(ctx context.Context) ([]User, error) {
	// Scan all user:{username} keys
	// Note: This relies on User struct being available in package
	// If User is different from UserEntity, we need to map.
	// Assuming User is the proto struct or similar.
	// But let's look at what we store. We store UserEntity (with password etc).
	// The return type is []User.

	var users []User

	iter := m.Rdb.Scan(ctx, 0, "user:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := m.Rdb.Get(ctx, key).Result()
		if err != nil {
			continue
		}

		// Unmarshal into... User? Or UserEntity?
		// Logic in mongo model was: Decode(&b) where b is []User.
		// So it expects User struct to be compatible with BSON/JSON.
		var u User
		if err := json.Unmarshal([]byte(val), &u); err == nil {
			users = append(users, u)
		}
	}

	return users, nil
}

// Finds a user by their username
func (m *UserModel) FindByName(ctx context.Context, username string) (*UserEntity, error) {
	key := fmt.Sprintf("user:%s", username)
	val, err := m.Rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var user UserEntity
	if err := json.Unmarshal([]byte(val), &user); err != nil {
		return nil, err
	}
	return &user, nil
}

// Inserts a new user into the collection
func (m *UserModel) Insert(ctx context.Context, user UserEntity) (string, error) {
	key := fmt.Sprintf("user:%s", user.Username)
	data, err := json.Marshal(user)
	if err != nil {
		return "", err
	}

	// Check if exists? Or overwrite?
	// API usually allows create.
	err = m.Rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return "", err
	}
	return user.Username, nil
}

// Updates an existing user's data
func (m *UserModel) Update(ctx context.Context, user UserEntity) (int64, error) {
	key := fmt.Sprintf("user:%s", user.Username)

	// Check existence
	exists, err := m.Rdb.Exists(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	if exists == 0 {
		return 0, fmt.Errorf("user not found")
	}

	data, err := json.Marshal(user)
	if err != nil {
		return 0, err
	}

	err = m.Rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return 0, err
	}
	return 1, nil
}

// Deletes a user by their username
func (m *UserModel) Delete(ctx context.Context, username string) (int64, error) {
	key := fmt.Sprintf("user:%s", username)
	return m.Rdb.Del(ctx, key).Result()
}

// FindByID - Deprecated in Redis adaptation as we don't have ObjectIDs
// If needed, we could scan or generic ID lookup, but "id" string was usually ObjectID hex.
func (m *UserModel) FindByID(ctx context.Context, id string) (*User, error) {
	return nil, fmt.Errorf("FindByID not implemented for Redis")
}
