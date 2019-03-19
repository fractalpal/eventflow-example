### Payments Query Part
Payment query service HTTP REST API used for fetching of the payment model.
It's implemented using Event Sourcing patter under the hood.

#### Configuration
Used env vars:
```
# possible env vars, below shown default vaules:
export QUERY_LISTEN_HOST="localhost"
export QUERY_LISTEN_PORT="8081"
export QUERY_MONGO_URL="mongodb://usr:pwd@localhost:27017"
```
