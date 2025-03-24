package sqlite

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"

	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/database/schema"
)

func MustInitDB() *sql.DB {
	db, err := sql.Open("sqlite", "file::memory:?_pragma=foreign_keys(on)")
	if err != nil {
		panic(err)
	}
	// create tables
	if _, err := db.ExecContext(context.TODO(), schema.DDL_Schema); err != nil {
		panic(err)
	}
	// initFunction for "default" handled in repo.
	database.New(db).InsertRegularChannel(context.TODO(), "default")

	return db
}
