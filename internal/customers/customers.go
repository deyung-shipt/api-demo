package customers

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"strconv"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

// ID ...
type ID string

// Scan implements the Scanner interface.
func (i *ID) Scan(value interface{}) error {
	vI64, ok := value.(int64)
	if !ok {
		return errors.New("id.Scan: expected int64")
	}
	*i = ID(strconv.Itoa(int(vI64)))
	return nil
}

// Value implements the driver Valuer interface.
func (i ID) Value() (driver.Value, error) {
	iAsInt, err := strconv.Atoi(string(i))
	return int64(iAsInt), err
}

// Customer ...
type Customer struct {
	ID        ID     `db:"id"`
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

// CreateParams ...
type CreateParams struct {
	Email     string `db:"email"`
	FirstName string `db:"first_name"`
	LastName  string `db:"last_name"`
}

// CreateResult ...
type CreateResult struct {
	ID string
}

// ErrNotFound ...
var (
	ErrEmailAlreadyExists = errors.New("customer with provided email already exists")
	ErrNotFound           = errors.New("customer not found")
)

// DB ...
type DB struct {
	client *sqlx.DB
}

func (d DB) Close() error {
	return d.client.Close()
}

// DBConfig ...
type DBConfig struct {
	Client *sqlx.DB
}

// NewDB ...
func NewDB(cfg DBConfig) *DB {
	return &DB{
		client: cfg.Client,
	}
}

// CreateCustomer creates a new customer ...
func (d *DB) CreateCustomer(ctx context.Context, params CreateParams) (*CreateResult, error) {
	rows, err := d.client.NamedQueryContext(ctx, createCustomer, &params)
	if err != nil {
		return nil, d.checkConstraintErr(err)
	}
	defer rows.Close()

	var res CreateResult
	if !rows.Next() {
		return nil, sql.ErrNoRows
	}
	if err := rows.Scan(&res.ID); err != nil {
		return nil, err
	}
	return &res, nil
}

// GetCustomer returns a customer record by ID.  If no Customer is found, returned error
// is set to ErrNotFound ...
func (d *DB) GetCustomer(ctx context.Context, id string) (*Customer, error) {
	var c Customer
	err := d.client.
		QueryRowxContext(ctx, getCustomer, ID(id)).
		StructScan(&c)
	if err == sql.ErrNoRows {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return &c, nil
}

// HealthCheck ...
func (d *DB) HealthCheck(ctx context.Context) error {
	var i int
	return d.client.QueryRowContext(ctx, "SELECT 1").Scan(&i)
}

// checkConstraintErr checks for named constraint violations
// and bubbles up a well defined error for any recognized volations
func (d *DB) checkConstraintErr(err error) error {
	const (
		constraintCustomerEmail = "customers_email_idx"
	)
	pqErr, ok := err.(*pq.Error)
	if !ok {
		return err
	}
	switch pqErr.Constraint {
	case constraintCustomerEmail:
		return ErrEmailAlreadyExists
	}
	return err
}
