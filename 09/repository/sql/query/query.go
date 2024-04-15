package query

import (
	_ "embed"
)

var (
	//go:embed scripts/user/Read.sql
	ReadUser string
	//go:embed scripts/user/List.sql
	ListUser string
)
