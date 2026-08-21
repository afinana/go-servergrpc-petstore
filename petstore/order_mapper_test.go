package petstore

import (
	"testing"
	"time"
)

func TestOrderMapper(t *testing.T) {
	now := time.Now().UTC().Truncate(time.Second)
	dateStr := now.Format(time.RFC3339)

	t.Run("CreateOrderEntity", func(t *testing.T) {
		dto := &Order{
			Id:       501,
			PetId:    101,
			Quantity: 2,
			ShipDate: dateStr,
			Status:   Order_ORDER_STATUS_PLACED,
			Complete: true,
		}

		entity := createOrderEntity(dto)

		if entity.Id != 501 {
			t.Errorf("Expected Id 501, got %d", entity.Id)
		}
		if entity.PetId != 101 {
			t.Errorf("Expected PetId 101, got %d", entity.PetId)
		}
		if entity.Quantity != 2 {
			t.Errorf("Expected Quantity 2, got %d", entity.Quantity)
		}
		if !entity.ShipDate.Equal(now) {
			t.Errorf("Expected ShipDate %v, got %v", now, entity.ShipDate)
		}
		if entity.Status != "ORDER_STATUS_PLACED" {
			t.Errorf("Expected Status 'ORDER_STATUS_PLACED', got '%s'", entity.Status)
		}
		if !entity.Complete {
			t.Errorf("Expected Complete to be true")
		}
	})

	t.Run("CreateOrderDTO", func(t *testing.T) {
		entity := &OrderEntity{
			Id:       502,
			PetId:    102,
			Quantity: 5,
			ShipDate: now,
			Status:   "ORDER_STATUS_DELIVERED",
			Complete: true,
		}

		dto := createOrderDTO(entity)

		if dto.Id != 502 {
			t.Errorf("Expected Id 502, got %d", dto.Id)
		}
		if dto.Status != Order_ORDER_STATUS_DELIVERED {
			t.Errorf("Expected ORDER_STATUS_DELIVERED, got %v", dto.Status)
		}
		if dto.ShipDate != dateStr {
			t.Errorf("Expected ShipDate '%s', got '%s'", dateStr, dto.ShipDate)
		}
	})
}
