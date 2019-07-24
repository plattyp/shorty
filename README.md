# Shorty - Minimal Persisted URL Shortener [![Build Status](https://travis-ci.org/plattyp/shorty.svg?branch=master)](https://travis-ci.org/plattyp/shorty)

A pretty simple URL shortener written in Golang

## Dependencies

    Golang (1.12 or greater)
    PostgreSQL (10 or greater)

## Getting Started

    # Will create the necessary database locally and run all migrations
    make install

    # Will build the Golang app for the first time and pull all dependent modules
    make build

## Envir

  - Using a PostgreSQL database as a datastore (Install Postgres locally)
  
## Options

  - By default, we'll generate shortened urls with 10 random characters. This can be modified by setting `URL_SLUG_LENGTH` to a different number. With 10 characters, there are 327,234,915,316,109,000 [permutations](https://stattrek.com/online-calculator/combinations-permutations.aspx)
  
## Deploying

  - The `.env` file in stored in this repo is just used for local development. On your deployed application, you should use environment variables
  - When deploying, your destination server will need to have a `DATABASE_URL` environment variable. You can then run `make migrate-deployed` to invoke the up migrations on that `DATABASE_URL`.

## Commands

### Creating DB / Running Migrations & Seeds

    make install

### If you need to run migrations, to execute them on your local DB

    make migrate

### If you need to run migrations on your deployed server, this will run them against the DATABASE_URL environment variable

    make migrate-deployed

### Building It

    make

### Running It (exposed on port 4100 by default)

    make run

### Initializing Test Environment

    make init-test

### Testing It (make sure to initialize prior to running tests)

    make test

## Endpoints

### GET /

Used as a health check with no authentication required.

### POST /api/shorten

Accepts a single body parameter of `url`. This requires authorization using the `USERNAME` and `PASSWORD` environment variables. If successful, it'll return back information about the shortened URL to be able to use.

Sample Request:
```json
{
  "url": "https://www.google.com/"
}
```

Sample Response:
```json
{
  "id": 1,
  "shortened_url": "http://localhost:4100/I4s1Wkk89tAOEzcXFjJ3"
}
```

### GET /:slug

Used for redirecting the shortened URL slug into the original destination URL.
