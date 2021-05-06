package sqltracing_test

import (
	"context"
	"database/sql"
	"testing"

	"github.com/siddhant2408/golang-libraries/mysqltest"
	"github.com/siddhant2408/golang-libraries/mysqltracing"
	"github.com/siddhant2408/golang-libraries/sqltracing"
	"github.com/siddhant2408/golang-libraries/testutils"
)

func TestMySQL(t *testing.T) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	db := mysqltest.GetDatabase(t)
	testMySQLTable(ctx, t, db)
	testMySQLInsert(ctx, t, db)
	testMySQLRow(ctx, t, db)
	testMySQLRows(ctx, t, db)
}

func testMySQLTable(ctx context.Context, t *testing.T, db *sql.DB) {
	_, err := db.ExecContext(ctx, `CREATE TABLE test (test TEXT);`)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func testMySQLInsert(ctx context.Context, t *testing.T, db *sql.DB) {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	stmt, err := tx.PrepareContext(ctx, `INSERT INTO test (test) VALUES (?)`)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	s1 := "test"
	_, err = stmt.ExecContext(ctx, s1)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = stmt.Close()
	if err != nil {
		testutils.FatalErr(t, err)
	}
	err = tx.Commit()
	if err != nil {
		testutils.FatalErr(t, err)
	}
}

func testMySQLRow(ctx context.Context, t *testing.T, db *sql.DB) {
	row := db.QueryRowContext(ctx, `SELECT * FROM test;`)
	var s string
	err := sqltracing.RowScan(ctx, row, &s)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	expected := "test"
	if s != expected {
		t.Fatalf("unexpected result: got %q, want %q", s, expected)
	}
}

func testMySQLRows(ctx context.Context, t *testing.T, db *sql.DB) {
	rows, err := db.QueryContext(ctx, `SELECT * FROM test;`)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	var called testutils.CallCounter
	err = sqltracing.RowsIterate(ctx, rows, func(ctx context.Context, rows *sql.Rows) error {
		called.Call()
		return nil
	})
	if err != nil {
		testutils.FatalErr(t, err)
	}
	called.AssertCount(t, 1)
}

func TestMySQLOpenDSN(t *testing.T) {
	ctx := context.Background()
	mysqltest.CheckAvailable(t)
	dsn := mysqltest.GetDSN()
	db, err := sql.Open(mysqltracing.DriverName, dsn)
	if err != nil {
		testutils.FatalErr(t, err)
	}
	defer db.Close() //nolint:errcheck
	err = db.PingContext(ctx)
	if err != nil {
		testutils.FatalErr(t, err)
	}
}
