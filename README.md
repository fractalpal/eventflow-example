[![Go Report Card](https://goreportcard.com/badge/github.com/fractalpal/eventflow-example)](https://goreportcard.com/report/github.com/fractalpal/eventflow-example)
[![CircleCI](https://circleci.com/gh/fractalpal/eventflow-example/tree/master.svg?style=svg)](https://circleci.com/gh/fractalpal/eventflow-example/tree/master)

## eventflow showcase project
This is an example of usage of eventflow and implementation of event sourcing pattern in go.

Payments HTTP REST API (json) implemented using CQRS and Event Sourcing pattern.

CQRS gives us separation of concerns for our writes and reads. 
This way we can easily scale traffic for creations and queries independently.

Event Sourcing gives us all information about changes that happened in our service.
This way we've got detailed view what's changes in our system (audit log). 
We can create state at particular time if needed or re-apply from beginning to get current state

##### Notice
This demo does not implement all unit and/or integration/e2e tests. 
Some of them are just for showing of how the may look like. 

### How to run
Command part uses postgres as underlying storage.
Run docker image:
```
docker run --rm --name es-postgres -d -e POSTGRES_DB=events -e POSTGRES_USER=usr -e POSTGRES_PASSWORD=pwd -p 5432:5432 postgres:11.2
```

Query part uses mongodb as underlying storage.
Run docker image:
```
docker run --rm --name es-mongo -d -e MONGO_INITDB_ROOT_USERNAME=usr -e MONGO_INITDB_ROOT_PASSWORD=pwd -p 27017:27017 mongo:3.4
```

### Build and run
````
# build
$ go build
# run 
$ ./eventflow-example
````


### Payments REST API

#### Command Part

###### Create new payment
```
# returns 'id' in response header: 'Location'
curl -X POST http://localhost:8080 -d '{"type":"instant", "attributes":{"amount":"100.50","currency":"EUR","beneficiary_party":{"account_name":"Ben","account_number":"123"},"debtor_party":{"account_name":"Deb","account_number":"987"},"payment_id":"92aaf311-a2fe-4022-86ef-162f314149df","payment_type":"credit_card","processing_date":"2019-03-13T21:20:57+01:00","reference":"f1eef151-87c0-4aac-afbc-3e68c8f5807c"}}' -v
```
###### Update beneficiary party for 'id'
```
curl -X PUT http://localhost:8080/7db061e8-81c5-499d-9f93-321cde3575be/beneficiary-party -d '{"account_name":"Mr. Ben","account_number":"123456789"}' -v
```

###### Update debtor party for 'id'
```
curl -X PUT http://localhost:8080/7db061e8-81c5-499d-9f93-321cde3575be/debtor-party -d '{"account_name":"Account Inc.","account_number":"987654321"}' -v
```

###### Delete payment
```
curl -X DELETE http://localhost:8080/7db061e8-81c5-499d-9f93-321cde3575be -v
```

#### Query Part
###### Get payment for 'id'
```
curl -X GET http://localhost:8081/7db061e8-81c5-499d-9f93-321cde3575be -v
```
###### Get payments list
```
# use query 'page' and 'limit' for paging
curl -X GET http://localhost:8081?page=1&limit=10 -v
```

