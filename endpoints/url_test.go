package endpoints_test

import (
	"os"
	"shorty/endpoints"
	"testing"

	"shorty/accessors"
	"shorty/db"

	"github.com/stretchr/testify/suite"
)

type UrlTestSuite struct {
	suite.Suite
}

func (suite *UrlTestSuite) SetupTest() {
	if os.Getenv("TRAVIS") != "" {
		return
	}

	os.Setenv("DATABASE_URL", "postgres://localhost:5432/shorty-test")
}

func createShortenedURL(slug, url string) error {
	dbConn, err := db.NewDatabaser(os.Getenv("DATABASE_URL"))
	if err != nil {
		return err
	}
	defer dbConn.Close()

	u := accessors.URLDataAccessor{Databaser: dbConn}
	_, err = u.CreateURL(url, slug)
	if err != nil {
		return err
	}

	return nil
}

func (suite *UrlTestSuite) TestRandomSlugGeneratorReturnsRandomStringOfNCharacters() {
	result := endpoints.RandomGeneratedSlug(20)
	suite.Equal(20, len(result))

	result = endpoints.RandomGeneratedSlug(31)
	suite.Equal(31, len(result))
}

func (suite *UrlTestSuite) TestNewSlugAvoidsCollisionByCheckingAgainstUsedSlugs() {
	dbConn, err := db.NewDatabaser(os.Getenv("DATABASE_URL"))
	if !suite.NoError(err) {
		return
	}
	defer dbConn.Close()

	u := accessors.URLDataAccessor{Databaser: dbConn}

	generatedSlugs := map[string]bool{}
	for i := 0; i < 1000; i++ {
		slug, err := endpoints.NewSlug(u, 10)
		if !suite.NoError(err) {
			return
		}

		_, ok := generatedSlugs[slug]
		suite.Equal(false, ok)
		generatedSlugs[slug] = true
	}
}

func (suite *UrlTestSuite) TestNewSlugReturnsErrorIfMaxRetriesAreMadeAndNoUniqueSlugCanBeGenerated() {
	randomSlug := endpoints.RandomGeneratedSlug(20)

	// This will lock the random generator to only return the same thing
	os.Setenv("STATIC_RANDOM_SLUG", randomSlug)

	// Create the URL
	err := createShortenedURL(randomSlug, "https://www.google.com")
	if !suite.NoError(err) {
		return
	}

	dbConn, err := db.NewDatabaser(os.Getenv("DATABASE_URL"))
	if !suite.NoError(err) {
		return
	}
	defer dbConn.Close()

	u := accessors.URLDataAccessor{Databaser: dbConn}

	_, err = endpoints.NewSlug(u, 10)
	suite.Error(err)
	suite.Equal(endpoints.ErrUnableToGenerateUniqueSlug, err)

	// Reset the generator
	os.Unsetenv("STATIC_RANDOM_SLUG")
}

func TestUrlTestSuite(t *testing.T) {
	suite.Run(t, new(UrlTestSuite))
}
