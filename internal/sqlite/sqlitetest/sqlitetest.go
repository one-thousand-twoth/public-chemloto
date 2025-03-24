package sqlitetest

import (
	"database/sql"

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
	}

	return db, fn

}
