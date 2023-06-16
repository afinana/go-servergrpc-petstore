package petstore

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// OrderModel represent a mgo database session with a order data model
type OrderModel struct {
	C *redis.Client
}

// All method will be used to get all records from orders table
func (m *OrderModel) All() ([]OrderEntity, error) {
	// Define variables
	ctx := context.TODO()
	orders := []OrderEntity{}

	// Find all item
	iter := m.C.Scan(ctx, 0, "orders:*", 0).Iterator()

	for iter.Next(ctx) {
		// Find item by id
		data, err := m.C.Get(ctx, iter.Val()).Result()
		if err != nil {
			// Checks if the user was not found
			return nil, err
		}
		var order OrderEntity
		err = json.Unmarshal([]byte(data), &order)
		// Checks if the user was not found
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	} // for
	return orders, nil
}

// FindByID will be used to find a order registry by id
func (m *OrderModel) FindByID(id string) (*OrderEntity, error) {

	ctx := context.Background()

	// Find order by id
	data, err := m.C.Get(ctx, id).Result()
	if err != nil {
		// Checks if the order was not found
		return nil, err
	}
	var order OrderEntity
	json.Unmarshal([]byte(data), &order)
	return &order, nil

}

// Insert will be used to insert a new order registry
func (m *OrderModel) Insert(order OrderEntity) (*OrderEntity, error) {
	// Add user
	data, err := json.Marshal(order)
	ctx := context.Background()

	// Find user by id
	err = m.C.Set(ctx, fmt.Sprintf("order:%v", order.Id), data, 0).Err()
	if err != nil {
		// Checks if the user was not found
		return nil, err
	}
	return &order, nil
}

// Delete will be used to delete a order registry
func (m *OrderModel) Delete(id string) error {
	ctx := context.Background()
	// Delete user by id
	err := m.C.Del(ctx, id).Err()
	if err != nil {
		// Checks if the user was not found
		return err
	}
	return nil
}
