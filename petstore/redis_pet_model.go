package petstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// PetModel handles database operations for Pet entities using Redis
type PetModel struct {
	Rdb *redis.Client
}

// Returns all pets in the collection
func (m *PetModel) All(ctx context.Context) ([]PetEntity, error) {
	// Find all keys matching pet:* (excluding indexes)
	// This is inefficient for large datasets, but matches the "All" semantics.
	// A better approach would be to maintain a "all_pets" set.
	// For now, let's use a set "pets:all" to store IDs.

	ids, err := m.Rdb.SMembers(ctx, "pets:all").Result()
	if err != nil {
		return nil, err
	}

	var pets []PetEntity
	for _, id := range ids {
		val, err := m.Rdb.Get(ctx, fmt.Sprintf("pet:%s", id)).Result()
		if err != nil {
			if err == redis.Nil {
				continue // Should not happen if data is consistent
			}
			return nil, err
		}

		var pet PetEntity
		if err := json.Unmarshal([]byte(val), &pet); err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	}

	return pets, nil
}

// Finds a pet by its numerical ID
func (m *PetModel) FindByID(ctx context.Context, id int64) (*PetEntity, error) {
	key := fmt.Sprintf("pet:%d", id)
	val, err := m.Rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var pet PetEntity
	if err := json.Unmarshal([]byte(val), &pet); err != nil {
		return nil, err
	}

	return &pet, nil
}

// Inserts a new pet into the collection
// Returns the inserted ID (int64)
func (m *PetModel) Insert(ctx context.Context, pet PetEntity) (int64, error) {
	// Generate ID if not provided (or even if provided, strictly speaking should be unique)
	// For simplicity, let's generate a new ID if it's 0.
	var id int64
	if pet.Id == 0 {
		var err error
		id, err = m.Rdb.Incr(ctx, "pet:id").Result()
		if err != nil {
			return 0, err
		}
		pet.Id = id
	} else {
		id = pet.Id
	}
	// Also keys
	key := fmt.Sprintf("pet:%d", id)

	data, err := json.Marshal(pet)
	if err != nil {
		return 0, err
	}

	// Transaction to save pet and update indexes
	pipe := m.Rdb.TxPipeline()
	pipe.Set(ctx, key, data, 0)
	pipe.SAdd(ctx, "pets:all", id)
	if pet.Status != "" {
		pipe.SAdd(ctx, fmt.Sprintf("pet:status:%s", pet.Status), id)
	}
	for _, tag := range pet.Tags {
		pipe.SAdd(ctx, fmt.Sprintf("pet:tag:%s", tag.Name), id)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Updates an existing pet's data
func (m *PetModel) Update(ctx context.Context, pet PetEntity) (int64, error) {
	// Check if exists first to handle "Not Found" correctly if needed,
	// or just overwrite. API usually expects 404 if not found.
	key := fmt.Sprintf("pet:%d", pet.Id)

	// Get old pet to remove from old indexes if status/tags changed
	oldVal, err := m.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, fmt.Errorf("pet not found")
	}
	if err != nil {
		return 0, err
	}

	var oldPet PetEntity
	json.Unmarshal([]byte(oldVal), &oldPet) // Ignore error, best effort

	data, err := json.Marshal(pet)
	if err != nil {
		return 0, err
	}

	pipe := m.Rdb.TxPipeline()
	pipe.Set(ctx, key, data, 0)

	// Mange indexes
	if oldPet.Status != pet.Status {
		if oldPet.Status != "" {
			pipe.SRem(ctx, fmt.Sprintf("pet:status:%s", oldPet.Status), pet.Id)
		}
		if pet.Status != "" {
			pipe.SAdd(ctx, fmt.Sprintf("pet:status:%s", pet.Status), pet.Id)
		}
	}

	// Tags - inefficient to check diff, just remove all old and add all new?
	// Or just add new. SRem old ones.
	for _, tag := range oldPet.Tags {
		pipe.SRem(ctx, fmt.Sprintf("pet:tag:%s", tag.Name), pet.Id)
	}
	for _, tag := range pet.Tags {
		pipe.SAdd(ctx, fmt.Sprintf("pet:tag:%s", tag.Name), pet.Id)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return 1, nil // Modified count
}

// Deletes a pet by its numerical ID
func (m *PetModel) DeleteByID(ctx context.Context, id int64) (int64, error) {
	key := fmt.Sprintf("pet:%d", id)

	// Get pet to remove from indexes
	val, err := m.Rdb.Get(ctx, key).Result()
	if err == redis.Nil {
		return 0, nil // Nothing to delete
	}
	if err != nil {
		return 0, err
	}

	var pet PetEntity
	json.Unmarshal([]byte(val), &pet)

	pipe := m.Rdb.TxPipeline()
	pipe.Del(ctx, key)
	pipe.SRem(ctx, "pets:all", id)
	if pet.Status != "" {
		pipe.SRem(ctx, fmt.Sprintf("pet:status:%s", pet.Status), id)
	}
	for _, tag := range pet.Tags {
		pipe.SRem(ctx, fmt.Sprintf("pet:tag:%s", tag.Name), id)
	}

	_, err = pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}

	return 1, nil
}

// Finds pets by a list of statuses
func (m *PetModel) FindByStatus(ctx context.Context, status []string) ([]PetEntity, error) {
	// Union of sets
	var keys []string
	for _, s := range status {
		keys = append(keys, fmt.Sprintf("pet:status:%s", s))
	}

	ids, err := m.Rdb.SUnion(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	return m.fetchPetsByIDs(ctx, ids)
}

// Finds pets by a list of tags
func (m *PetModel) FindBytags(ctx context.Context, tags []string) ([]PetEntity, error) {
	var keys []string
	for _, t := range tags {
		keys = append(keys, fmt.Sprintf("pet:tag:%s", t))
	}

	// Spec says "Tags to filter by". Usually OR? Swagger implies OR usually for this endpoint.
	ids, err := m.Rdb.SUnion(ctx, keys...).Result()
	if err != nil {
		return nil, err
	}

	return m.fetchPetsByIDs(ctx, ids)
}

func (m *PetModel) fetchPetsByIDs(ctx context.Context, ids []string) ([]PetEntity, error) {
	var pets []PetEntity
	for _, idStr := range ids {
		// idStr from redis set is string
		val, err := m.Rdb.Get(ctx, fmt.Sprintf("pet:%s", idStr)).Result()
		if err != nil {
			continue
		}
		var pet PetEntity
		if err := json.Unmarshal([]byte(val), &pet); err == nil {
			pets = append(pets, pet)
		}
	}
	return pets, nil
}

// Deletes a pet by its MongoDB ObjectID - DEPRECATED/UNUSED in Redis version, but kept if interface requires it?
// Actually, I'll remove it as I'm decoupling from Mongo.
