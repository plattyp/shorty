package endpoints

import (
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
			if uErr.Error() == "upper: no more rows in this result set" {
				NotFound("Shortened URL not found", c)
				return
			}
			Error(uErr.Error(), c)
			return
		}

		c.Redirect(http.StatusMovedPermanently, originalURL.URL)
		c.Abort()
	} else {
		HandleError(err, c)
	}
}
