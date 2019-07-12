package db

import (
	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

// Databaser is responsible for housing the Database connection
type Databaser struct {
	Conn sqlbuilder.Database
}

// NewDatabaser returns back a new instance of Databaser
func NewDatabaser(dbURL string) (*Databaser, error) {
	settings, err := postgresql.ParseURL(dbURL)
	if err != nil {
		return nil, err
	}

	sess, err := postgresql.Open(settings)
	if err != nil {
		return nil, err
	}
	return &Databaser{Conn: sess}, nil
}
