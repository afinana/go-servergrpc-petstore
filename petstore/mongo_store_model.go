package petstore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StoreModel handles database operations for Order entities
type StoreModel struct {
	C *mongo.Collection
}

// Finds an order by its numerical ID
func (m *StoreModel) FindByID(ctx context.Context, id int64) (*OrderEntity, error) {
	var order OrderEntity
	err := m.C.FindOne(ctx, bson.M{"id": id}).Decode(&order)
	if err != nil {
		return nil, err
	}
	return &order, nil
}

// Inserts a new order into the collection
func (m *StoreModel) Insert(ctx context.Context, order OrderEntity) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(ctx, order)
}

// Deletes an order by its numerical ID
func (m *StoreModel) DeleteByID(ctx context.Context, id int64) (*mongo.DeleteResult, error) {
	return m.C.DeleteOne(ctx, bson.M{"id": id})
}

// Aggregates pet inventory by status
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

// EnsureIndexes creates necessary indexes for the stores collection
func (m *StoreModel) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "id", Value: 1}},
			Options: options.Index().SetUnique(true).SetSparse(true),
		},
	}
	_, err := m.C.Indexes().CreateMany(ctx, models)
	return err
}
