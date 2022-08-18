package customers

import (
	"encoding/json"
	"net/http"

	"github.com/shipt/tempest/httpt"
	"github.com/shipt/tempest/httpt/httpres"

	"github.com/shipt/tempest-template/internal/customers"
)

var errCustomerEmailAlreadyExists = httpres.Error{
	Message: "a customer with the requested email address already exists",
	Type:    "customer_email_already_exists",
	Status:  http.StatusUnprocessableEntity,
}

var errCustomerNotFound = httpres.Error{
	Message: "the customer you requested was not found",
	Type:    "customer_not_found",
	Status:  http.StatusNotFound,
}

type handler struct {
	http.Handler

	customersBaseURL string
	db               *customers.DB
}

// HandlerConfig ...
type HandlerConfig struct {
	CustomersBaseURL string
	DB               *customers.DB
}

type createCustomerRequest struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func (c *createCustomerRequest) validate() error {
	// TODO: implement validation ...
	return nil
}

type customerResponse struct {
	ID        string `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// NewHandler ...
func NewHandler(cfg HandlerConfig) http.Handler {
	h := handler{
		customersBaseURL: cfg.CustomersBaseURL,
		db:               cfg.DB,
	}

	// setup handler routing ...
	mux := httpt.NewMux()
	mux.HandleFunc("/v1/customers/", http.MethodPost, h.createCustomer)
	mux.HandleFunc("/v1/customers/{id}", http.MethodGet, h.getCustomer)

	h.Handler = mux
	return &h
}

func (h *handler) createCustomer(w http.ResponseWriter, r *http.Request) {
	var reqBody createCustomerRequest

	// read the request body, returning a 400 - Bad Request if
	// the body is empty or isn't valid JSON
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		httpres.WriteError(w, r, httpres.Error{
			Type:    "invalid_json_body",
			Message: err.Error(),
			Status:  http.StatusBadRequest,
		})
		return
	}

	// validate the request body ...
	if err := reqBody.validate(); err != nil {
		httpres.WriteError(w, r, err)
		return
	}

	// create the customer in the DB ...
	cRes, err := h.db.CreateCustomer(r.Context(), customers.CreateParams{
		Email:     reqBody.Email,
		FirstName: reqBody.FirstName,
		LastName:  reqBody.LastName,
	})
	if err == customers.ErrEmailAlreadyExists {
		httpres.WriteError(w, r, errCustomerEmailAlreadyExists)
		return
	}
	if err != nil {
		httpres.WriteError(w, r, err)
		return
	}

	// redirect to the newly created customer resource ..
	httpres.Write(w, r,
		httpres.WithHeader("Location", h.customerURL(cRes.ID)),
		httpres.WithStatus(http.StatusCreated),
	)
}

func (h *handler) getCustomer(w http.ResponseWriter, r *http.Request) {
	// fetch the customer from the database
	c, err := h.db.GetCustomer(r.Context(), httpt.PathVars(r)["id"])

	if err == customers.ErrNotFound {
		httpres.WriteError(w, r, errCustomerNotFound)
		return
	}
	if err != nil {
		httpres.WriteError(w, r, err)
		return
	}
	resObj := customerResponse{
		ID:        string(c.ID),
		Email:     c.Email,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}
	httpres.Write(w, r,
		httpres.WithStatus(http.StatusOK),
		httpres.WithJSON(resObj),
	)
}

func (h *handler) customerURL(id string) string {
	return h.customersBaseURL + "/v1/customers/" + id
}
