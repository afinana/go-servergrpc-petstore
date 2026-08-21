package petstore

import (
	"context"
	"io"
	"log"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func newTestAppWithMocks() (*Application, *MockPetRepository, *MockStoreRepository, *MockUserRepository) {
	infoLog := log.New(io.Discard, "INFO\t", 0)
	errLog := log.New(io.Discard, "ERROR\t", 0)

	pets := NewMockPetRepository()
	stores := NewMockStoreRepository()
	users := NewMockUserRepository()

	app := NewLog(infoLog, errLog, pets, stores, users)
	return app, pets, stores, users
}

// ---------------- Pet Handler Unit Tests ----------------

func TestUnit_PetHandlers(t *testing.T) {
	app, _, _, _ := newTestAppWithMocks()
	ctx := context.Background()

	t.Run("AddPet and GetPetById", func(t *testing.T) {
		pet := &Pet{
			Id:        101,
			Name:      "Buddy",
			Status:    Pet_STATUS_AVAILABLE,
			PhotoUrls: []string{"http://example.com/buddy.jpg"},
			Category:  &Category{Id: 1, Name: "Dogs"},
		}

		_, err := app.AddPet(ctx, &AddPetRequest{Body: pet})
		if err != nil {
			t.Fatalf("AddPet failed: %v", err)
		}

		got, err := app.GetPetById(ctx, &GetPetByIdRequest{PetId: 101})
		if err != nil {
			t.Fatalf("GetPetById failed: %v", err)
		}
		if got.Name != "Buddy" || got.Id != 101 {
			t.Errorf("Expected pet Buddy/101, got %+v", got)
		}
	})

	t.Run("GetPetById_NotFound", func(t *testing.T) {
		_, err := app.GetPetById(ctx, &GetPetByIdRequest{PetId: 99999})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", status.Code(err))
		}
	})

	t.Run("FindPetsByStatus", func(t *testing.T) {
		_, _ = app.AddPet(ctx, &AddPetRequest{Body: &Pet{Id: 102, Name: "SoldPet", Status: Pet_STATUS_SOLD}})
		resp, err := app.FindPetsByStatus(ctx, &FindPetsByStatusRequest{
			Status: []FindPetsByStatusRequest_Status{FindPetsByStatusRequest_STATUS_SOLD},
		})
		if err != nil {
			t.Fatalf("FindPetsByStatus failed: %v", err)
		}
		if len(resp.Items) != 1 || resp.Items[0].Name != "SoldPet" {
			t.Errorf("Expected SoldPet, got %+v", resp.Items)
		}
	})

	t.Run("FindPetsByTags", func(t *testing.T) {
		_, _ = app.AddPet(ctx, &AddPetRequest{Body: &Pet{
			Id:   103,
			Name: "TaggedPet",
			Tags: []*Tag{{Id: 1, Name: "friendly"}},
		}})
		resp, err := app.FindPetsByTags(ctx, &FindPetsByTagsRequest{Tags: []string{"friendly"}})
		if err != nil {
			t.Fatalf("FindPetsByTags failed: %v", err)
		}
		if len(resp.Items) != 1 || resp.Items[0].Name != "TaggedPet" {
			t.Errorf("Expected TaggedPet, got %+v", resp.Items)
		}
	})

	t.Run("UpdatePet", func(t *testing.T) {
		_, err := app.UpdatePet(ctx, &UpdatePetRequest{
			Body: &Pet{Id: 101, Name: "BuddyRenamed", Status: Pet_STATUS_SOLD},
		})
		if err != nil {
			t.Fatalf("UpdatePet failed: %v", err)
		}

		got, err := app.GetPetById(ctx, &GetPetByIdRequest{PetId: 101})
		if err != nil {
			t.Fatalf("GetPetById failed: %v", err)
		}
		if got.Name != "BuddyRenamed" {
			t.Errorf("Expected BuddyRenamed, got %s", got.Name)
		}
	})

	t.Run("UpdatePetWithForm", func(t *testing.T) {
		_, err := app.UpdatePetWithForm(ctx, &UpdatePetWithFormRequest{
			PetId:  101,
			Name:   "BuddyFormUpdated",
			Status: "STATUS_AVAILABLE",
		})
		if err != nil {
			t.Fatalf("UpdatePetWithForm failed: %v", err)
		}

		got, err := app.GetPetById(ctx, &GetPetByIdRequest{PetId: 101})
		if err != nil {
			t.Fatalf("GetPetById failed: %v", err)
		}
		if got.Name != "BuddyFormUpdated" {
			t.Errorf("Expected BuddyFormUpdated, got %s", got.Name)
		}
	})

	t.Run("UpdatePetWithForm_NotFound", func(t *testing.T) {
		_, err := app.UpdatePetWithForm(ctx, &UpdatePetWithFormRequest{
			PetId: 88888,
			Name:  "NonExistent",
		})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", status.Code(err))
		}
	})

	t.Run("UploadFile_Success", func(t *testing.T) {
		resp, err := app.UploadFile(ctx, &UploadFileRequest{
			PetId:              101,
			File:               "test_base64_data",
			AdditionalMetadata: "profile_pic",
		})
		if err != nil {
			t.Fatalf("UploadFile failed: %v", err)
		}
		if resp.Code != 200 {
			t.Errorf("Expected code 200, got %d", resp.Code)
		}
	})

	t.Run("UploadFile_MissingPetId", func(t *testing.T) {
		_, err := app.UploadFile(ctx, &UploadFileRequest{
			PetId: 0,
			File:  "data",
		})
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %v", status.Code(err))
		}
	})

	t.Run("UploadFile_MissingFile", func(t *testing.T) {
		_, err := app.UploadFile(ctx, &UploadFileRequest{
			PetId: 101,
			File:  "",
		})
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %v", status.Code(err))
		}
	})

	t.Run("DeletePet", func(t *testing.T) {
		_, err := app.DeletePet(ctx, &DeletePetRequest{PetId: 101})
		if err != nil {
			t.Fatalf("DeletePet failed: %v", err)
		}

		_, err = app.GetPetById(ctx, &GetPetByIdRequest{PetId: 101})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound after deletion, got %v", status.Code(err))
		}
	})
}

