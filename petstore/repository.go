package petstore

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

// PetRepository defines persistence operations for Pet entities
type PetRepository interface {
	All(ctx context.Context) ([]PetEntity, error)
	FindByObjectID(ctx context.Context, id string) (*PetEntity, error)
	FindByID(ctx context.Context, id int64) (*PetEntity, error)
	Insert(ctx context.Context, pet PetEntity) (*mongo.InsertOneResult, error)
	Update(ctx context.Context, pet PetEntity) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, id string) (*mongo.DeleteResult, error)
	DeleteByID(ctx context.Context, id int64) (*mongo.DeleteResult, error)
	FindByStatus(ctx context.Context, status []string) ([]PetEntity, error)
	FindBytags(ctx context.Context, tags []string) ([]PetEntity, error)
	EnsureIndexes(ctx context.Context) error
	GetCollection() *mongo.Collection
}

// StoreRepository defines persistence operations for Store/Order entities
type StoreRepository interface {
	FindByID(ctx context.Context, id int64) (*OrderEntity, error)
	Insert(ctx context.Context, order OrderEntity) (*mongo.InsertOneResult, error)
	DeleteByID(ctx context.Context, id int64) (*mongo.DeleteResult, error)
	GetInventory(ctx context.Context, petCollection *mongo.Collection) (map[string]int32, error)
	EnsureIndexes(ctx context.Context) error
}

// UserRepository defines persistence operations for User entities
type UserRepository interface {
	All(ctx context.Context) ([]UserEntity, error)
	FindByID(ctx context.Context, id string) (*UserEntity, error)
	Insert(ctx context.Context, user UserEntity) (*mongo.InsertOneResult, error)
	FindByName(ctx context.Context, username string) (*UserEntity, error)
	Update(ctx context.Context, user UserEntity) (*mongo.UpdateResult, error)
	Delete(ctx context.Context, username string) (*mongo.DeleteResult, error)
	EnsureIndexes(ctx context.Context) error
}
