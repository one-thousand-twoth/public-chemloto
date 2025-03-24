package sqlite

import (
	"testing"

	_ "modernc.org/sqlite" // SQLite driver
)

func BenchmarkMustInitDB(b *testing.B) {
	for i := 0; i < b.N; i++ {
		db := MustInitDB() // Initialize DB
		db.Close()         // Ensure proper cleanup
	}
}
