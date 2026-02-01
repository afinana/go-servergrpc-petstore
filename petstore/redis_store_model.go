package petstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
)

// StoreModel handles database operations for Order entities
type StoreModel struct {
	Rdb *redis.Client
}

// Finds an order by its numerical ID
func (m *StoreModel) FindByID(ctx context.Context, id int64) (*OrderEntity, error) {
	key := fmt.Sprintf("order:%d", id)
	val, err := m.Rdb.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrNotFound
		}
		return nil, err
	}

	var order OrderEntity
	if err := json.Unmarshal([]byte(val), &order); err != nil {
		return nil, err
	}
	return &order, nil
}

// Inserts a new order into the collection
func (m *StoreModel) Insert(ctx context.Context, order OrderEntity) (int64, error) {
	var id int64
	if order.Id == 0 {
		var err error
		id, err = m.Rdb.Incr(ctx, "order:id").Result()
		if err != nil {
			return 0, err
		}
		order.Id = id
	} else {
		id = order.Id
	}

	// Set ShipDate if missing, though typically business logic handles this.
	// Just standardizing what Mongo might have done (nothing special).

	key := fmt.Sprintf("order:%d", id)
	data, err := json.Marshal(order)
	if err != nil {
		return 0, err
	}

	err = m.Rdb.Set(ctx, key, data, 0).Err()
	if err != nil {
		return 0, err
	}

	return id, nil
}

// Deletes an order by its numerical ID
func (m *StoreModel) DeleteByID(ctx context.Context, id int64) (int64, error) {
	key := fmt.Sprintf("order:%d", id)
	result, err := m.Rdb.Del(ctx, key).Result()
	if err != nil {
		return 0, err
	}
	return result, nil
}

// Aggregates pet inventory by status
// We need to know which statuses exist.
// Ideally PetModel maintains a "known_statuses" set.
// For now, let's assume standard statuses or scan keys (slow but works for small/test).
// Or better, since we are in the same Redis, we can just look for keys matching `pet:status:*`
// But SCAN is better.
func (m *StoreModel) GetInventory(ctx context.Context) (map[string]int32, error) {
	inventory := make(map[string]int32)

	// Scan for keys with pattern pet:status:*
	iter := m.Rdb.Scan(ctx, 0, "pet:status:*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		// key is like "pet:status:available"
		// Extract status string?
		// Actually we just need the count.

		status := key[len("pet:status:"):] // simple substring

		count, err := m.Rdb.SCard(ctx, key).Result()
		if err != nil {
			continue // ignore error?
		}
		inventory[status] = int32(count)
	}
	if err := iter.Err(); err != nil {
		return nil, err
	}

	return inventory, nil
}
