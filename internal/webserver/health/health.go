package health

import (
	"net/http"

	"github.com/shipt/tempest/httpt"
	"github.com/shipt/tempest/httpt/httpres"

	"github.com/shipt/tempest-template/internal/customers"
)

type handler struct {
	http.Handler

	customersDB *customers.DB
}

// HandlerConfig ...
type HandlerConfig struct {
	CustomersDB *customers.DB
	// add any configuration or dependencies here ...
}

func notReadyErr(err error) error {
	return httpres.Error{
		Message: err.Error(),
		Type:    "service_unavailable",
		Status:  http.StatusServiceUnavailable,
	}
}

// NewHandler ...
func NewHandler(cfg HandlerConfig) http.Handler {
	h := handler{
		customersDB: cfg.CustomersDB,
	}

	mux := httpt.NewMux()
	mux.HandleFunc("/health/live", http.MethodGet, h.getLiveness)
	mux.HandleFunc("/health/ready", http.MethodGet, h.getReadiness)

	h.Handler = mux
	return h
}

func (h *handler) getLiveness(w http.ResponseWriter, r *http.Request) {
	// add more comprehensive liveness checks here
	// https://github.com/shipt/infraspec-proto/blob/master/INFRASPEC-DOCS.md#health
	httpres.Write(w, r)
}

func (h *handler) getReadiness(w http.ResponseWriter, r *http.Request) {
	if err := h.customersDB.HealthCheck(r.Context()); err != nil {
		httpres.WriteError(w, r, notReadyErr(err))
	}
	// add more comprehensive readiness checks here
	// https://github.com/shipt/infraspec-proto/blob/master/INFRASPEC-DOCS.md#health
	httpres.Write(w, r)
}
