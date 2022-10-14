package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
	"github.com/vbph/bank/utils"
)

var testQueries *Queries
var testDb *sql.DB

func TestMain(m *testing.M) {
	cfg, err := utils.LoadConfig("../../")
	if err != nil {
		log.Fatalf("load config failed: %+v\n", err)
	}

	testDb, err = sql.Open(cfg.DbDriver, cfg.DbSource)
	if err != nil {
		log.Fatalf("connect db failed: %+v\n", err)
	}

	testQueries = New(testDb)

	os.Exit(m.Run())
}
