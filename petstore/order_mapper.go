package petstore

import (
	"time"
)

// Converts Order DTO to OrderEntity
func createOrderEntity(in *Order) *OrderEntity {
	shipDate, _ := time.Parse(time.RFC3339, in.ShipDate)
	return &OrderEntity{
		Id:       in.Id,
		PetId:    in.PetId,
		Quantity: in.Quantity,
		ShipDate: shipDate,
		Status:   in.Status.String(),
		Complete: in.Complete,
	}
}

// Converts OrderEntity to Order DTO
func createOrderDTO(in *OrderEntity) *Order {
	var status Order_OrderStatus
	if s, ok := Order_OrderStatus_value[in.Status]; ok {
		status = Order_OrderStatus(s)
	}

	return &Order{
		Id:       in.Id,
		PetId:    in.PetId,
		Quantity: in.Quantity,
		ShipDate: in.ShipDate.Format(time.RFC3339),
		Status:   status,
		Complete: in.Complete,
	}
}
