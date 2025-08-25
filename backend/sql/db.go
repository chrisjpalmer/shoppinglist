package sql

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	_ "embed"

	"github.com/chrisjpalmer/shoppinglist/backend/generated"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

//go:embed schema.sql
var ddl string

func Connect(ctx context.Context) (*generated.Queries, error) {
	dbPath := "local/local.db"

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbPath))
	if err != nil {
		return nil, err
	}

	exists, err := dbExists(dbPath)
	if err != nil {
		return nil, err
	}

	// create the tables for the first time if they don't exist
	if !exists {
		if _, err := db.ExecContext(ctx, ddl); err != nil {
			return nil, err
		}
	}

	return generated.New(db), nil
}

func dbExists(dbPath string) (bool, error) {
	_, err := os.Stat(dbPath)
	if err != nil && !os.IsNotExist(err) {
		return false, err
	}

	return !os.IsNotExist(err), nil
}
