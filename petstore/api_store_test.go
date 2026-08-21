package petstore

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func TestStoreLifecycle(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()
	orderID := int64(1001)
	shipDate := time.Now().Truncate(time.Second)

	order := &Order{
		Id:       orderID,
		PetId:    5001,
		Quantity: 3,
		ShipDate: shipDate.Format(time.RFC3339),
		Status:   Order_ORDER_STATUS_PLACED,
		Complete: false,
	}

	// 1. Place Order
	t.Run("PlaceOrder", func(t *testing.T) {
		req := &PlaceOrderRequest{Body: order}
		resp, err := testApp.PlaceOrder(ctx, req)
		if err != nil {
			t.Fatalf("PlaceOrder failed: %v", err)
		}
		if resp.Id != orderID {
			t.Errorf("Expected id %d, got %d", orderID, resp.Id)
		}
	})

	// 2. Get Order By ID
	t.Run("GetOrderById", func(t *testing.T) {
		req := &GetOrderByIdRequest{OrderId: orderID}
		resp, err := testApp.GetOrderById(ctx, req)
		if err != nil {
			t.Fatalf("GetOrderById failed: %v", err)
		}
		if resp.Id != orderID {
			t.Errorf("Expected id %d, got %d", orderID, resp.Id)
		}
		if resp.Quantity != 3 {
			t.Errorf("Expected quantity 3, got %d", resp.Quantity)
		}
	})

	// 3. Get Inventory
	t.Run("GetInventory", func(t *testing.T) {
		// Add a pet to ensure inventory has something
		_, err := testApp.AddPet(ctx, &AddPetRequest{Body: &Pet{Id: 9999, Status: Pet_STATUS_AVAILABLE}})
		if err != nil {
			t.Fatalf("AddPet for inventory test failed: %v", err)
		}

		resp, err := testApp.GetInventory(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("GetInventory failed: %v", err)
		}

		found := false
		for _, item := range resp.Items {
			if item.Name == "STATUS_AVAILABLE" && item.Value > 0 {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected available pets in inventory, items: %v", resp.Items)
		}
	})

	// 4. Delete Order
	t.Run("DeleteOrder", func(t *testing.T) {
		req := &DeleteOrderRequest{OrderId: orderID}
		_, err := testApp.DeleteOrder(ctx, req)
		if err != nil {
			t.Fatalf("DeleteOrder failed: %v", err)
		}

		// Verify deletion
		getReq := &GetOrderByIdRequest{OrderId: orderID}
		_, err = testApp.GetOrderById(ctx, getReq)
		if status.Code(err) != codes.NotFound {
			t.Errorf("Expected NotFound error, got %v", err)
		}
	})
}

func TestPlaceOrder_NilBody(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	req := &PlaceOrderRequest{Body: nil}
	_, err := testApp.PlaceOrder(ctx, req)
	if err == nil {
		t.Fatalf("Expected error for nil body, got nil")
	}
	if status.Code(err) != codes.InvalidArgument {
		t.Errorf("Expected InvalidArgument code, got %v", status.Code(err))
	}
}

func TestGetOrderById_NotFound(t *testing.T) {
	skipIfNoMongo(t)
	ctx := context.Background()

	req := &GetOrderByIdRequest{OrderId: 888888}
	_, err := testApp.GetOrderById(ctx, req)
	if err == nil {
		t.Fatalf("Expected error for nonexistent order, got nil")
	}
	if status.Code(err) != codes.NotFound {
		t.Errorf("Expected NotFound code, got %v", status.Code(err))
	}
}
