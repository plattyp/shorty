package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"shorty/accessors"
	"shorty/db"

	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

func getCredentials() (string, string) {
	return os.Getenv("USERNAME"), os.Getenv("PASSWORD")
}

func getRandomString() string {
	rand.Seed(time.Now().UnixNano())
	min := 100000000
	max := 999999999
	return fmt.Sprintf("%d", (rand.Intn(max-min) + min))
}

func performGETRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func createShortenedURL(url string) (string, error) {
	randomSlug := getRandomString()

	dbConn, err := db.NewDatabaser(os.Getenv("DATABASE_URL"))
	if err != nil {
		return "", err
	}
	defer dbConn.Close()

	u := accessors.URLDataAccessor{Databaser: dbConn}
	createdURL, err := u.CreateURL(url, randomSlug)
	if err != nil {
		return "", err
	}

	return createdURL.Slug, nil
}

func (suite *MainTestSuite) SetupTest() {
	if os.Getenv("TRAVIS") != "" {
		return
	}

	os.Setenv("DATABASE_URL", "postgres://localhost:5432/shorty-test")
}

func (suite *MainTestSuite) TestIndexReturnsStandardResult() {
	// Grab our router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w := performGETRequest(router, "/")
	suite.Equal(http.StatusOK, w.Code)

	// Convert the JSON response to a map
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if !suite.NoError(err) {
		return
	}

	// Grab the value & whether or not it exists
	statusValue, _ := response["status"].(bool)
	messageValue, _ := response["message"].(string)

	// Make some assertions on the correctness of the response.
	suite.Equal(true, statusValue)
	suite.Equal("All systems go.", messageValue)
}

func (suite *MainTestSuite) TestPOSTShortenWithoutCredentialsReturns401() {
	// Grab our router
	router := SetupRouter()

	// Perform a POST request with that handler.
	payload := `{
    "url": "https://google.com"
  }`
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(payload))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *MainTestSuite) TestPOSTShortenWithInvalidCredentialsReturns401() {
	// Grab our router
	router := SetupRouter()

	// Perform a POST request with that handler.
	payload := `{
    "url": "https://google.com"
  }`
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(payload))
	req.SetBasicAuth("NOT", "REAL")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusUnauthorized, w.Code)
}

func (suite *MainTestSuite) TestPOSTShortenWithValidCredentialsReturns201AndSavesURLIfValidURLProvided() {
	// This is so we can get a static slug returned
	randomSlug := getRandomString()
	os.Setenv("STATIC_RANDOM_SLUG", randomSlug)

	// Grab our router
	router := SetupRouter()

	// Perform a POST request with that handler.
	payload := `{
    "url": "https://google.com"
  }`
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(payload))
	username, password := getCredentials()
	req.SetBasicAuth(username, password)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusCreated, w.Code)

	// Convert the JSON response to a map
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if !suite.NoError(err) {
		return
	}

	// Grab the value & whether or not it exists
	idValue, _ := response["id"].(int)
	shortenedURLValue, _ := response["shortened_url"].(string)

	// Make some assertions on the correctness of the response.
	suite.NotNil(idValue)
	suite.Equal("http://localhost:4100/"+randomSlug, shortenedURLValue)

	// Reset
	os.Unsetenv("STATIC_RANDOM_SLUG")
}

func (suite *MainTestSuite) TestPOSTShortenReturns400IfNoURLIsProvided() {
	// Grab our router
	router := SetupRouter()

	// Perform a POST request with that handler.
	payload := `{}`
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(payload))
	username, password := getCredentials()
	req.SetBasicAuth(username, password)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)

	// Convert the JSON response to a map
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if !suite.NoError(err) {
		return
	}

	// Grab the value & whether or not it exists
	statusValue, _ := response["status"].(bool)
	messageValue, _ := response["message"].(string)

	// Make some assertions on the correctness of the response.
	suite.Equal(false, statusValue)
	suite.Equal("You are missing a required field: URL", messageValue)
}

func (suite *MainTestSuite) TestPOSTShortenReturns400IfInvalidURLIsProvided() {
	// Grab our router
	router := SetupRouter()

	// Perform a POST request with that handler.
	payload := `{
		"url": "TOTAL CRAP"
	}`
	req, _ := http.NewRequest("POST", "/api/shorten", strings.NewReader(payload))
	username, password := getCredentials()
	req.SetBasicAuth(username, password)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	suite.Equal(http.StatusBadRequest, w.Code)

	// Convert the JSON response to a map
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if !suite.NoError(err) {
		return
	}

	// Grab the value & whether or not it exists
	statusValue, _ := response["status"].(bool)
	messageValue, _ := response["message"].(string)

	// Make some assertions on the correctness of the response.
	suite.Equal(false, statusValue)
	suite.Equal("You are passing an invalid value for field: URL", messageValue)
}

func (suite *MainTestSuite) TestGETRedirectReturns404IfNoSlugCorresponds() {
	// Grab our router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w := performGETRequest(router, "/notrealslug")
	suite.Equal(http.StatusNotFound, w.Code)

	// Convert the JSON response to a map
	var response map[string]interface{}
	err := json.Unmarshal([]byte(w.Body.String()), &response)
	if !suite.Nil(err) {
		return
	}

	// Grab the value & whether or not it exists
	statusValue, _ := response["status"].(bool)
	messageValue, _ := response["message"].(string)

	// Make some assertions on the correctness of the response.
	suite.Equal(false, statusValue)
	suite.Equal("Shortened URL not found", messageValue)
}

func (suite *MainTestSuite) TestGETRedirectReturnsAPermanentRedirectToTargetURLIfExists() {
	originalURL := "https://www.google.com"

	slug, err := createShortenedURL(originalURL)
	if !suite.NoError(err) {
		return
	}

	// Grab our router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w := performGETRequest(router, "/"+slug)
	suite.Equal(http.StatusMovedPermanently, w.Code)

	location := w.Header().Get("location")
	suite.Equal(originalURL, location)
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
