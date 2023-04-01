package repository

import (
	"context"
	"time"

	"data-persistence/generics/model"
	"data-persistence/generics/repository/query"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type Repository struct {
	database database
}

func New(ctx context.Context, dsn string) (Repository, error) {
	db, err := newDatabase(ctx, dsn)
	if err != nil {
		return Repository{}, err
	}
	return Repository{database: db}, nil
}

func (r Repository) CreateCustomerWithAddress(ctx context.Context, customer *model.Customer, address *model.CustomerAddress) error {
	return withTransaction(ctx, r.database.pool, func(q querier) error {
		err := execOne(ctx, q, query.CreateCustomer, pgx.NamedArgs{
			"id":         customer.ID,
			"name":       customer.Name,
			"email":      customer.Email,
			"created_at": customer.CreatedAt,
			"updated_at": customer.UpdatedAt,
		})
		if err != nil {
			return err
		}

		err = execOne(ctx, q, query.CreateCustomerAddress, pgx.NamedArgs{
			"id":            address.ID,
			"customer_id":   address.CustomerID,
			"location_name": address.LocationName,
			"address":       address.Address,
			"created_at":    address.CreatedAt,
			"updated_at":    address.UpdatedAt,
		})
		if err != nil {
			return err
		}
		return nil
	})
}

func (r Repository) CreateProducts(ctx context.Context, products []model.Product) error {
	return withConnection(ctx, r.database.pool, func(q querier) error {
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
		result := q.SendBatch(ctx, batch)
		for i := 0; i < len(products); i++ {
			if _, err := result.Exec(); err != nil {
				return err
			}
		}
		if err := result.Close(); err != nil {
			return err
		}
		return nil
	})
}

func (r Repository) AddProductsToCart(ctx context.Context, customerID uuid.UUID, productIDs []uuid.UUID) error {
	insertedAt := time.Now()
	return withConnection(ctx, r.database.pool, func(q querier) error {
		batch := &pgx.Batch{}
		for _, id := range productIDs {
			batch.Queue(query.AddProductToCart, pgx.NamedArgs{
				"customer_id": customerID,
				"product_id":  id,
				"inserted_at": insertedAt,
			})
		}
		result := q.SendBatch(ctx, batch)
		for i := 0; i < len(productIDs); i++ {
			if _, err := result.Exec(); err != nil {
				return err
			}
		}
		if err := result.Close(); err != nil {
			return err
		}
		return nil
	})
}

func (r Repository) ListCustomers(ctx context.Context) ([]model.Customer, error) {
	return withConnectionResult(ctx, r.database.pool, func(q querier) ([]model.Customer, error) {
		customers, err := list[customer](ctx, q, query.ListCostumers)
		if err != nil {
			return nil, err
		}
		return toCustomers(customers), nil
	})
}

func (r Repository) ListCustomerAddresses(ctx context.Context, customerID uuid.UUID) ([]model.CustomerAddress, error) {
	return withConnectionResult(ctx, r.database.pool, func(q querier) ([]model.CustomerAddress, error) {
		customerAddresses, err := list[customerAddress](ctx, q, query.ListCustomerAddresses, pgx.NamedArgs{
			"customer_id": customerID,
		})
		if err != nil {
			return nil, err
		}
		return toCustomerAddresses(customerAddresses), nil
	})
}

func (r Repository) ListCartItems(ctx context.Context, customerID uuid.UUID) ([]model.CartItem, error) {
	return withConnectionResult(ctx, r.database.pool, func(q querier) ([]model.CartItem, error) {
		cartItems, err := list[cartItem](ctx, q, query.ListCartItems, pgx.NamedArgs{
			"customer_id": customerID,
		})
		if err != nil {
			return nil, err
		}
		return toCartItems(cartItems), err
	})
}
