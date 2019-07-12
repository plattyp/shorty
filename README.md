# Shorty - Minimal Persisted URL Shortener

A pretty simple URL shortener written in Golang

## Getting Started

    createdb shorty
    make migrate

## Assumptions

  - Using a PostgreSQL database as a datastore (Install Postgres locally)
  - The ENVs in this are just used as an example, real world you'd probably use a vault or store them only on the deployed environment

## Structure

## Creating DB / Running Migrations & Seeds

    createdb shorty
    make migrate

## Building It

    make

## Running It (exposed on port 4100 by default)

    make run
