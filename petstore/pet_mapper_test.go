package petstore

import (
	"testing"
)

func TestPetMapper_createPetEntity(t *testing.T) {
	t.Run("FullPetMapping", func(t *testing.T) {
		dto := &Pet{
			Id:   100,
			Name: "Max",
			Category: &Category{
				Id:   1,
				Name: "Dogs",
			},
			PhotoUrls: []string{"http://example.com/dog1.jpg", "http://example.com/dog2.jpg"},
			Tags: []*Tag{
				{Id: 10, Name: "friendly"},
				{Id: 20, Name: "vaccinated"},
			},
			Status: Pet_STATUS_AVAILABLE,
		}

		entity := createPetEntity(dto)

		if entity.Id != 100 {
			t.Errorf("Expected Id 100, got %d", entity.Id)
		}
		if entity.Name != "Max" {
			t.Errorf("Expected Name 'Max', got '%s'", entity.Name)
		}
		if entity.Category == nil || entity.Category.Id != 1 || entity.Category.Name != "Dogs" {
			t.Errorf("Category mapped incorrectly: %+v", entity.Category)
		}
		if len(entity.PhotoUrls) != 2 {
			t.Errorf("Expected 2 photo URLs, got %d", len(entity.PhotoUrls))
		}
		if len(entity.Tags) != 2 || entity.Tags[0].Name != "friendly" {
			t.Errorf("Tags mapped incorrectly: %+v", entity.Tags)
		}
		if entity.Status != "STATUS_AVAILABLE" {
			t.Errorf("Expected Status 'STATUS_AVAILABLE', got '%s'", entity.Status)
		}
	})

	t.Run("NilCategoryAndEmptyCollections", func(t *testing.T) {
		dto := &Pet{
			Id:     200,
			Name:   "Shadow",
			Status: Pet_STATUS_PENDING,
		}

		entity := createPetEntity(dto)

		if entity.Category != nil {
			t.Errorf("Expected nil category, got %+v", entity.Category)
		}
		if len(entity.Tags) != 0 {
			t.Errorf("Expected empty tags, got %d", len(entity.Tags))
		}
		if len(entity.PhotoUrls) != 0 {
			t.Errorf("Expected empty photo URLs, got %d", len(entity.PhotoUrls))
		}
		if entity.Status != "STATUS_PENDING" {
			t.Errorf("Expected Status 'STATUS_PENDING', got '%s'", entity.Status)
		}
	})
}

func TestPetMapper_createPetDTO(t *testing.T) {
	entity := &PetEntity{
		Id:   300,
		Name: "Bella",
		Category: &CategoryEntity{
			Id:   2,
			Name: "Cats",
		},
		PhotoUrls: []string{"http://example.com/cat.jpg"},
		Tags: []TagEntity{
			{Id: 30, Name: "playful"},
		},
		Status: "STATUS_SOLD",
	}

	dto := createPetDTO(entity)

	if dto.Id != 300 {
		t.Errorf("Expected Id 300, got %d", dto.Id)
	}
	if dto.Name != "Bella" {
		t.Errorf("Expected Name 'Bella', got '%s'", dto.Name)
	}
	if dto.Category == nil || dto.Category.Name != "Cats" {
		t.Errorf("Category mapped incorrectly: %+v", dto.Category)
	}
	if dto.Status != Pet_STATUS_SOLD {
		t.Errorf("Expected Pet_STATUS_SOLD, got %v", dto.Status)
	}
}

func TestPetMapper_convertPetStatusList(t *testing.T) {
	input := []FindPetsByStatusRequest_Status{
		FindPetsByStatusRequest_STATUS_AVAILABLE,
		FindPetsByStatusRequest_STATUS_SOLD,
	}

	output := convertPetStatusList(input)

	if len(output) != 2 {
		t.Fatalf("Expected 2 strings, got %d", len(output))
	}
	if output[0] != "STATUS_AVAILABLE" || output[1] != "STATUS_SOLD" {
		t.Errorf("Status list conversion incorrect: %v", output)
	}
}

func TestPetMapper_CreatePetListDTO(t *testing.T) {
	entities := []PetEntity{
		{Id: 1, Name: "Pet1", Status: "STATUS_AVAILABLE"},
		{Id: 2, Name: "Pet2", Status: "STATUS_SOLD"},
	}

	dtos := CreatePetListDTO(entities)

	if len(dtos) != 2 {
		t.Fatalf("Expected 2 DTOs, got %d", len(dtos))
	}
	if dtos[0].Name != "Pet1" || dtos[1].Name != "Pet2" {
		t.Errorf("List mapping incorrect: %+v", dtos)
	}
}
