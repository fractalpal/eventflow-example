lint:
	@if [ -n "$$(find . -type f -name \*.go -exec goimports -l {} \;)" ]; then \
		echo "Go code is not formatted:" ; \
		find . -type f -name \*.go -exec goimports -l {} \; ; \
		exit 1; \
	fi

format:
	find . -type f -name \*.go -exec goimports -w {} \;

test:
	go test ./... --cover -tags test

test_integration:
	@bash ./payment/integration-tests/setup_db_and_run.sh
	@bash ./payment-query/integration-tests/setup_db_and_run.sh

