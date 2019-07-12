package endpoints

import (
	"fmt"
	"net/http"
	"shorty/accessors"

	"github.com/gin-gonic/gin"
)

// RequestURL is for binding to the redirect endpoint
type RequestURL struct {
	Slug string `uri:"slug" binding:"required"`
}

// RedirectURL is for redirecting to URL
func (e *Endpointer) RedirectURL(c *gin.Context) {
	var r RequestURL

	err := c.ShouldBindUri(&r)
	if err == nil {
		slug := r.Slug

		// Create URL
		u := accessors.URLDataAccessor{Databaser: e.databaser}
		originalURL, uErr := u.GetURLBySlug(slug)
		if uErr != nil {
			Error(uErr.Error(), c)
			return
		}

		fmt.Println(originalURL)

		c.Redirect(http.StatusMovedPermanently, originalURL.URL)
		c.Abort()
	} else {
		validationErr := HandleError(err)
		ValidationError(validationErr, c)
	}
}
