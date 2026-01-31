package petstore

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// PetModel handles database operations for Pet entities
type PetModel struct {
	C *mongo.Collection
}

// Returns all pets in the collection
func (m *PetModel) All(ctx context.Context) ([]PetEntity, error) {
	// Define variables
	b := []PetEntity{}

	// Find all pets
	petCursor, err := m.C.Find(ctx, bson.M{})
	if err != nil {
		return nil, err
	}
	err = petCursor.All(ctx, &b)
	if err != nil {
		return nil, err
	}

	return b, err
}

// Finds a pet by its MongoDB ObjectID
func (m *PetModel) FindByObjectID(ctx context.Context, id string) (*PetEntity, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}

	// Find pet by id
	var pet = PetEntity{}
	err = m.C.FindOne(ctx, bson.M{"_id": p}).Decode(&pet)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

// Finds a pet by its numerical ID
func (m *PetModel) FindByID(ctx context.Context, id int64) (*PetEntity, error) {

	// Find pet by id
	var pet = PetEntity{}
	err := m.C.FindOne(ctx, bson.M{"id": id}).Decode(&pet)
	if err != nil {
		return nil, err
	}

	return &pet, nil
}

// Inserts a new pet into the collection
func (m *PetModel) Insert(ctx context.Context, pet PetEntity) (*mongo.InsertOneResult, error) {
	return m.C.InsertOne(ctx, pet)
}

// Updates an existing pet's data
func (m *PetModel) Update(ctx context.Context, pet PetEntity) (*mongo.UpdateResult, error) {
	// filter by id
	filter := bson.M{"id": pet.Id}
	update := bson.M{"$set": pet}
	return m.C.UpdateOne(ctx, filter, update)
}

// Deletes a pet by its MongoDB ObjectID
func (m *PetModel) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	p, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	return m.C.DeleteOne(ctx, bson.M{"_id": p})
}

// Deletes a pet by its numerical ID
func (m *PetModel) DeleteByID(ctx context.Context, id int64) (*mongo.DeleteResult, error) {
	return m.C.DeleteOne(ctx, bson.M{"id": id})
}

// Finds pets by a list of statuses
func (m *PetModel) FindByStatus(ctx context.Context, status []string) ([]PetEntity, error) {
	var filters []bson.M
	for _, item := range status {
		filters = append(filters, bson.M{"status": item})
	}
	filter := bson.M{"$or": filters}

	cursor, err := m.C.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var pets []PetEntity
	if err = cursor.All(ctx, &pets); err != nil {
		return nil, err
	}
	return pets, nil
}

// Finds pets by a list of tags
func (m *PetModel) FindBytags(ctx context.Context, tags []string) ([]PetEntity, error) {
	var filters []bson.M
	for _, tag := range tags {
		filters = append(filters, bson.M{"tags": bson.M{"$elemMatch": bson.M{"name": tag}}})
	}
	filter := bson.M{"$or": filters}

	cursor, err := m.C.Find(ctx, filter)
	if err != nil {
		return nil, err
	}

	var pets []PetEntity
	if err = cursor.All(ctx, &pets); err != nil {
		return nil, err
	}
	return pets, nil
}
