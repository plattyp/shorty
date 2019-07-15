package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/suite"
)

type MainTestSuite struct {
	suite.Suite
}

func getCredentials() (string, string) {
	return os.Getenv("USERNAME"), os.Getenv("PASSWORD")
}

func performGETRequest(r http.Handler, path string) *httptest.ResponseRecorder {
	req, _ := http.NewRequest("GET", path, nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

func (suite *MainTestSuite) SetupTest() {
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/shorty-test")
}

func (suite *MainTestSuite) TestIndexReturnsStandardResult() {
	// Grab our router
	router := SetupRouter()

	// Perform a GET request with that handler.
	w := performGETRequest(router, "/")

	//assert.Equal(suite.T(), http.StatusOK, w.Code)
	suite.Equal(http.StatusOK, w.Code)

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

	fmt.Println(w.Body.String())
}

func TestMainTestSuite(t *testing.T) {
	suite.Run(t, new(MainTestSuite))
}
