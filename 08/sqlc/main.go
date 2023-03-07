package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"data-persistence/sqlc/eshop"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	dsn = "postgres://root:root@localhost:5432/data-persistence?sslmode=disable"
)

func createCustomer(ctx context.Context, query *eshop.Queries) (*eshop.Customer, error) {
	now := time.Now()
	createCustomerParams := eshop.CreateCustomerParams{
		ID:        uuid.New(),
		Name:      "John Doe",
		Email:     "john.doe@gmail.com",
		CreatedAt: now,
		UpdatedAt: now,
	}
	customer, err := query.CreateCustomer(ctx, createCustomerParams)
	if err != nil {
		return nil, err
	}

	createCustomerAddressParams := eshop.CreateCustomerAddressParams{
		CustomerID:   customer.ID,
		LocationName: "Work",
		Address:      "Rohanské nábř. 678/23, 186 00 Karlín",
		CreatedAt:    now,
		UpdatedAt:    now,
	}
	if err = query.CreateCustomerAddress(ctx, createCustomerAddressParams); err != nil {
		return nil, err
	}

	return &customer, nil
}

func createProducts(ctx context.Context, query *eshop.Queries) ([]eshop.Product, error) {
	now := time.Now()
	params := eshop.CreateProductParams{
		ID:          uuid.New(),
		Name:        "Lagkapten",
		Description: "Desk, dark grey/black, 140x60 cm",
		Price:       1499,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	lagkapten, err := query.CreateProduct(ctx, params)
	if err != nil {
		return nil, err
	}

	params = eshop.CreateProductParams{
		ID:          uuid.New(),
		Name:        "Strandmon",
		Description: "Wing chair, Nordvalla dark grey",
		Price:       5490,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	strandmon, err := query.CreateProduct(ctx, params)
	if err != nil {
		return nil, err
	}

	products := []eshop.Product{lagkapten, strandmon}
	return products, nil
}

func addProductsToCart(ctx context.Context, query *eshop.Queries, customer *eshop.Customer, products ...eshop.Product) error {
	insertedAt := time.Now()
	for _, v := range products {
		params := eshop.AddProductToCartParams{
			CustomerID: customer.ID,
			ProductID:  v.ID,
			InsertedAt: insertedAt,
		}
		if err := query.AddProductToCart(ctx, params); err != nil {
			return err
		}
	}
	return nil
}

func createData(ctx context.Context, db *sql.DB) (err error) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	query := eshop.New(db).WithTx(tx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				err = errors.Join(rErr, err)
				return
			}
		}
		if cErr := tx.Commit(); err != nil {
			err = errors.Join(cErr, err)
		}
	}()

	customer, err := createCustomer(ctx, query)
	if err != nil {
		return err
	}

	products, err := createProducts(ctx, query)
	if err != nil {
		return err
	}

	err = addProductsToCart(ctx, query, customer, products...)
	if err != nil {
		return err
	}

	return nil
}

func readData(ctx context.Context, db *sql.DB, l *zap.Logger) error {
	query := eshop.New(db)

	customers, err := query.ListCustomers(ctx)
	if err != nil {
		return err
	}
	l.With(zap.Any("customers", customers)).Info("listing customers")

	for _, c := range customers {
		customerAddresses, err := query.ListCustomerAddresses(ctx, c.ID)
		if err != nil {
			return err
		}
		l.With(zap.Any("customer_addresses", customerAddresses)).Info("listing customer addresses")
	}

	for _, c := range customers {
		cartProducts, err := query.ListCartItems(ctx, c.ID)
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

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		l.Fatal("opening db connection", zap.Error(err))
	}

	if err = createData(ctx, db); err != nil {
		l.Fatal("creating data", zap.Error(err))
	}

	if err = readData(ctx, db, l); err != nil {
		l.Fatal("reading data", zap.Error(err))
	}
}
