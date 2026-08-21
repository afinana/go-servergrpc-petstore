package petstore

import (
	"context"
	"testing"

	"go.mongodb.org/mongo-driver/bson"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestAddPetAndGetPetById(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	// Clean up before test
	_, err := testApp.pets.GetCollection().DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up pets collection: %v", err)
	}

	// 1. Add Pet
	petID := int64(1001)
	petName := "Fluffy"
	pet := &Pet{
		Id:        petID,
		Name:      petName,
		Status:    Pet_STATUS_AVAILABLE,
		PhotoUrls: []string{"url1", "url2"},
		Tags:      []*Tag{{Id: 1, Name: "cute"}},
		Category:  &Category{Id: 1, Name: "Dogs"},
	}
	addReq := &AddPetRequest{
		Body: pet,
	}

	_, err = testApp.AddPet(ctx, addReq)
	if err != nil {
		t.Fatalf("AddPet failed: %v", err)
	}

	// 2. Get Pet By ID
	getReq := &GetPetByIdRequest{
		PetId: petID,
	}
	retrievedPet, err := testApp.GetPetById(ctx, getReq)
	if err != nil {
		t.Fatalf("GetPetById failed: %v", err)
	}

	if retrievedPet.Id != petID {
		t.Errorf("Expected pet ID %d, got %d", petID, retrievedPet.Id)
	}
	if retrievedPet.Name != petName {
		t.Errorf("Expected pet name %s, got %s", petName, retrievedPet.Name)
	}
	if retrievedPet.Category == nil || retrievedPet.Category.Name != "Dogs" {
		t.Errorf("Expected Category 'Dogs', got %+v", retrievedPet.Category)
	}
}

func TestGetPetById_NotFound(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	getReq := &GetPetByIdRequest{
		PetId: 999999,
	}
	_, err := testApp.GetPetById(ctx, getReq)
	if err == nil {
		t.Fatalf("Expected error for nonexistent pet, got nil")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound status code, got %v", status.Code(err))
	}
}

func TestFindPetsByStatus(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	// Clean up before test
	_, err := testApp.pets.GetCollection().DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up pets collection: %v", err)
	}

	// Add pets with different statuses
	pets := []*Pet{
		{Id: 2001, Name: "P1", Status: Pet_STATUS_AVAILABLE},
		{Id: 2002, Name: "P2", Status: Pet_STATUS_PENDING},
		{Id: 2003, Name: "P3", Status: Pet_STATUS_SOLD},
	}

	for _, p := range pets {
		_, err := testApp.AddPet(ctx, &AddPetRequest{Body: p})
		if err != nil {
			t.Fatalf("Failed to add pet %s: %v", p.Name, err)
		}
	}

	// Test Find Available
	req := &FindPetsByStatusRequest{
		Status: []FindPetsByStatusRequest_Status{FindPetsByStatusRequest_STATUS_AVAILABLE},
	}
	res, err := testApp.FindPetsByStatus(ctx, req)
	if err != nil {
		t.Fatalf("FindPetsByStatus failed: %v", err)
	}

	if len(res.Items) != 1 {
		t.Errorf("Expected 1 available pet, got %d", len(res.Items))
		return
	}
	if res.Items[0].Name != "P1" {
		t.Errorf("Expected pet P1, got %s", res.Items[0].Name)
	}
}

func TestFindPetsByTags(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	// Clean up before test
	_, err := testApp.pets.GetCollection().DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up pets collection: %v", err)
	}

	// Add pets with tags
	pets := []*Pet{
		{Id: 3001, Name: "TaggedPet1", Tags: []*Tag{{Id: 1, Name: "energetic"}}},
		{Id: 3002, Name: "TaggedPet2", Tags: []*Tag{{Id: 2, Name: "calm"}}},
	}

	for _, p := range pets {
		_, err := testApp.AddPet(ctx, &AddPetRequest{Body: p})
		if err != nil {
			t.Fatalf("Failed to add pet %s: %v", p.Name, err)
		}
	}

	// Find by tag
	req := &FindPetsByTagsRequest{Tags: []string{"energetic"}}
	res, err := testApp.FindPetsByTags(ctx, req)
	if err != nil {
		t.Fatalf("FindPetsByTags failed: %v", err)
	}

	if len(res.Items) != 1 {
		t.Fatalf("Expected 1 pet with tag energetic, got %d", len(res.Items))
	}
	if res.Items[0].Name != "TaggedPet1" {
		t.Errorf("Expected pet 'TaggedPet1', got '%s'", res.Items[0].Name)
	}
}

