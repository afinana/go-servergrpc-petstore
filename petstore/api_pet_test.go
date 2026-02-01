package petstore

import (
	"context"
	"testing"
)

func TestAddPetAndGetPetById(t *testing.T) {
	ctx := context.Background()

	// Clean up before test
	rdb.FlushDB(ctx)

	// 1. Add Pet
	petID := int64(1001)

	petName := "Fluffy"
	pet := &Pet{
		Id:        petID,
		Name:      petName,
		Status:    Pet_STATUS_AVAILABLE,
		PhotoUrls: []string{"url1", "url2"},
		Tags:      []*Tag{{Id: 1, Name: "cute"}},
	}
	addReq := &AddPetRequest{
		Body: pet,
	}

	_, err := testApp.AddPet(ctx, addReq)
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
}

func TestFindPetsByStatus(t *testing.T) {
	ctx := context.Background()

	// Clean up before test
	rdb.FlushDB(ctx)

	// Add a few pets with different statuses
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

func TestDeletePet(t *testing.T) {
	ctx := context.Background()

	// Clean up before test
	rdb.FlushDB(ctx)

	// Add a pet
	petID := int64(3001)
	pet := &Pet{
		Id:   petID,
		Name: "ToDelete",
	}
	_, err := testApp.AddPet(ctx, &AddPetRequest{Body: pet})
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
	// We might expect an error or nil here depending on implementation details
	// Looking at api_pet.go, it logs "Pets not found" and returns error if ErrNotFound
	_, err = testApp.GetPetById(ctx, getReq)
	if err == nil {
		t.Errorf("Expected error when getting deleted pet, got nil")
	}
}

func TestUpdatePetWithForm(t *testing.T) {
	ctx := context.Background()

	// Clean up before test
	rdb.FlushDB(ctx)

	// 1. Add Pet
	petID := int64(4001)
	petName := "OriginalName"
	pet := &Pet{
		Id:     petID,
		Name:   petName,
		Status: Pet_STATUS_AVAILABLE,
	}
	_, err := testApp.AddPet(ctx, &AddPetRequest{Body: pet})
	if err != nil {
		t.Fatalf("Failed to add pet: %v", err)
	}

	// 2. Update Pet With Form
	updatedName := "UpdatedName"
	updatedStatus := "STATUS_SOLD" // string status as per proto definition for UpdatePetWithForm

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
	getReq := &GetPetByIdRequest{
		PetId: petID,
	}
	retrievedPet, err := testApp.GetPetById(ctx, getReq)
	if err != nil {
		t.Fatalf("GetPetById failed: %v", err)
	}

	if retrievedPet.Name != updatedName {
		t.Errorf("Expected updated name %s, got %s", updatedName, retrievedPet.Name)
	}

	// API GetPetById returns Pet struct which has Status enum.
	// Our UpdatePetWithForm updates the DB with the string "sold".
	// The retrieval maps DB string to enum.
	// "sold" maps to Pet_PET_STATUS_SOLD (confirmed in pet_mapper.go usually)
	// Let's verify if the status enum matches.
	if retrievedPet.Status != Pet_STATUS_SOLD {
		t.Errorf("Expected status SOLD, got %v", retrievedPet.Status)
	}
}
