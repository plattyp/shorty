package resources

import (
	"time"

	"github.com/lib/pq"
)

// Resourcer is a generic resource interface
type Resourcer interface {
	GetValues() map[string]interface{}
	TableName() string
}

// Resource contains all basic functionality that all database models share
type Resource struct {
	Resourcer
	ID        int64       `db:"id,omitempty"`
	CreatedAt time.Time   `db:"created_at,omitempty"`
	DeletedAt pq.NullTime `db:"deleted_at,omitempty"`
}

// GetValues returns back a map of interfaces of the value of the resource
func (resource Resource) GetValues() map[string]interface{} {
	return map[string]interface{}{
		"id":         resource.ID,
		"created_at": resource.CreatedAt,
		"deleted_at": resource.DeletedAt,
	}
}

// TableName returns back the name of the table within the database
func (resource Resource) TableName() string {
	return ""
}
