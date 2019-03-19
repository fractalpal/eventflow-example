#!/bin/bash

set -o nounset
set -o errexit

# create container
echo "creating docker container with test mongo database"
ID=$(docker run --rm --name es-mongo-test -d -e MONGO_INITDB_ROOT_USERNAME=usr -e MONGO_INITDB_ROOT_PASSWORD=pwd -p 27018:27017 mongo:3.4)

echo "Waiting for mongo to be ready"
while ! nc -z localhost 27018; do sleep 1; done;

# trap function to cleaning stuff
function finish {
    # kill container
    echo "cleaning after tests"
    docker kill $ID
}
trap finish EXIT

echo "staring client integration test"
export PAYMENT_MONGO_URL=mongodb://usr:pwd@localhost:27018
go test ./... -tags integration
