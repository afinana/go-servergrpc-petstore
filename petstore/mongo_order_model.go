package petstore

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// OrderModel represent a mgo database session with a order data model
type OrderModel struct {
	C *mongo.Collection
}

// All method will be used to get all records from orders table
func (m *OrderModel) All() ([]OrderEntity, error) {
	// Define variables
	ctx := context.TODO()
	b := []OrderEntity{}

	// Find all orders
	orderCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = orderCursor.All(ctx, &b)
	if err != nil {
		return nil, err
	}

	return b, err
}

// FindByID will be used to find a order registry by id
func (m *OrderModel) FindByID(id string) (*OrderEntity, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find order by id
	var order = OrderEntity{}
	err = m.C.FindOne(context.TODO(), bson.M{"_id": p}).Decode(&order)
	if err != nil {
		// Checks if the order was not found
		if err == mongo.ErrNoDocuments {
			return nil, errors.New("ErrNoDocuments")
		}
		return nil, err
	}

	return &order, nil
}

// Insert will be used to insert a new order registry
func (m *OrderModel) Insert(order OrderEntity) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), order)
}

// Delete will be used to delete a order registry
func (m *OrderModel) Delete(id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(context.TODO(), bson.M{"_id": p})
}
