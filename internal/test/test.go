package test

import (
	"net/url"
	"os"

	"github.com/jmoiron/sqlx"
	"github.com/shipt/tempest-template/internal/customers"
	"github.com/shipt/tempest/database/sqlt"
)

// Env ...
type Env struct {
	CustomerDB *customers.DB
}

// WithEnv ...
func WithEnv(f func(*Env)) {
	sqlClient := sqlt.NewPostgresSQLxDB(sqlt.PostgresConfig{
		URL: safeGetPSQLURL(),
	})

	// clean the db prior to running tests
	mustCleanTable(sqlClient, "customers")

	f(&Env{
		CustomerDB: customers.NewDB(customers.DBConfig{
			Client: sqlClient,
		}),
	})
}

func mustCleanTable(client *sqlx.DB, table string) {
	_, err := client.Exec("TRUNCATE TABLE " + table)
	if err != nil {
		panic(err)
	}
}

func safeGetPSQLURL() string {
	allowedHostnames := map[string]bool{
		"localhost": true,
		"database":  true,
		"postgres":  true,
	}
	psqlURL := "postgres://pguser:pgpass@localhost:5432/tempest-template?sslmode=disable"
	if p := os.Getenv("PSQL_URL"); p != "" {
		psqlURL = p
	}
	parsed, err := url.Parse(psqlURL)
	if err != nil {
		panic("invalid psql url: " + psqlURL)
	}
	if !allowedHostnames[parsed.Hostname()] {
		panic("Not so fast!  are you sure that " + psqlURL + " is a test database?")
	}
	return psqlURL
}
