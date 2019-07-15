package db

import (
	"os"
	"strconv"

	"upper.io/db.v3/lib/sqlbuilder"
	"upper.io/db.v3/postgresql"
)

const dbAppName = "shorty"

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

	// Set SSLMode for deployed environments
	if isLiveEnvironment() {
		settings.Options["sslmode"] = "require"
	}

	// Set Connection Name To Track in Postgres
	settings.Options["application_name"] = dbAppName

	sess, err := postgresql.Open(settings)
	if err != nil {
		return nil, err
	}
	sess.SetMaxIdleConns(getMaxIdleConnections())
	sess.SetMaxOpenConns(getMaxOpenConnections())
	sess.SetLogging(enableDatabaseLogging())
	return &Databaser{Conn: sess}, nil
}

// Close returns the connection back to the pool
func (d Databaser) Close() error {
	return d.Conn.Close()
}

func getMaxIdleConnections() int {
	if maxIdleConns, exists := os.LookupEnv("MAX_IDLE_CONNECTIONS"); exists {
		i64, err := strconv.ParseInt(maxIdleConns, 10, 32)
		if err == nil {
			return int(i64)
		}
	}
	return 10
}

func getMaxOpenConnections() int {
	if maxOpenConns, exists := os.LookupEnv("MAX_OPEN_CONNECTIONS"); exists {
		i64, err := strconv.ParseInt(maxOpenConns, 10, 32)
		if err == nil {
			return int(i64)
		}
	}
	return 25
}

func enableDatabaseLogging() bool {
	if databaseLogging, exists := os.LookupEnv("ENABLE_DATABASE_LOGGING"); exists {
		dbLogging, err := strconv.ParseBool(databaseLogging)
		if err == nil {
			return dbLogging
		}
	}
	return false
}

func isLiveEnvironment() bool {
	return os.Getenv("SHORT_ENV") == "production"
}
