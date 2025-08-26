package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

func main() {
	dbPath := "./local/local.db"

	db, err := sql.Open("sqlite3", fmt.Sprintf("file:%s", dbPath))
	if err != nil {
		log.Fatal(err)
	}

	_, err = db.Exec("ALTER TABLE meals ADD COLUMN recipe_url TEXT NOT NULL DEFAULT ''")
	if err != nil {
		log.Fatal(err)
	}
}
