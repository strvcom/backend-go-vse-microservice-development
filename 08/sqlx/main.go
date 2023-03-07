package main

import (
	"context"
	"errors"
	"time"

	"data-persistence/sqlx/query"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	dsn = "postgres://root:root@localhost:5432/data-persistence?sslmode=disable"
)

func createCustomer(tx *sqlx.Tx) (*Customer, error) {
	customer := NewCustomer("John Doe", "john.doe@gmail.com")
	if _, err := tx.NamedExec(query.CreateCustomer, customer); err != nil {
		return nil, err
	}

	customerAddress := NewCustomerAddress(
		customer.ID,
		"Work",
		"Rohanské nábř. 678/23, 186 00 Karlín",
	)
	_, err := tx.NamedExec(
		query.CreateCustomerAddress,
		map[string]any{
			"customer_id":   customerAddress.CustomerID,
			"location_name": customerAddress.LocationName,
			"address":       customerAddress.Address,
			"created_at":    customerAddress.CreatedAt,
			"updated_at":    customer.UpdatedAt,
		},
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func createProducts(tx *sqlx.Tx) ([]Product, error) {
	products := []Product{
		*NewProduct("Lagkapten", "Desk, dark grey/black, 140x60 cm", 1499),
		*NewProduct("Strandmon", "Wing chair, Nordvalla dark grey", 5490),
	}
	if _, err := tx.NamedExec(query.CreateProduct, products); err != nil {
		return nil, err
	}
	return products, nil
}

func addProductsToCart(tx *sqlx.Tx, customer *Customer, products []Product) error {
	insertedAt := time.Now()
	var params []map[string]any
	for _, v := range products {
		params = append(params, map[string]any{
			"customer_id": customer.ID,
			"product_id":  v.ID,
			"inserted_at": insertedAt,
		})
	}
	if _, err := tx.NamedExec(query.AddProductToCart, params); err != nil {
		return err
	}
	return nil
}

func createData(ctx context.Context, db *sqlx.DB) (err error) {
	tx, err := db.BeginTxx(ctx, nil)
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

	customer, err := createCustomer(tx)
	if err != nil {
		return err
	}
	products, err := createProducts(tx)
	if err != nil {
		return err
	}
	err = addProductsToCart(tx, customer, products)
	if err != nil {
		return err
	}

	return nil
}

func listCustomers(db *sqlx.DB) ([]Customer, error) {
	var customers []Customer
	if err := db.Select(&customers, query.ListCostumers); err != nil {
		return nil, err
	}
	return customers, nil
}

func listCustomerAddresses(db *sqlx.DB, customerID uuid.UUID) ([]CustomerAddress, error) {
	var customerAddresses []CustomerAddress
	if err := db.Select(&customerAddresses, query.ListCustomerAddresses, customerID); err != nil {
		return nil, err
	}
	return customerAddresses, nil
}

func listCartItems(db *sqlx.DB, customerID uuid.UUID) ([]CartItem, error) {
	var cartProducts []CartItem
	if err := db.Select(&cartProducts, query.ListCartItems, customerID); err != nil {
		return nil, err
	}

	return cartProducts, nil
}

func readData(db *sqlx.DB, l *zap.Logger) error {
	customers, err := listCustomers(db)
	if err != nil {
		return err
	}
	l.With(zap.Any("customers", customers)).Info("listing customers")

	for _, c := range customers {
		customerAddresses, err := listCustomerAddresses(db, c.ID)
		if err != nil {
			return err
		}
		l.With(zap.Any("customer_addresses", customerAddresses)).Info("listing customer addresses")
	}

	for _, c := range customers {
		cartProducts, err := listCartItems(db, c.ID)
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
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db, err := sqlx.Open("pgx", dsn)
	if err != nil {
		l.Fatal("opening db connection", zap.Error(err))
	}

	if err = createData(context.Background(), db); err != nil {
		l.Fatal("creating data", zap.Error(err))
	}

	if err = readData(db, l); err != nil {
		l.Fatal("reading data", zap.Error(err))
	}
}
