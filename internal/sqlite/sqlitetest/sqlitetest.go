package sqlitetest

import (
	"context"
	"database/sql"

	"github.com/anrew1002/Tournament-ChemLoto/internal/database"
	"github.com/anrew1002/Tournament-ChemLoto/internal/sqlite"
)

func GetTestDatabase() (*sql.DB, func()) {

	db := sqlite.MustInitDB()

	fn := func() {
		db.Exec(`
		DELETE FROM users;
		DELETE FROM rooms;
		DELETE FROM channels;
		DELETE FROM channel_subscribers;
		`)
		if _, err := database.New(db).InsertRegularChannel(context.TODO(), "default"); err != nil {
			panic("ERROR WHILE INSERT DEFAULT CHANNEL")
		}
	}

	return db, fn

}
