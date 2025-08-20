package db

import (
	"context"
	"database/sql"

	_ "embed"

	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

//go:embed ../sql/schema.sql
var ddl string

func Connect(ctx context.Context) (*generated.Queries, error) {
	db, err := sql.Open("sqlite3", "file:local/local.db")
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	return generated.New(db), nil
}
