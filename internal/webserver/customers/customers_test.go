package customers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/shipt/tempest-template/internal/test"
	"github.com/shipt/tempest/httpt/httpres"
)

func Test_CustomerCRUD(t *testing.T) {
	// setup the test environment
	test.WithEnv(func(e *test.Env) {
		handler := NewHandler(HandlerConfig{
			CustomersBaseURL: "http://customer-service.shipt.test",
			DB:               e.CustomerDB,
		})

		var customerURL string

		t.Run("create a customer", func(t *testing.T) {
			var reqBody bytes.Buffer

			err := json.NewEncoder(&reqBody).Encode(&createCustomerRequest{
				Email:     "bobloblaw@thebobloblawlawblog.com",
				FirstName: "Bob",
				LastName:  "Loblaw",
			})
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/v1/customers/", &reqBody)
			assert.NoError(t, err)

			// make the request
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// check the response status
			assert.Equal(t, http.StatusCreated, rec.Code)

			// grab the location header from the response
			// we'll use this to GET the customer in a
			// follow up test ...
			customerURL = rec.Header().Get("Location")
		})

		t.Run("create a customer with same email", func(t *testing.T) {
			var reqBody bytes.Buffer

			err := json.NewEncoder(&reqBody).Encode(&createCustomerRequest{
				Email:     "bobloblaw@thebobloblawlawblog.com",
				FirstName: "Bob",
				LastName:  "Loblaw",
			})
			assert.NoError(t, err)

			req, err := http.NewRequest(http.MethodPost, "/v1/customers/", &reqBody)
			assert.NoError(t, err)

			// make the request
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// check the response status
			assert.Equal(t, http.StatusUnprocessableEntity, rec.Code)

			// check the error type
			var e struct {
				Error httpres.Error `json:"error"`
			}
			assert.NoError(t, json.NewDecoder(rec.Body).Decode(&e))
			assert.Equal(t, "customer_email_already_exists", e.Error.Type)
		})

		t.Run("get the customer", func(t *testing.T) {
			// follow the location header to get the newly created customer ...
			req, err := http.NewRequest(http.MethodGet, customerURL, http.NoBody)
			assert.NoError(t, err)

			// make the request
			rec := httptest.NewRecorder()
			handler.ServeHTTP(rec, req)

			// check the response status and body
			assert.Equal(t, http.StatusOK, rec.Code)

			var c customerResponse
			assert.NoError(t, json.NewDecoder(rec.Body).Decode(&c))

			assert.Equal(t, "bobloblaw@thebobloblawlawblog.com", c.Email)
			assert.Equal(t, "Bob", c.FirstName)
			assert.Equal(t, "Loblaw", c.LastName)
		})
	})
}