// ---------------- Store Handler Unit Tests ----------------

func TestUnit_StoreHandlers(t *testing.T) {
	app, _, _, _ := newTestAppWithMocks()
	ctx := context.Background()

	t.Run("PlaceOrder and GetOrderById", func(t *testing.T) {
		order := &Order{
			Id:       201,
			PetId:    101,
			Quantity: 2,
			ShipDate: time.Now().Format(time.RFC3339),
			Status:   Order_ORDER_STATUS_PLACED,
			Complete: true,
		}

		placed, err := app.PlaceOrder(ctx, &PlaceOrderRequest{Body: order})
		if err != nil {
			t.Fatalf("PlaceOrder failed: %v", err)
		}
		if placed.Id != 201 {
			t.Errorf("Expected id 201, got %d", placed.Id)
		}

		got, err := app.GetOrderById(ctx, &GetOrderByIdRequest{OrderId: 201})
		if err != nil {
			t.Fatalf("GetOrderById failed: %v", err)
		}
		if got.Quantity != 2 {
			t.Errorf("Expected quantity 2, got %d", got.Quantity)
		}
	})

	t.Run("PlaceOrder_NilBody", func(t *testing.T) {
		_, err := app.PlaceOrder(ctx, &PlaceOrderRequest{Body: nil})
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %v", status.Code(err))
		}
	})

	t.Run("GetOrderById_NotFound", func(t *testing.T) {
		_, err := app.GetOrderById(ctx, &GetOrderByIdRequest{OrderId: 9999})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", status.Code(err))
		}
	})

	t.Run("GetInventory", func(t *testing.T) {
		inv, err := app.GetInventory(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("GetInventory failed: %v", err)
		}
		if len(inv.Items) == 0 {
			t.Errorf("Expected inventory items, got empty")
		}
	})

	t.Run("DeleteOrder", func(t *testing.T) {
		_, err := app.DeleteOrder(ctx, &DeleteOrderRequest{OrderId: 201})
		if err != nil {
			t.Fatalf("DeleteOrder failed: %v", err)
		}

		_, err = app.GetOrderById(ctx, &GetOrderByIdRequest{OrderId: 201})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound after delete, got %v", status.Code(err))
		}
	})
}

// ---------------- User Handler Unit Tests ----------------

