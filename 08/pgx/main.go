package main

import (
	"context"
	"errors"
	"time"

	"github.com/georgysavva/scany/v2/pgxscan"

	"data-persistence/pgx/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	dsn = "postgres://root:root@localhost:5432/data-persistence?sslmode=disable"
)

type Querier interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, optionsAndArgs ...any) (pgx.Rows, error)
	SendBatch(ctx context.Context, b *pgx.Batch) pgx.BatchResults
}

func createCustomer(ctx context.Context, querier Querier) (*Customer, error) {
	customer := NewCustomer("John Doe", "john.doe@gmail.com")
	_, err := querier.Exec(ctx, query.CreateCustomer, pgx.NamedArgs{
		"id":         customer.ID,
		"name":       customer.Name,
		"email":      customer.Email,
		"created_at": customer.CreatedAt,
		"updated_at": customer.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	customerAddress := NewCustomerAddress(
		customer.ID,
		"Work",
		"Rohanské nábř. 678/23, 186 00 Karlín",
	)
	_, err = querier.Exec(ctx, query.CreateCustomerAddress, pgx.NamedArgs{
		"id":            customerAddress.ID,
		"customer_id":   customerAddress.CustomerID,
		"location_name": customerAddress.LocationName,
		"address":       customerAddress.Address,
		"created_at":    customerAddress.CreatedAt,
		"updated_at":    customerAddress.UpdatedAt,
	})
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func createProducts(ctx context.Context, querier Querier) ([]Product, error) {
	products := []Product{
		*NewProduct("Lagkapten", "Desk, dark grey/black, 140x60 cm", 1499),
		*NewProduct("Strandmon", "Wing chair, Nordvalla dark grey", 5490),
	}

	batch := &pgx.Batch{}
	for _, v := range products {
		batch.Queue(query.CreateProduct, pgx.NamedArgs{
			"id":          v.ID,
			"name":        v.Name,
			"description": v.Description,
			"price":       v.Price,
			"created_at":  v.CreatedAt,
			"updated_at":  v.UpdatedAt,
		})
	}

	result := querier.SendBatch(ctx, batch)
	for i := 0; i < len(products); i++ {
		if _, err := result.Exec(); err != nil {
			return nil, err
		}
	}

	if err := result.Close(); err != nil {
		return nil, err
	}

	return products, nil
}

func addProductsToCart(ctx context.Context, querier Querier, customer *Customer, products []Product) error {
	insertedAt := time.Now()
	batch := &pgx.Batch{}
	for _, v := range products {
		batch.Queue(query.AddProductToCart, pgx.NamedArgs{
			"customer_id": customer.ID,
			"product_id":  v.ID,
			"inserted_at": insertedAt,
		})
	}

	result := querier.SendBatch(ctx, batch)
	for i := 0; i < len(products); i++ {
		if _, err := result.Exec(); err != nil {
			return err
		}
	}

	if err := result.Close(); err != nil {
		return err
	}

	return nil
}

func createData(ctx context.Context, db *pgxpool.Pool) (err error) {
	tx, err := db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(ctx); rErr != nil {
				err = errors.Join(rErr, err)
			}
			return
		}
		if cErr := tx.Commit(ctx); err != nil {
			err = errors.Join(cErr, err)
		}
	}()

	customer, err := createCustomer(ctx, tx)
	if err != nil {
		return err
	}
	products, err := createProducts(ctx, tx)
	if err != nil {
		return err
	}
	err = addProductsToCart(ctx, tx, customer, products)
	if err != nil {
		return err
	}

	return nil
}

func listCustomers(ctx context.Context, querier Querier) ([]Customer, error) {
	var customers []Customer
	if err := pgxscan.Select(ctx, querier, &customers, query.ListCostumers); err != nil {
		return nil, err
	}
	return customers, nil
}

func listCustomerAddresses(ctx context.Context, querier Querier, customerID uuid.UUID) ([]CustomerAddress, error) {
	var customerAddresses []CustomerAddress
	err := pgxscan.Select(ctx, querier, &customerAddresses, query.ListCustomerAddresses, pgx.NamedArgs{
		"customer_id": customerID,
	})
	if err != nil {
		return nil, err
	}
	return customerAddresses, nil
}

func listCartItems(ctx context.Context, querier Querier, customerID uuid.UUID) ([]CartItem, error) {
	var cartProducts []CartItem
	err := pgxscan.Select(ctx, querier, &cartProducts, query.ListCartItems, pgx.NamedArgs{
		"customer_id": customerID,
	})
	if err != nil {
		return nil, err
	}
	return cartProducts, nil
}

func readData(ctx context.Context, db *pgxpool.Pool, l *zap.Logger) error {
	customers, err := listCustomers(ctx, db)
	if err != nil {
		return err
	}
	l.With(zap.Any("customers", customers)).Info("listing customers")

	for _, c := range customers {
		customerAddresses, err := listCustomerAddresses(ctx, db, c.ID)
		if err != nil {
			return err
		}
		l.With(zap.Any("customer_addresses", customerAddresses)).Info("listing customer addresses")
	}

	for _, c := range customers {
		cartProducts, err := listCartItems(ctx, db, c.ID)
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

	db, err := pgxpool.New(ctx, dsn)
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
