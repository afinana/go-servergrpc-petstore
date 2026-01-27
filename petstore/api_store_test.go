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
		_, err := testApp.AddPet(ctx, &AddPetRequest{Body: &Pet{Id: 9999, Status: Pet_PET_STATUS_AVAILABLE}})
		if err != nil {
			t.Fatalf("AddPet for inventory test failed: %v", err)
		}

		resp, err := testApp.GetInventory(ctx, &emptypb.Empty{})
		if err != nil {
			t.Fatalf("GetInventory failed: %v", err)
		}

		found := false
		for _, item := range resp.Items {
			if item.Name == "PET_STATUS_AVAILABLE" && item.Value > 0 {
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