func TestUpdatePet(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	// Clean up
	_, err := testApp.pets.GetCollection().DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up pets collection: %v", err)
	}

	// Add pet
	petID := int64(4001)
	pet := &Pet{Id: petID, Name: "Original", Status: Pet_STATUS_AVAILABLE}
	_, err = testApp.AddPet(ctx, &AddPetRequest{Body: pet})
	if err != nil {
		t.Fatalf("AddPet failed: %v", err)
	}

	// Update pet
	updated := &Pet{Id: petID, Name: "UpdatedName", Status: Pet_STATUS_SOLD}
	_, err = testApp.UpdatePet(ctx, &UpdatePetRequest{Body: updated})
	if err != nil {
		t.Fatalf("UpdatePet failed: %v", err)
	}

	// Verify
	retrieved, err := testApp.GetPetById(ctx, &GetPetByIdRequest{PetId: petID})
	if err != nil {
		t.Fatalf("GetPetById failed: %v", err)
	}
	if retrieved.Name != "UpdatedName" {
		t.Errorf("Expected name 'UpdatedName', got '%s'", retrieved.Name)
	}
	if retrieved.Status != Pet_STATUS_SOLD {
		t.Errorf("Expected status SOLD, got %v", retrieved.Status)
	}
}

func TestDeletePet(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	// Clean up before test
	_, err := testApp.pets.GetCollection().DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up pets collection: %v", err)
	}

	// Add a pet
	petID := int64(5001)
	pet := &Pet{
		Id:   petID,
		Name: "ToDelete",
	}
	_, err = testApp.AddPet(ctx, &AddPetRequest{Body: pet})
	if err != nil {
		t.Fatalf("Failed to add pet: %v", err)
	}

	// Delete Pet
	delReq := &DeletePetRequest{
		PetId: petID,
	}
	_, err = testApp.DeletePet(ctx, delReq)
	if err != nil {
		t.Fatalf("DeletePet failed: %v", err)
	}

	// Verify deletion
	getReq := &GetPetByIdRequest{
		PetId: petID,
	}
	_, err = testApp.GetPetById(ctx, getReq)
	if err == nil {
		t.Errorf("Expected error when getting deleted pet, got nil")
	}
}

func TestUpdatePetWithForm(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	// Clean up before test
	_, err := testApp.pets.GetCollection().DeleteMany(ctx, bson.M{})
	if err != nil {
		t.Fatalf("Failed to clean up pets collection: %v", err)
	}

	// 1. Add Pet
	petID := int64(6001)
	petName := "OriginalName"
	pet := &Pet{
		Id:     petID,
		Name:   petName,
		Status: Pet_STATUS_AVAILABLE,
	}
	_, err = testApp.AddPet(ctx, &AddPetRequest{Body: pet})
	if err != nil {
		t.Fatalf("Failed to add pet: %v", err)
	}

	// 2. Update Pet With Form
	updatedName := "UpdatedFormName"
	updatedStatus := "STATUS_SOLD"

	req := &UpdatePetWithFormRequest{
		PetId:  petID,
		Name:   updatedName,
		Status: updatedStatus,
	}

	_, err = testApp.UpdatePetWithForm(ctx, req)
	if err != nil {
		t.Fatalf("UpdatePetWithForm failed: %v", err)
	}

	// 3. Verify Update
	retrievedPet, err := testApp.GetPetById(ctx, &GetPetByIdRequest{PetId: petID})
	if err != nil {
		t.Fatalf("GetPetById failed: %v", err)
	}

	if retrievedPet.Name != updatedName {
		t.Errorf("Expected updated name %s, got %s", updatedName, retrievedPet.Name)
	}
	if retrievedPet.Status != Pet_STATUS_SOLD {
		t.Errorf("Expected status SOLD, got %v", retrievedPet.Status)
	}
}

func TestUploadFile(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	t.Run("ValidUpload", func(t *testing.T) {
		req := &UploadFileRequest{
			PetId:              1001,
			AdditionalMetadata: "avatar image",
			File:               "binary_data_placeholder",
		}
		resp, err := testApp.UploadFile(ctx, req)
		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}
		if resp.Code != 200 {
			t.Errorf("Expected code 200, got %d", resp.Code)
		}
	})
}
