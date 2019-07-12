package resources

// URL is the struct around the Urls table
type URL struct {
	Resource
	URL  string `db:"url"`
	Slug string `db:"slug"`
}

// GetValues returns back a map of values about the User resource
func (url URL) GetValues() map[string]interface{} {
	return map[string]interface{}{
		"id":         url.ID,
		"url":        url.URL,
		"slug":       url.Slug,
		"created_at": url.CreatedAt,
		"deleted_at": url.DeletedAt,
	}
}

// Table returns the table associated with the resource
func (url URL) Table() string {
	return "urls"
}
