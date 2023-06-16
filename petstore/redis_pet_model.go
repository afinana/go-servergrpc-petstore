package petstore

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

// PetModel represent a mgo database session with a pet data model
type PetModel struct {
	C *redis.Client
}

// All method will be used to get all records from pets table
func (m *PetModel) All() ([]PetEntity, error) {
	// Define variables
	ctx := context.TODO()
	pets := []PetEntity{}

	// Find all item
	iter := m.C.Scan(ctx, 0, "pets:*", 0).Iterator()

	for iter.Next(ctx) {
		// Find item by id
		data, err := m.C.Get(ctx, iter.Val()).Result()
		if err != nil {
			// Checks if the user was not found
			return nil, err
		}
		var pet PetEntity
		err = json.Unmarshal([]byte(data), &pet)
		// Checks if the user was not found
		if err != nil {
			return nil, err
		}
		pets = append(pets, pet)
	} // for
	return pets, nil
}

// FindByID will be used to find a pet registry by id
func (m *PetModel) FindByObjectID(id string) (*PetEntity, error) {
	ctx := context.Background()

	// Find order by id
	data, err := m.C.Get(ctx, id).Result()
	if err != nil {
		// Checks if the order was not found
		return nil, err
	}
	var pet PetEntity
	json.Unmarshal([]byte(data), &pet)
	return &pet, nil
}
func (m *PetModel) FindByID(id int64) (*PetEntity, error) {

	// Find pet by id
	petId := strconv.FormatInt(id, 10)
	return m.FindByObjectID(petId)

}

// Insert will be used to insert a new pet registry
func (m *PetModel) Insert(pet PetEntity) (*PetEntity, error) {
	// Add pet
	json, err := json.Marshal(pet)
	if err != nil {
		// Checks if the pet was not found
		return nil, err
	}

	ctx := context.Background()
	// Add pet with id
	err = m.C.Set(ctx, fmt.Sprintf("pet:%v", pet.Id), json, 0).Err()
	if err != nil {
		// Checks if the pet was not found
		return nil, err
	}

	// Add status to hset with id
	status_tag := fmt.Sprintf("pet_status:%v", pet.Status)
	pet_key := fmt.Sprintf("pet:%v", pet.Id)

	_, err = m.C.HSet(ctx, status_tag, pet_key, pet.Status).Result()
	if err != nil {
		// Checks if the hset was not found
		return nil, err
	}

	// Add tags to hset with id
	tags := pet.Tags
	for _, tag := range tags {
		_, err = m.C.HSet(ctx, fmt.Sprintf("pet_tags:%v", tag.Name), pet_key, tag.Name).Result()
		if err != nil {
			// Checks if the hset was not found
			return nil, err
		}
	}

	return &pet, nil
}

// Insert will be used to insert a new pet registry
func (m *PetModel) Update(pet PetEntity) (*PetEntity, error) {
	log.Printf("Update::FindByID of id:%d \n", pet.Id)

	// Clean pet register
	m.DeleteByRedisID(fmt.Sprintf("%v", pet.Id))

	return m.Insert(pet)
}

// Delete will be used to delete a pet registry
func (m *PetModel) Delete(id string) error {
	ctx := context.Background()
	// Delete pet by id
	err := m.C.Del(ctx, fmt.Sprintf("pet:%v", id)).Err()
	return err
}

// FindByStatus will be used to find a pet registry by status
func (m *PetModel) FindByStatus(status []string) ([]PetEntity, error) {

	return m.FindByTagsRedis("pet_status:", status)
}

// FindByTagsRedis will be used to find a pet registry by a list of statuses or tags
func (m *PetModel) FindByTagsRedis(prefix string, tags []string) ([]PetEntity, error) {

	// begin find
	ctx := context.Background()
	var pets []PetEntity
	for _, tag := range tags {

		key := fmt.Sprintf("%v%v", prefix, tag)

		log.Printf("FindByTagsRedis::HGet of keys=%s \n", key)
		// Get all ids of the given tag
		ids, err := m.C.HKeys(ctx, key).Result()
		if err != nil {
			// Checks if the pet was not found
			return nil, err
		}
		for _, id := range ids {
			log.Printf("FindByTagsRedis::FindByID of id=%s \n", id)
			pet, err := m.FindByRedisID(id)
			if err != nil {
				// Checks if the pet was not found
				break
			}
			pets = append(pets, *pet)

		}

	}
	return pets, nil
}

// FindByID will be used to find a pet registry by id
func (m *PetModel) FindByRedisID(id string) (*PetEntity, error) {

	ctx := context.Background()

	// Find pet by id
	data, err := m.C.Get(ctx, id).Result()
	if err != nil {
		// Checks if the pet was not found
		return nil, err
	}

	pet := PetEntity{}
	err = json.Unmarshal([]byte(data), &pet)
	if err != nil {
		panic(err)
	}
	return &pet, nil
}

// FindByID will be used to find a pet registry by id
func (m *PetModel) DeleteByRedisID(id string) error {

	ctx := context.Background()

	pet, err := m.FindByRedisID(id)
	if err == nil {

		// Clean old pet registry
		status_tag := fmt.Sprintf("pet_status:%v", pet.Status)
		pet_key := fmt.Sprintf("pet:%v", pet.Id)
		_ = m.C.Del(ctx, pet_key)
		_ = m.C.HDel(ctx, status_tag, pet_key)
		for _, tag := range pet.Tags {
			tag_key := fmt.Sprintf("%v%v", "pet_tags:", tag)
			_ = m.C.HDel(ctx, tag_key, pet_key)

		}

	}
	return err
}
