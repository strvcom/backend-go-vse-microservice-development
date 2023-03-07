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
	ID        uuid.UUID `gorm:"primaryKey"`
	Name      string    `gorm:"size:255;not null"`
	Email     string    `gorm:"size:255;not null,unique"`
	CreatedAt time.Time `gorm:"not null"`
	UpdatedAt time.Time `gorm:"not null"`
	Addresses []CustomerAddress
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
	CustomerID   uuid.UUID `gorm:"primaryKey"`
	LocationName string    `gorm:"size:255;not null"`
	Address      string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"not null"`
	UpdatedAt    time.Time `gorm:"not null"`
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
	ID          uuid.UUID `gorm:"primaryKey"`
	Name        string    `gorm:"size:255;not null"`
	Description string    `gorm:"size:1024;not null"`
	Price       float32   `gorm:"type:float;not null"`
	CreatedAt   time.Time `gorm:"not null"`
	UpdatedAt   time.Time `gorm:"not null"`
}

type CartItem struct {
	CustomerID uuid.UUID `gorm:"primaryKey;not null"`
	InsertedAt time.Time `gorm:"not null"`
	ProductID  uuid.UUID `gorm:"primaryKey;not null"`
	Product    Product
}
