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

## Assumptions

  - Using a PostgreSQL database as a datastore (Install Postgres locally)
  - The ENVs in this are just used as an example, real world you'd probably use a vault or store them only on the deployed environment

## Commands

### Creating DB / Running Migrations & Seeds

    make install

### If you add migrations, to execute them on your target DB

    make migrate

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

