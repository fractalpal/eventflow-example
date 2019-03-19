#!/bin/bash

set -o nounset
set -o errexit

# create container
echo "creating docker container with test postgres database"
ID=$(docker run --rm --name es-postgres-test -d -e POSTGRES_DB=events_test -e POSTGRES_USER=usr -e POSTGRES_PASSWORD=pwd -p 5555:5432 postgres:11.2)

echo "Waiting for postgres to be ready"
while ! nc -z localhost 5555; do sleep 1; done;
sleep 2

# trap function to cleaning stuff
function finish {
    # kill container
    echo "cleaning after tests"
    docker kill $ID
}
trap finish EXIT

echo "staring client integration test"
export PAYMENT_POSTGRES_MIGRATIONS_PATH=file://../adapters/store/migrations
export PAYMENT_POSTGRES_URL=postgres://usr:pwd@localhost:5555/events_test?sslmode=disable
go test ./... -tags integration
