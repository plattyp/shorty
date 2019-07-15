package endpoints

import (
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
	URL string `json:"url" binding:"required"`
}

// ShortenedURLResponse is for return to the URL request
type ShortenedURLResponse struct {
	ID           int64  `json:"id"`
	ShortenedURL string `json:"shortened_url"`
}

const slugGeneratedLength = 20

// NewShortenedURL is for creating a new shortened URL
func (e *Endpointer) NewShortenedURL(c *gin.Context) {
	var json OriginURL

	err := c.ShouldBindWith(&json, binding.JSON)
	if err == nil {
		url := json.URL

		// Create URL
		u := accessors.URLDataAccessor{Databaser: e.databaser}
		rand.Seed(time.Now().UnixNano())
		randomSlug := randomGeneratedSlug(slugGeneratedLength)
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
		validationErr := HandleError(err)
		ValidationError(validationErr, c)
	}
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ123456789"

func randomGeneratedSlug(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func constructShortURL(slug string) string {
	return fmt.Sprintf("%s/%s", os.Getenv("SHORTY_HOST"), slug)
}
