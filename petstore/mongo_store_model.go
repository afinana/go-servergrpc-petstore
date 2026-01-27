package petstore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// StoreModel represent a mgo database session with an order data model
type StoreModel struct {
	C *mongo.Collection
}

// FindByID will be used to find an order registry by id
func (m *StoreModel) FindByID(id int64) (*OrderEntity, error) {
	var order OrderEntity
	err := m.C.FindOne(context.TODO(), bson.M{"id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Insert will be used to insert a new order registry
func (m *StoreModel) Insert(order OrderEntity) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(context.TODO(), order)
}

// DeleteByID will be used to delete an order registry by id
func (m *StoreModel) DeleteByID(id int64) (*mongo.DeleteResult, error) {
	return m.C.DeleteOne(context.TODO(), bson.M{"id": id})
}

// GetInventory will be used to get pet inventory by status
// This logic might belong here or in PetModel. Since it's /store/inventory, it's here.
// It needs to access the pets collection though.
func (m *StoreModel) GetInventory(ctx context.Context, petCollection *mongo.Collection) (map[string]int32, error) {
	pipeline := mongo.Pipeline{
		{{Key: "$group", Value: bson.D{{Key: "_id", Value: "$status"}, {Key: "count", Value: bson.D{{Key: "$sum", Value: 1}}}}}},
	}

	cursor, err := petCollection.Aggregate(ctx, pipeline)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	inventory := make(map[string]int32)
	for cursor.Next(ctx) {
		var result struct {
			ID    string `bson:"_id"`
			Count int32  `bson:"count"`
		}
		if err := cursor.Decode(&result); err != nil {
			return nil, err
		}
		inventory[result.ID] = result.Count
	}

	return inventory, nil
}
