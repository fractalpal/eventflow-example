### Payments Command Part
Payment service with HTTP REST API used for creation / updates of the payment model.
It's implemented using Event Sourcing patter under the hood. 

#### Configuration
Used env vars:
```
# possible env vars, below shown default vaules:
export PAYMENT_LISTEN_HOST="localhost"
export PAYMENT_LISTEN_PORT="8080"
export PAYMENT_POSTGRES_URL="postgres://usr:pwd@localhost:5432/events?sslmode=disable"
export PAYMENT_POSTGRES_MIGRATIONS_PATH="file://payment/adapters/store/migrations"
```
