package query

import (
	_ "embed"
)

var (
	//go:embed scripts/CreateCustomer.sql
	CreateCustomer string
	//go:embed scripts/CreateCustomerAddress.sql
	CreateCustomerAddress string
	//go:embed scripts/CreateProduct.sql
	CreateProduct string
	//go:embed scripts/AddProductToCart.sql
	AddProductToCart string

	//go:embed scripts/ListCustomers.sql
	ListCostumers string
	//go:embed scripts/ListCustomerAddresses.sql
	ListCustomerAddresses string
	//go:embed scripts/ListCartItems.sql
	ListCartItems string
)
