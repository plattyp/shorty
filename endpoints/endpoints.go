package endpoints

import (
	"fmt"
	"net/http"

	"shorty/db"

	"github.com/gin-gonic/gin"
	validator "gopkg.in/go-playground/validator.v8"
)

// Endpointer holds all context for the endpoint to use
type Endpointer struct {
	databaser *db.Databaser
}

// NewEndpointer returns a new endpointer to be used
func NewEndpointer(d *db.Databaser) *Endpointer {
	return &Endpointer{databaser: d}
}

// Success returns generic success message
func Success(message string, c *gin.Context) {
	c.JSON(
		http.StatusCreated,
		gin.H{
			"status":  true,
			"message": message,
		},
	)
}

// SuccessOK returns generic success message with an OK status
func SuccessOK(message string, c *gin.Context) {
	c.JSON(
		http.StatusOK,
		gin.H{
			"status":  true,
			"message": message,
		},
	)
}

// Error returns generic error message
func Error(e string, c *gin.Context) {
	c.JSON(
		http.StatusBadRequest,
		gin.H{
			"status":  false,
			"message": e,
		},
	)
}

// NotFound returns generic error message
func NotFound(e string, c *gin.Context) {
	c.JSON(
		http.StatusNotFound,
		gin.H{
			"status":  false,
			"message": e,
		},
	)
}

// HandleError formats binding errors to a single FieldError
func HandleError(err error, c *gin.Context) {
	errs, ok := err.(validator.ValidationErrors)
	if ok {
		for _, valErr := range errs {
			if valErr.Tag == "required" {
				Error(fmt.Sprintf("You are missing a required field: %s", valErr.Field), c)
			} else if valErr.Tag == "url" {
				Error(fmt.Sprintf("You are passing an invalid value for field: %s", valErr.Field), c)
			} else {
				Error(fmt.Sprintf("Validation error due to the following tag '%s': %s", valErr.Tag, valErr.Field), c)
			}
		}
	} else {
		Error(err.Error(), c)
	}
}
