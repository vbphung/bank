package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgresql://root:postgres_password@localhost:5432/simplebank?sslmode=disable"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	var err error

	testDb, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("connect db failed:", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
