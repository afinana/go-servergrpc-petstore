package petstore

import (
	"testing"
)

func TestUserMapper(t *testing.T) {
	t.Run("CreateUserEntity", func(t *testing.T) {
		dto := &User{
			Id:         10,
			Username:   "johndoe",
			FirstName:  "John",
			LastName:   "Doe",
			Email:      "john@example.com",
			Password:   "secret",
			Phone:      "555-1234",
			UserStatus: 1,
		}

		entity := createUserEntity(dto)

		if entity.Id != 10 {
			t.Errorf("Expected Id 10, got %d", entity.Id)
		}
		if entity.Username != "johndoe" {
			t.Errorf("Expected Username 'johndoe', got '%s'", entity.Username)
		}
		if entity.Email != "john@example.com" {
			t.Errorf("Expected Email 'john@example.com', got '%s'", entity.Email)
		}
	})

	t.Run("CreateUserEntity_Nil", func(t *testing.T) {
		entity := createUserEntity(nil)
		if entity != nil {
			t.Errorf("Expected nil entity, got %+v", entity)
		}
	})

	t.Run("CreateUserDTO", func(t *testing.T) {
		entity := &UserEntity{
			Id:         20,
			Username:   "janedoe",
			FirstName:  "Jane",
			LastName:   "Doe",
			Email:      "jane@example.com",
			Password:   "secret2",
			Phone:      "555-5678",
			UserStatus: 2,
		}

		dto := createUserDTO(entity)

		if dto.Id != 20 {
			t.Errorf("Expected Id 20, got %d", dto.Id)
		}
		if dto.Username != "janedoe" {
			t.Errorf("Expected Username 'janedoe', got '%s'", dto.Username)
		}
	})

	t.Run("CreateUserDTO_Nil", func(t *testing.T) {
		dto := createUserDTO(nil)
		if dto != nil {
			t.Errorf("Expected nil DTO, got %+v", dto)
		}
	})

	t.Run("CreateUserListDTO", func(t *testing.T) {
		entities := []UserEntity{
			{Id: 1, Username: "user1"},
			{Id: 2, Username: "user2"},
		}

		dtos := CreateUserListDTO(entities)

		if len(dtos) != 2 {
			t.Fatalf("Expected 2 DTOs, got %d", len(dtos))
		}
		if dtos[0].Username != "user1" || dtos[1].Username != "user2" {
			t.Errorf("List mapping incorrect: %+v", dtos)
		}
	})
}
