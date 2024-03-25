package main

import (
	"time"

	"github.com/google/uuid"
)

func NewCustomer(name, email string) *Customer {
	now := time.Now()
	return &Customer{
		ID:        uuid.New(),
		Name:      name,
		Email:     email,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

type Customer struct {
	ID        uuid.UUID `db:"id"`
	Name      string    `db:"name"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func NewCustomerAddress(customerID uuid.UUID, locationName, address string) *CustomerAddress {
	now := time.Now()
	return &CustomerAddress{
		ID:           uuid.New(),
		CustomerID:   customerID,
		LocationName: locationName,
		Address:      address,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type CustomerAddress struct {
	ID           uuid.UUID `db:"id"`
	CustomerID   uuid.UUID `db:"customer_id"`
	LocationName string    `db:"location_name"`
	Address      string    `db:"address"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func NewProduct(name, description string, price float32) *Product {
	now := time.Now()
	return &Product{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		Price:       price,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

type Product struct {
	ID          uuid.UUID `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	Price       float32   `db:"price"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
}

type CartItem struct {
	Product
	InsertedAt time.Time `db:"inserted_at"`
}
