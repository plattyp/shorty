package accessors

import (
	"errors"

	"shorty/resources"

	"shorty/db"

	upper "upper.io/db.v3"
)

// URLAccessor is used to fetch/create urls
type URLAccessor interface {
	GetURL(id int64) (*resources.URL, error)
	GetURLBySlug(slug string) (*resources.URL, error)
	CreateURL(url string, slug string) (*resources.URL, error)
}

const urlTableName = "urls"

// URLDataAccessor is used to interacts with urls
type URLDataAccessor struct {
	Databaser *db.Databaser
}

// CreateURL will create a url based on the params
func (u URLDataAccessor) CreateURL(url string, slug string) (*resources.URL, error) {
	createdURL := resources.URL{
		URL:  url,
		Slug: slug,
	}

	result, err := u.urlsTable().Insert(&createdURL)
	if err != nil {
		return nil, err
	}

	if id, ok := result.(int64); ok {
		user, err := u.GetURL(id)
		if err != nil {
			return nil, err
		}
		return user, nil
	}

	return nil, errors.New("Failed to create url")
}

// GetURL retrieves a URL by its ID
func (u URLDataAccessor) GetURL(id int64) (*resources.URL, error) {
	var url resources.URL

	err := u.urlsTable().Find(upper.Cond{"id": id}).One(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

// GetURLBySlug retrieves a URL by its ID
func (u URLDataAccessor) GetURLBySlug(slug string) (*resources.URL, error) {
	var url resources.URL

	err := u.urlsTable().Find(upper.Cond{"slug": slug}).One(&url)
	if err != nil {
		return nil, err
	}

	return &url, nil
}

// usersTable returns back a collection
func (u URLDataAccessor) urlsTable() upper.Collection {
	return u.Databaser.Conn.Collection(urlTableName)
}
