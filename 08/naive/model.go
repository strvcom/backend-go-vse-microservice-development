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
	ID        uuid.UUID
	Name      string
	Email     string
	CreatedAt time.Time
	UpdatedAt time.Time
}

func NewCustomerAddress(customerID uuid.UUID, locationName, address string) *CustomerAddress {
	now := time.Now()
	return &CustomerAddress{
		CustomerID:   customerID,
		LocationName: locationName,
		Address:      address,
		CreatedAt:    now,
		UpdatedAt:    now,
	}
}

type CustomerAddress struct {
	CustomerID   uuid.UUID
	LocationName string
	Address      string
	CreatedAt    time.Time
	UpdatedAt    time.Time
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
	ID          uuid.UUID
	Name        string
	Description string
	Price       float32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CartItem struct {
	Product
	InsertedAt time.Time
}
