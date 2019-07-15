package endpoints

import (
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"time"

	"shorty/accessors"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// OriginURL is for binding to the URL request
type OriginURL struct {
	URL string `json:"url" binding:"required,url"`
}

// ShortenedURLResponse is for return to the URL request
type ShortenedURLResponse struct {
	ID           int64  `json:"id"`
	ShortenedURL string `json:"shortened_url"`
}

// ErrUnableToGenerateUniqueSlug used to explain that it was unable to generate a unique slug
var ErrUnableToGenerateUniqueSlug = errors.New("Unable to generate a unique slug")

const slugGeneratedLength = 20

// NewShortenedURL is for creating a new shortened URL
func (e *Endpointer) NewShortenedURL(c *gin.Context) {
	var json OriginURL

	err := c.ShouldBindWith(&json, binding.JSON)
	if err == nil {
		url := json.URL

		// Create URL
		u := accessors.URLDataAccessor{Databaser: e.databaser}
		randomSlug, nErr := NewSlug(u, 10)
		if nErr != nil {
			Error(nErr.Error(), c)
			return
		}

		createdURL, uErr := u.CreateURL(url, randomSlug)
		if uErr != nil {
			Error(uErr.Error(), c)
			return
		}

		// Create the response
		shortenedURL := constructShortURL(createdURL.Slug)
		response := ShortenedURLResponse{
			ID:           createdURL.ID,
			ShortenedURL: shortenedURL,
		}

		c.JSON(
			http.StatusCreated,
			response,
		)
	} else {
		HandleError(err, c)
	}
}

// NewSlug returns a unique slug that currently is not used in the DB
func NewSlug(u accessors.URLAccessor, maxRetries int) (string, error) {
	for i := 0; i < maxRetries; i++ {
		slug := RandomGeneratedSlug(slugGeneratedLength)

		exists, err := u.SlugExists(slug)
		if err != nil {
			return "", err
		}

		if !exists {
			return slug, nil
		}
	}

	return "", ErrUnableToGenerateUniqueSlug
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

// RandomGeneratedSlug returns a randomly made string of N characters
func RandomGeneratedSlug(n int) string {
	rand.Seed(time.Now().UnixNano())

	// This will be used to make testing easier
	if os.Getenv("STATIC_RANDOM_SLUG") != "" {
		return os.Getenv("STATIC_RANDOM_SLUG")
	}

	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func constructShortURL(slug string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("SHORTY_HOST"), slug)
}
