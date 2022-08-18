package webserver

import (
	"context"
	"github.com/shipt/bubinga"
	"net/http"

	"github.com/shipt/tempest/database/sqlt"
	"github.com/shipt/tempest/httpt"

	"github.com/shipt/tempest-template/doc"
	"github.com/shipt/tempest-template/internal/customers"
	wscustomers "github.com/shipt/tempest-template/internal/webserver/customers"
	wshealth "github.com/shipt/tempest-template/internal/webserver/health"
)

// Config ...
type Config struct {
	CustomersBaseURL string              `config:"CUSTOMERS_BASE_URL"`
	PostgresConfig   sqlt.PostgresConfig `config:",squash"`
}

// Run ...
func Run(ctx context.Context, cfg Config) error {
	customersDB := customers.NewDB(customers.DBConfig{
		Client: sqlt.NewPostgresSQLxDB(cfg.PostgresConfig),
	})

	defer func() {
		if err := customersDB.Close(); err != nil {
			bubinga.Error(ctx, "unable to close db", err)
		}

	}()

	mux := http.NewServeMux()
	mux.Handle("/v1/customers/", wscustomers.NewHandler(wscustomers.HandlerConfig{
		CustomersBaseURL: cfg.CustomersBaseURL,
		DB:               customersDB,
	}))

	// static content (e.g. documentation) handler
	mux.Handle("/doc/", doc.Handler)

	//run the server with a custom handler for the health routes, terminating either on error or when ctx is cancelled
	srv := httpt.NewServer(
		httpt.ServerConfig{Handler: mux},
		httpt.WithHealthHandler(
			"/health/",
			wshealth.NewHandler(wshealth.HandlerConfig{CustomersDB: customersDB}),
		),
	)
	return srv.Run(ctx)
}
