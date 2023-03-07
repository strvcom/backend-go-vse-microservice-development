package repository

import (
	"time"

	"data-persistence/generics/model"

	"github.com/google/uuid"
)

type customer struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (c customer) toCustomer() *model.Customer {
	return &model.Customer{
		ID:        c.ID,
		Name:      c.Name,
		Email:     c.Email,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func toCustomers(customers []customer) []model.Customer {
	result := make([]model.Customer, 0, len(customers))
	for _, c := range customers {
		result = append(result, *c.toCustomer())
	}
	return result
}

type customerAddress struct {
	CustomerID   uuid.UUID `db:"customer_id"`
	LocationName string    `db:"location_name"`
	Address      string    `db:"address"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func (c customerAddress) toCustomerAddress() *model.CustomerAddress {
	return &model.CustomerAddress{
		CustomerID:   c.CustomerID,
		LocationName: c.LocationName,
		Address:      c.Address,
		CreatedAt:    c.CreatedAt,
		UpdatedAt:    c.UpdatedAt,
	}
}

func toCustomerAddresses(addresses []customerAddress) []model.CustomerAddress {
	result := make([]model.CustomerAddress, 0, len(addresses))
	for _, c := range addresses {
		result = append(result, *c.toCustomerAddress())
	}
	return result
}

type product struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float32   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type cartItem struct {
	product
	InsertedAt time.Time `db:"inserted_at"`
}

func (c cartItem) toCartItem() *model.CartItem {
	return &model.CartItem{
		Product: model.Product{
			ID:          c.ID,
			Name:        c.Name,
			Description: c.Description,
			Price:       c.Price,
			CreatedAt:   c.CreatedAt,
			UpdatedAt:   c.UpdatedAt,
		},
		InsertedAt: c.InsertedAt,
	}
}

func toCartItems(cartItems []cartItem) []model.CartItem {
	result := make([]model.CartItem, 0, len(cartItems))
	for _, c := range cartItems {
		result = append(result, *c.toCartItem())
	}
	return result
}
