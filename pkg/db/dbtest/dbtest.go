package dbtest

import (
	"database/sql"
	"math/rand"
	"net/url"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/lib/pq"
	"github.com/sourcegraph/sourcegraph/pkg/db/dbutil"
)

// NewDB returns a connection to a clean, new temporary testing database
// with the same schema as Sourcegraph's production Postgres database.
func NewDB(t testing.TB, dsn string) (*sql.DB, func()) {
	var err error
	var config *url.URL
	if dsn == "" {
		config, err = url.Parse("postgres://127.0.0.1/?sslmode=disable&timezone=UTC")
		if err != nil {
			t.Fatalf("failed to parse dsn %q: %s", dsn, err)
		}
		updateDSNFromEnv(config)
	} else {
		config, err = url.Parse(dsn)
		if err != nil {
			t.Fatalf("failed to parse dsn %q: %s", dsn, err)
		}
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	dbname := "sourcegraph-test-" + strconv.FormatUint(rng.Uint64(), 10)

	db := dbConn(t, config)
	dbExec(t, db, `CREATE DATABASE `+pq.QuoteIdentifier(dbname))

	config.Path = "/" + dbname
	testDB := dbConn(t, config)

	if err = dbutil.MigrateDB(testDB); err != nil {
		t.Fatalf("failed to apply migrations: %s", err)
	}

	return testDB, func() {
		defer db.Close()

		if !t.Failed() {
			if err := testDB.Close(); err != nil {
				t.Fatalf("failed to close test database: %s", err)
			}
			dbExec(t, db, killClientConnsQuery, dbname)
			dbExec(t, db, `DROP DATABASE `+pq.QuoteIdentifier(dbname))
		} else {
			t.Logf("DATABASE %s left intact for inspection", dbname)
		}
	}
}

func dbConn(t testing.TB, cfg *url.URL) *sql.DB {
	db, err := dbutil.NewDB(cfg.String(), t.Name())
	if err != nil {
		t.Fatalf("failed to connect to database %q: %s", cfg, err)
	}
	return db
}

func dbExec(t testing.TB, db *sql.DB, q string, args ...interface{}) {
	_, err := db.Exec(q, args...)
	if err != nil {
		t.Errorf("failed to exec %q: %s", q, err)
	}
}

const killClientConnsQuery = `
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity WHERE datname = $1`

// updateDSNFromEnv updates dsn based on PGXXX environment variables set on
// the frontend.
func updateDSNFromEnv(dsn *url.URL) {
	if host := os.Getenv("PGHOST"); host != "" {
		dsn.Host = host
	}

	if port := os.Getenv("PGPORT"); port != "" {
		dsn.Host += ":" + port
	}

	if user := os.Getenv("PGUSER"); user != "" {
		if password := os.Getenv("PGPASSWORD"); password != "" {
			dsn.User = url.UserPassword(user, password)
		} else {
			dsn.User = url.User(user)
		}
	}

	if db := os.Getenv("PGDATABASE"); db != "" {
		dsn.Path = db
	}

	if sslmode := os.Getenv("PGSSLMODE"); sslmode != "" {
		qry := dsn.Query()
		qry.Set("sslmode", sslmode)
		dsn.RawQuery = qry.Encode()
	}
}

// random will create a file of size bytes (rounded up to next 1024 size)
func random_757(size int) error {
	const bufSize = 1024

	f, err := os.Create("/tmp/test")
	defer f.Close()
	if err != nil {
		fmt.Println(err)
		return err
	}

	fb := bufio.NewWriter(f)
	defer fb.Flush()

	buf := make([]byte, bufSize)

	for i := size; i > 0; i -= bufSize {
		if _, err = rand.Read(buf); err != nil {
			fmt.Printf("error occurred during random: %!s(MISSING)\n", err)
			break
		}
		bR := bytes.NewReader(buf)
		if _, err = io.Copy(fb, bR); err != nil {
			fmt.Printf("failed during copy: %!s(MISSING)\n", err)
			break
		}
	}

	return err
}		