func TestUnit_UserHandlers(t *testing.T) {
	app, _, _, _ := newTestAppWithMocks()
	ctx := context.Background()

	t.Run("CreateUser and GetUserByName", func(t *testing.T) {
		user := &User{
			Id:         301,
			Username:   "alice",
			FirstName:  "Alice",
			LastName:   "Smith",
			Email:      "alice@example.com",
			Password:   "p@ssword",
			Phone:      "1234567890",
			UserStatus: 1,
		}

		_, err := app.CreateUser(ctx, &CreateUserRequest{Body: user})
		if err != nil {
			t.Fatalf("CreateUser failed: %v", err)
		}

		got, err := app.GetUserByName(ctx, &GetUserByNameRequest{Username: "alice"})
		if err != nil {
			t.Fatalf("GetUserByName failed: %v", err)
		}
		if got.Email != "alice@example.com" {
			t.Errorf("Expected email alice@example.com, got %s", got.Email)
		}
	})

	t.Run("CreateUser_NilBody", func(t *testing.T) {
		_, err := app.CreateUser(ctx, &CreateUserRequest{Body: nil})
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %v", status.Code(err))
		}
	})

	t.Run("GetUserByName_NotFound", func(t *testing.T) {
		_, err := app.GetUserByName(ctx, &GetUserByNameRequest{Username: "nonexistent"})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound, got %v", status.Code(err))
		}
	})

	t.Run("LoginUser_Success", func(t *testing.T) {
		resp, err := app.LoginUser(ctx, &LoginUserRequest{Username: "alice", Password: "p@ssword"})
		if err != nil {
			t.Fatalf("LoginUser failed: %v", err)
		}
		if resp.Code != 200 {
			t.Errorf("Expected code 200, got %d", resp.Code)
		}
	})

	t.Run("LoginUser_InvalidCredentials", func(t *testing.T) {
		_, err := app.LoginUser(ctx, &LoginUserRequest{Username: "alice", Password: "wrong_password"})
		if status.Code(err) != codes.Unauthenticated {
			t.Errorf("Expected Unauthenticated, got %v", status.Code(err))
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		_, err := app.UpdateUser(ctx, &UpdateUserRequest{
			Username: "alice",
			Body: &User{
				Username:  "alice",
				FirstName: "AliceUpdated",
				LastName:  "Smith",
				Email:     "alice.updated@example.com",
			},
		})
		if err != nil {
			t.Fatalf("UpdateUser failed: %v", err)
		}

		got, err := app.GetUserByName(ctx, &GetUserByNameRequest{Username: "alice"})
		if err != nil {
			t.Fatalf("GetUserByName failed: %v", err)
		}
		if got.FirstName != "AliceUpdated" {
			t.Errorf("Expected FirstName AliceUpdated, got %s", got.FirstName)
		}
	})

	t.Run("UpdateUser_NilBody", func(t *testing.T) {
		_, err := app.UpdateUser(ctx, &UpdateUserRequest{Username: "alice", Body: nil})
		if status.Code(err) != codes.InvalidArgument {
			t.Errorf("Expected InvalidArgument, got %v", status.Code(err))
		}
	})

	t.Run("CreateUsersWithArrayInput", func(t *testing.T) {
		batch := []*User{
			{Id: 302, Username: "bob", Email: "bob@example.com"},
			{Id: 303, Username: "charlie", Email: "charlie@example.com"},
		}
		_, err := app.CreateUsersWithArrayInput(ctx, &CreateUsersWithArrayInputRequest{Body: batch})
		if err != nil {
			t.Fatalf("CreateUsersWithArrayInput failed: %v", err)
		}

		got, err := app.GetUserByName(ctx, &GetUserByNameRequest{Username: "bob"})
		if err != nil || got.Username != "bob" {
			t.Errorf("Expected bob, got %+v (err: %v)", got, err)
		}
	})

	t.Run("CreateUsersWithListInput", func(t *testing.T) {
		batch := []*User{
			{Id: 304, Username: "david", Email: "david@example.com"},
		}
		_, err := app.CreateUsersWithListInput(ctx, &CreateUsersWithListInputRequest{Body: batch})
		if err != nil {
			t.Fatalf("CreateUsersWithListInput failed: %v", err)
		}

		got, err := app.GetUserByName(ctx, &GetUserByNameRequest{Username: "david"})
		if err != nil || got.Username != "david" {
			t.Errorf("Expected david, got %+v (err: %v)", got, err)
		}
	})

	t.Run("LogoutUser", func(t *testing.T) {
		_, err := app.LogoutUser(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("LogoutUser failed: %v", err)
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		_, err := app.DeleteUser(ctx, &DeleteUserRequest{Username: "alice"})
		if err != nil {
			t.Fatalf("DeleteUser failed: %v", err)
		}

		_, err = app.GetUserByName(ctx, &GetUserByNameRequest{Username: "alice"})
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound after delete, got %v", status.Code(err))
		}
	})
}
