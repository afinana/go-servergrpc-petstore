package petstore

import "github.com/redis/go-redis/v9"

// StoreModel represent a mgo database session with a pet data model
type StoreModel struct {
	C *redis.Client
}
