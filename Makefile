.PHONY: test test-unit test-integration build init-local-dynamodb serve deploy

test: test-unit test-integration

test-unit:
	ginkgo -r -skipPackage=integration

test-integration:
	ginkgo -r integration

build:
	sam build

init-local-dynamodb:
	./scripts/init-local-dynamodb.sh

serve: build
	sam local start-api

deploy: build
	sam deploy
