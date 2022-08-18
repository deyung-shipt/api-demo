package customers

import (
	"regexp"
	"strings"
)

// space eater and compact strip extraneous white space
// from queries so they are more readable in logs and traces
var spaceEater = regexp.MustCompile(`\s+`)

func compact(q string) string {
	q = strings.TrimSpace(q)
	return spaceEater.ReplaceAllString(q, " ")
}

var createCustomer = compact(`
	insert into customers (email, first_name, last_name, created_at, updated_at) values 
		( 
			:email, 
			:first_name, 
			:last_name, 
			CURRENT_TIMESTAMP, 
			CURRENT_TIMESTAMP
		)
	returning id
`)

var getCustomer = compact(`
	select
		id,
		email,
		first_name,
		last_name
	from
		customers
	where
		id = $1
`)
