package petstore

import (
	"context"
	"testing"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestUserLifecycle(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()
	username := "testuser_lifecycle"
	user := &User{
		Id:         1001,
		Username:   username,
		FirstName:  "Test",
		LastName:   "User",
		Email:      "test@example.com",
		Password:   "password123",
		Phone:      "1234567890",
		UserStatus: 1,
	}

	// 1. Create User
	t.Run("CreateUser", func(t *testing.T) {
		req := &CreateUserRequest{Body: user}
		_, err := testApp.CreateUser(ctx, req)
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}
	})

	// 2. Get User By Name
	t.Run("GetUserByName", func(t *testing.T) {
		req := &GetUserByNameRequest{Username: username}
		resp, err := testApp.GetUserByName(ctx, req)
		if err != nil {
			t.Fatalf("GetUserByName failed: %v", err)
		}
		if resp.Username != username {
			t.Errorf("Expected username %s, got %s", username, resp.Username)
		}
		if resp.Email != user.Email {
			t.Errorf("Expected email %s, got %s", user.Email, resp.Email)
		}
	})

	// 3. Login User
	t.Run("LoginUser", func(t *testing.T) {
		req := &LoginUserRequest{Username: username, Password: user.Password}
		resp, err := testApp.LoginUser(ctx, req)
		if err != nil {
			t.Fatalf("LoginUser failed: %v", err)
		}
		if resp.Code != 200 {
			t.Errorf("Expected code 200, got %d", resp.Code)
		}
	})

	// 4. Update User
	t.Run("UpdateUser", func(t *testing.T) {
		user.FirstName = "Updated"
		req := &UpdateUserRequest{Username: username, Body: user}
		_, err := testApp.UpdateUser(ctx, req)
		if err != nil {
			t.Fatalf("UpdateUser failed: %v", err)
		}

		// Verify update in DB
		getReq := &GetUserByNameRequest{Username: username}
		resp, err := testApp.GetUserByName(ctx, getReq)
		if err != nil {
			t.Fatalf("GetUserByName after update failed: %v", err)
		}
		if resp.FirstName != "Updated" {
			t.Errorf("Expected FirstName Updated, got %s", resp.FirstName)
		}
	})

	// 5. Logout User
	t.Run("LogoutUser", func(t *testing.T) {
		req := &emptypb.Empty{}
		_, err := testApp.LogoutUser(ctx, req)
		if err != nil {
			t.Fatalf("LogoutUser failed: %v", err)
		}
	})

	// 6. Delete User
	t.Run("DeleteUser", func(t *testing.T) {
		req := &DeleteUserRequest{Username: username}
		_, err := testApp.DeleteUser(ctx, req)
		if err != nil {
			t.Fatalf("DeleteUser failed: %v", err)
		}

		// Verify deletion
		getReq := &GetUserByNameRequest{Username: username}
		_, err = testApp.GetUserByName(ctx, getReq)
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound error for deleted user, got %v", err)
		}
	})
}

func TestCreateUsersWithArrayInput(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	users := []*User{
		{Id: 2001, Username: "batch_array_1", Email: "b1@example.com"},
		{Id: 2002, Username: "batch_array_2", Email: "b2@example.com"},
	}

	req := &CreateUsersWithArrayInputRequest{Body: users}
	_, err := testApp.CreateUsersWithArrayInput(ctx, req)
	if err != nil {
		t.Fatalf("CreateUsersWithArrayInput failed: %v", err)
	}

	// Verify retrieval
	u1, err := testApp.GetUserByName(ctx, &GetUserByNameRequest{Username: "batch_array_1"})
	if err != nil {
		t.Fatalf("Failed to retrieve batch_array_1: %v", err)
	}
	if u1.Email != "b1@example.com" {
		t.Errorf("Expected email 'b1@example.com', got '%s'", u1.Email)
	}
}

func TestCreateUsersWithListInput(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	users := []*User{
		{Id: 3001, Username: "batch_list_1", Email: "l1@example.com"},
	}

	req := &CreateUsersWithListInputRequest{Body: users}
	_, err := testApp.CreateUsersWithListInput(ctx, req)
	if err != nil {
		t.Fatalf("CreateUsersWithListInput failed: %v", err)
	}

	u1, err := testApp.GetUserByName(ctx, &GetUserByNameRequest{Username: "batch_list_1"})
	if err != nil {
		t.Fatalf("Failed to retrieve batch_list_1: %v", err)
	}
	if u1.Email != "l1@example.com" {
		t.Errorf("Expected email 'l1@example.com', got '%s'", u1.Email)
	}
}

func TestLoginInvalidUser(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()
	req := &LoginUserRequest{Username: "nonexistent", Password: "password"}
	_, err := testApp.LoginUser(ctx, req)
	if status.Code(err) != codes.Unauthenticated {
		t.Errorf("Expected Unauthenticated error, got %v", err)
	}
}

func TestCreateUser_NilBody(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	req := &CreateUserRequest{Body: nil}
	_, err := testApp.CreateUser(ctx, req)
	if err == nil {
		t.Fatalf("Expected error for nil body, got nil")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument code, got %v", status.Code(err))
	}
}

func TestGetUserByName_NotFound(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	req := &GetUserByNameRequest{Username: "user_does_not_exist_404"}
	_, err := testApp.GetUserByName(ctx, req)
	if err == nil {
		t.Fatalf("Expected error for nonexistent user, got nil")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound code, got %v", status.Code(err))
	}
}
