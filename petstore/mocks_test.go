package petstore

import (
	"context"
	"fmt"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// MockPetRepository implements PetRepository in-memory for unit testing
type MockPetRepository struct {
	mu   sync.RWMutex
	pets map[int64]PetEntity
}

func NewMockPetRepository() *MockPetRepository {
	return &MockPetRepository{
		pets: make(map[int64]PetEntity),
	}
}

func (m *MockPetRepository) All(ctx context.Context) ([]PetEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]PetEntity, 0, len(m.pets))
	for _, p := range m.pets {
		list = append(list, p)
	}
	return list, nil
}

func (m *MockPetRepository) FindByObjectID(ctx context.Context, id string) (*PetEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, p := range m.pets {
		if p.ID.Hex() == id {
			return &p, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (m *MockPetRepository) FindByID(ctx context.Context, id int64) (*PetEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	p, ok := m.pets[id]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	return &p, nil
}

func (m *MockPetRepository) Insert(ctx context.Context, pet PetEntity) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if pet.ID.IsZero() {
		pet.ID = primitive.NewObjectID()
	}
	m.pets[pet.Id] = pet
	return &mongo.InsertOneResult{InsertedID: pet.ID}, nil
}

func (m *MockPetRepository) Update(ctx context.Context, pet PetEntity) (*mongo.UpdateResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.pets[pet.Id]; !ok {
		return &mongo.UpdateResult{MatchedCount: 0, ModifiedCount: 0}, nil
	}
	m.pets[pet.Id] = pet
	return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
}

func (m *MockPetRepository) Delete(ctx context.Context, id string) (*mongo.DeleteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	for k, p := range m.pets {
		if p.ID.Hex() == id {
			delete(m.pets, k)
			return &mongo.DeleteResult{DeletedCount: 1}, nil
		}
	}
	return &mongo.DeleteResult{DeletedCount: 0}, nil
}

func (m *MockPetRepository) DeleteByID(ctx context.Context, id int64) (*mongo.DeleteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.pets[id]; !ok {
		return &mongo.DeleteResult{DeletedCount: 0}, nil
	}
	delete(m.pets, id)
	return &mongo.DeleteResult{DeletedCount: 1}, nil
}

func (m *MockPetRepository) FindByStatus(ctx context.Context, status []string) ([]PetEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	matchMap := make(map[string]bool)
	for _, s := range status {
		matchMap[s] = true
	}
	var res []PetEntity
	for _, p := range m.pets {
		if matchMap[p.Status] {
			res = append(res, p)
		}
	}
	return res, nil
}

func (m *MockPetRepository) FindBytags(ctx context.Context, tags []string) ([]PetEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	tagMap := make(map[string]bool)
	for _, t := range tags {
		tagMap[t] = true
	}
	var res []PetEntity
	for _, p := range m.pets {
		for _, tag := range p.Tags {
			if tagMap[tag.Name] {
				res = append(res, p)
				break
			}
		}
	}
	return res, nil
}

func (m *MockPetRepository) EnsureIndexes(ctx context.Context) error {
	return nil
}

func (m *MockPetRepository) GetCollection() *mongo.Collection {
	return nil
}

// MockStoreRepository implements StoreRepository in-memory for unit testing
type MockStoreRepository struct {
	mu     sync.RWMutex
	orders map[int64]OrderEntity
}

func NewMockStoreRepository() *MockStoreRepository {
	return &MockStoreRepository{
		orders: make(map[int64]OrderEntity),
	}
}

func (m *MockStoreRepository) FindByID(ctx context.Context, id int64) (*OrderEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	order, ok := m.orders[id]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	return &order, nil
}

func (m *MockStoreRepository) Insert(ctx context.Context, order OrderEntity) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.orders[order.Id] = order
	return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
}

func (m *MockStoreRepository) DeleteByID(ctx context.Context, id int64) (*mongo.DeleteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.orders[id]; !ok {
		return &mongo.DeleteResult{DeletedCount: 0}, nil
	}
	delete(m.orders, id)
	return &mongo.DeleteResult{DeletedCount: 1}, nil
}

func (m *MockStoreRepository) GetInventory(ctx context.Context, petCollection *mongo.Collection) (map[string]int32, error) {
	return map[string]int32{"STATUS_AVAILABLE": 5, "STATUS_PENDING": 2}, nil
}

func (m *MockStoreRepository) EnsureIndexes(ctx context.Context) error {
	return nil
}

// MockUserRepository implements UserRepository in-memory for unit testing
type MockUserRepository struct {
	mu    sync.RWMutex
	users map[string]UserEntity
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		users: make(map[string]UserEntity),
	}
}

func (m *MockUserRepository) All(ctx context.Context) ([]UserEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	list := make([]UserEntity, 0, len(m.users))
	for _, u := range m.users {
		list = append(list, u)
	}
	return list, nil
}

func (m *MockUserRepository) FindByID(ctx context.Context, id string) (*UserEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	for _, u := range m.users {
		if u.ID.Hex() == id {
			return &u, nil
		}
	}
	return nil, mongo.ErrNoDocuments
}

func (m *MockUserRepository) Insert(ctx context.Context, user UserEntity) (*mongo.InsertOneResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if user.Username == "" {
		return nil, fmt.Errorf("username required")
	}
	m.users[user.Username] = user
	return &mongo.InsertOneResult{InsertedID: primitive.NewObjectID()}, nil
}

func (m *MockUserRepository) FindByName(ctx context.Context, username string) (*UserEntity, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	u, ok := m.users[username]
	if !ok {
		return nil, mongo.ErrNoDocuments
	}
	return &u, nil
}

func (m *MockUserRepository) Update(ctx context.Context, user UserEntity) (*mongo.UpdateResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[user.Username]; !ok {
		return &mongo.UpdateResult{MatchedCount: 0, ModifiedCount: 0}, nil
	}
	m.users[user.Username] = user
	return &mongo.UpdateResult{MatchedCount: 1, ModifiedCount: 1}, nil
}

func (m *MockUserRepository) Delete(ctx context.Context, username string) (*mongo.DeleteResult, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, ok := m.users[username]; !ok {
		return &mongo.DeleteResult{DeletedCount: 0}, nil
	}
	delete(m.users, username)
	return &mongo.DeleteResult{DeletedCount: 1}, nil
}

func (m *MockUserRepository) EnsureIndexes(ctx context.Context) error {
	return nil
}
