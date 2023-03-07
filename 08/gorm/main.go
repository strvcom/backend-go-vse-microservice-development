package main

import (
	"time"

	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	dsn = "postgres://root:root@localhost:5432/data-persistence?sslmode=disable"
)

func createCustomer(db *gorm.DB) (*Customer, error) {
	customer := NewCustomer("John Doe", "john.doe@gmail.com")
	customerAddress := NewCustomerAddress(
		customer.ID,
		"Work",
		"Rohanské nábř. 678/23, 186 00 Karlín",
	)

	if err := db.Create(customer).Error; err != nil {
		return nil, err
	}
	if err := db.Create(customerAddress).Error; err != nil {
		return nil, err
	}

	return customer, nil
}

func createProducts(db *gorm.DB) ([]Product, error) {
	products := []Product{
		*NewProduct("Lagkapten", "Desk, dark grey/black, 140x60 cm", 1499),
		*NewProduct("Strandmon", "Wing chair, Nordvalla dark grey", 5490),
	}

	if err := db.Create(products).Error; err != nil {
		return nil, err
	}

	return products, nil
}

func addProductsToCart(db *gorm.DB, customer *Customer, products ...Product) error {
	insertedAt := time.Now()
	var cartProducts []CartItem
	for _, p := range products {
		cartProducts = append(cartProducts, CartItem{
			CustomerID: customer.ID,
			Product:    Product{ID: p.ID},
			InsertedAt: insertedAt,
		})
	}

	if err := db.Create(cartProducts).Error; err != nil {
		return err
	}

	return nil
}

func createData(db *gorm.DB) error {
	err := db.Transaction(func(tx *gorm.DB) error {
		customer, err := createCustomer(tx)
		if err != nil {
			return err
		}
		products, err := createProducts(tx)
		if err != nil {
			return err
		}
		err = addProductsToCart(tx, customer, products...)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil
}

func readData(db *gorm.DB, l *zap.Logger) error {
	var customers []Customer
	if err := db.Find(&customers).Error; err != nil {
		return err
	}
	l.With(zap.Any("customers", customers)).Info("listing customers")

	for _, c := range customers {
		var customerAddresses []CustomerAddress
		if err := db.Find(&customerAddresses, "customer_id = ?", c.ID).Error; err != nil {
			return err
		}
		l.With(zap.Any("customer_addresses", customerAddresses)).Info("listing customer addresses")
	}

	for _, c := range customers {
		var cartItems []CartItem
		if err := db.Preload("Product").Find(&cartItems, "customer_id = ?", c.ID).Error; err != nil {
			return err
		}
		l.With(
			zap.String("customer_id", c.ID.String()),
			zap.Any("cart_items", cartItems),
		).Info("listing products in cart")
	}

	return nil
}

func main() {
	l, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})
	if err != nil {
		l.Fatal("opening db connection", zap.Error(err))
	}
	err = db.AutoMigrate(
		&Customer{},
		&CustomerAddress{},
		&Product{},
		&CartItem{},
	)
	if err != nil {
		l.Fatal("auto-migrating database schema", zap.Error(err))
	}

	if err = createData(db); err != nil {
		l.Fatal("creating data", zap.Error(err))
	}

	if err = readData(db, l); err != nil {
		l.Fatal("reading data", zap.Error(err))
	}
}
