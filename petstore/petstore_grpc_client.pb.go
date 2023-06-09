// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: petstore.pb

package petstore

import (
	context "context"
	grpc "google.golang.org/grpc"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// SwaggerPetstoreServiceClient is the client API for SwaggerPetstoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SwaggerPetstoreServiceClient interface {
	// Add a new pet to the store
	AddPet(ctx context.Context, in *AddPetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Create user
	//
	
	// This can only be done by the logged in user.
	CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Creates list of users with given input array
	CreateUsersWithArrayInput(ctx context.Context, in *CreateUsersWithArrayInputRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Creates list of users with given input array
	CreateUsersWithListInput(ctx context.Context, in *CreateUsersWithListInputRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete purchase order by ID
	//
	// For valid response try integer IDs with positive integer value. Negative or non-integer values will generate API errors
	DeleteOrder(ctx context.Context, in *DeleteOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Deletes a pet
	DeletePet(ctx context.Context, in *DeletePetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Delete user
	//
	// This can only be done by the logged in user.
	DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Finds Pets by status
	//
	// Multiple status values can be provided with comma separated strings
	FindPetsByStatus(ctx context.Context, in *FindPetsByStatusRequest, opts ...grpc.CallOption) (*FindPetsByStatusResponse, error)
	// Finds Pets by tags
	//
	// Multiple tags can be provided with comma separated strings. Use tag1, tag2, tag3 for testing.
	FindPetsByTags(ctx context.Context, in *FindPetsByTagsRequest, opts ...grpc.CallOption) (*FindPetsByTagsResponse, error)
	// Returns pet inventories by status
	//
	// Returns a map of status codes to quantities
	GetInventory(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetInventoryResponse, error)
	// Find purchase order by ID
	//
	// For valid response try integer IDs with value >= 1 and <= 10. Other values will generated exceptions
	GetOrderById(ctx context.Context, in *GetOrderByIdRequest, opts ...grpc.CallOption) (*Order, error)
	// Find pet by ID
	//
	// Returns a single pet
	GetPetById(ctx context.Context, in *GetPetByIdRequest, opts ...grpc.CallOption) (*Pet, error)
	// Get user by user name
	GetUserByName(ctx context.Context, in *GetUserByNameRequest, opts ...grpc.CallOption) (*User, error)
	// Logs user into the system
	LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*ApiResponse, error)
	// Logs out current logged in user session
	LogoutUser(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Place an order for a pet
	PlaceOrder(ctx context.Context, in *PlaceOrderRequest, opts ...grpc.CallOption) (*Order, error)
	// Update an existing pet
	UpdatePet(ctx context.Context, in *UpdatePetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Updates a pet in the store with form data
	UpdatePetWithForm(ctx context.Context, in *UpdatePetWithFormRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// Updated user
	//
	// This can only be done by the logged in user.
	UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	// uploads an image
	UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*ApiResponse, error)
}

type swaggerPetstoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewSwaggerPetstoreServiceClient(cc grpc.ClientConnInterface) SwaggerPetstoreServiceClient {
	return &swaggerPetstoreServiceClient{cc}
}

func (c *swaggerPetstoreServiceClient) AddPet(ctx context.Context, in *AddPetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/AddPet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) CreateUser(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/CreateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) CreateUsersWithArrayInput(ctx context.Context, in *CreateUsersWithArrayInputRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/CreateUsersWithArrayInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) CreateUsersWithListInput(ctx context.Context, in *CreateUsersWithListInputRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/CreateUsersWithListInput", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) DeleteOrder(ctx context.Context, in *DeleteOrderRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/DeleteOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) DeletePet(ctx context.Context, in *DeletePetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/DeletePet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) DeleteUser(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/DeleteUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) FindPetsByStatus(ctx context.Context, in *FindPetsByStatusRequest, opts ...grpc.CallOption) (*FindPetsByStatusResponse, error) {
	out := new(FindPetsByStatusResponse)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/FindPetsByStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) FindPetsByTags(ctx context.Context, in *FindPetsByTagsRequest, opts ...grpc.CallOption) (*FindPetsByTagsResponse, error) {
	out := new(FindPetsByTagsResponse)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/FindPetsByTags", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) GetInventory(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*GetInventoryResponse, error) {
	out := new(GetInventoryResponse)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/GetInventory", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) GetOrderById(ctx context.Context, in *GetOrderByIdRequest, opts ...grpc.CallOption) (*Order, error) {
	out := new(Order)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/GetOrderById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) GetPetById(ctx context.Context, in *GetPetByIdRequest, opts ...grpc.CallOption) (*Pet, error) {
	out := new(Pet)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/GetPetById", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) GetUserByName(ctx context.Context, in *GetUserByNameRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/GetUserByName", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) LoginUser(ctx context.Context, in *LoginUserRequest, opts ...grpc.CallOption) (*ApiResponse, error) {
	out := new(ApiResponse)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/LoginUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) LogoutUser(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/LogoutUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) PlaceOrder(ctx context.Context, in *PlaceOrderRequest, opts ...grpc.CallOption) (*Order, error) {
	out := new(Order)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/PlaceOrder", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) UpdatePet(ctx context.Context, in *UpdatePetRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/UpdatePet", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) UpdatePetWithForm(ctx context.Context, in *UpdatePetWithFormRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/UpdatePetWithForm", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) UpdateUser(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/UpdateUser", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *swaggerPetstoreServiceClient) UploadFile(ctx context.Context, in *UploadFileRequest, opts ...grpc.CallOption) (*ApiResponse, error) {
	out := new(ApiResponse)
	err := c.cc.Invoke(ctx, "/swaggerpetstore.SwaggerPetstoreService/UploadFile", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

