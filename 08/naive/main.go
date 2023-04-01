package main

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"data-persistence/naive/query"

	"github.com/google/uuid"
	_ "github.com/jackc/pgx/v5/stdlib"
	"go.uber.org/zap"
)

const (
	dsn = "postgres://root:root@localhost:5432/data-persistence?sslmode=disable"
)

func createCustomer(tx *sql.Tx) (*Customer, error) {
	customer := NewCustomer("John Doe", "john.doe@gmail.com")
	_, err := tx.Exec(
		query.CreateCustomer,
		customer.ID,
		customer.Name,
		customer.Email,
		customer.CreatedAt,
		customer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	customerAddress := NewCustomerAddress(
		customer.ID,
		"Work",
		"Rohanské nábř. 678/23, 186 00 Karlín",
	)
	_, err = tx.Exec(
		query.CreateCustomerAddress,
		customerAddress.ID,
		customerAddress.CustomerID,
		customerAddress.LocationName,
		customerAddress.Address,
		customerAddress.CreatedAt,
		customer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return customer, nil
}

func createProducts(tx *sql.Tx) ([]Product, error) {
	products := []Product{
		*NewProduct("Lagkapten", "Desk, dark grey/black, 140x60 cm", 1499),
		*NewProduct("Strandmon", "Wing chair, Nordvalla dark grey", 5490),
	}
	for _, v := range products {
		_, err := tx.Exec(
			query.CreateProduct,
			v.ID,
			v.Name,
			v.Description,
			v.Price,
			v.CreatedAt,
			v.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
	}

	return products, nil
}

func addProductsToCart(tx *sql.Tx, customer *Customer, products []Product) error {
	insertedAt := time.Now()
	for _, v := range products {
		if _, err := tx.Exec(query.AddProductToCart, customer.ID, v.ID, insertedAt); err != nil {
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
	defer func() {
		if err != nil {
			if rErr := tx.Rollback(); rErr != nil {
				err = errors.Join(rErr, err)
			}
			return
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

func listCustomers(db *sql.DB) ([]Customer, error) {
	rows, err := db.Query(query.ListCostumers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customers []Customer
	for rows.Next() {
		c := Customer{}
		err = rows.Scan(
			&c.ID,
			&c.Name,
			&c.Email,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customers = append(customers, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customers, nil
}

func listCustomerAddresses(db *sql.DB, customerID uuid.UUID) ([]CustomerAddress, error) {
	rows, err := db.Query(query.ListCustomerAddresses, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var customerAddresses []CustomerAddress
	for rows.Next() {
		c := CustomerAddress{}
		err = rows.Scan(
			&c.ID,
			&c.CustomerID,
			&c.LocationName,
			&c.Address,
			&c.CreatedAt,
			&c.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		customerAddresses = append(customerAddresses, c)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return customerAddresses, nil
}

func listCartItems(db *sql.DB, customerID uuid.UUID) ([]CartItem, error) {
	rows, err := db.Query(query.ListCartItems, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cartItems []CartItem
	for rows.Next() {
		p := CartItem{}
		err = rows.Scan(
			&p.ID,
			&p.Name,
			&p.Description,
			&p.Price,
			&p.InsertedAt,
			&p.CreatedAt,
			&p.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		cartItems = append(cartItems, p)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return cartItems, nil
}

func readData(db *sql.DB, l *zap.Logger) error {
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
		cartItems, err := listCartItems(db, c.ID)
		if err != nil {
			return err
		}
		l.With(
			zap.String("customer_id", c.ID.String()),
			zap.Any("cart_items", cartItems),
		).Info("listing items in cart")
	}

	return nil
}

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db, err := sql.Open("pgx", dsn)
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
