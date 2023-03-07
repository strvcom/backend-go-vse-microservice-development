package main

import (
	"context"

	"github.com/google/uuid"

	"data-persistence/generics/model"
	"data-persistence/generics/repository"

	"go.uber.org/zap"
)

const (
	dsn = "postgres://root:root@localhost:5432/data-persistence?sslmode=disable"
)

func createCustomer(ctx context.Context, r repository.Repository) (*model.Customer, error) {
	customer := model.NewCustomer("John Doe", "john.doe@gmail.com")
	address := model.NewCustomerAddress(
		customer.ID,
		"Work",
		"Rohanské nábř. 678/23, 186 00 Karlín",
	)
	if err := r.CreateCustomerWithAddress(ctx, customer, address); err != nil {
		return nil, err
	}
	return customer, nil
}

func createProducts(ctx context.Context, r repository.Repository) ([]model.Product, error) {
	products := []model.Product{
		*model.NewProduct("Lagkapten", "Desk, dark grey/black, 140x60 cm", 1499),
		*model.NewProduct("Strandmon", "Wing chair, Nordvalla dark grey", 5490),
	}
	if err := r.CreateProducts(ctx, products); err != nil {
		return nil, err
	}
	return products, nil
}

func addProductsToCart(ctx context.Context, r repository.Repository, customerID uuid.UUID, productIDs []uuid.UUID) error {
	return r.AddProductsToCart(ctx, customerID, productIDs)
}

func createData(ctx context.Context, r repository.Repository) error {
	customer, err := createCustomer(ctx, r)
	if err != nil {
		return err
	}

	products, err := createProducts(ctx, r)
	if err != nil {
		return err
	}

	var productIDs []uuid.UUID
	for _, p := range products {
		productIDs = append(productIDs, p.ID)
	}
	err = addProductsToCart(ctx, r, customer.ID, productIDs)
	if err != nil {
		return err
	}

	return nil
}

func readData(ctx context.Context, r repository.Repository, l *zap.Logger) error {
	customers, err := r.ListCustomers(ctx)
	if err != nil {
		return err
	}
	l.With(zap.Any("customers", customers)).Info("listing customers")

	for _, c := range customers {
		customerAddresses, err := r.ListCustomerAddresses(ctx, c.ID)
		if err != nil {
			return err
		}
		l.With(zap.Any("customer_addresses", customerAddresses)).Info("listing customer addresses")
	}

	for _, c := range customers {
		cartProducts, err := r.ListCartItems(ctx, c.ID)
		if err != nil {
			return err
		}
		l.With(
			zap.String("customer_id", c.ID.String()),
			zap.Any("cart_items", cartProducts),
		).Info("listing products in cart")
	}

	return nil
}

func main() {
	ctx := context.Background()

	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	r, err := repository.New(ctx, dsn)
	if err != nil {
		l.Fatal("new repository", zap.Error(err))
	}

	if err = createData(ctx, r); err != nil {
		l.Fatal("creating data", zap.Error(err))
	}

	if err = readData(ctx, r, l); err != nil {
		l.Fatal("reading data", zap.Error(err))
	}
}
