.PHONY: test test-unit test-integration build init-local-dynamodb serve deploy deploy-cicd

test: test-unit test-integration

test-unit:
	go generate ./...
	ginkgo -r -skipPackage=integration -keepGoing

test-integration:
	ginkgo -r integration

build:
	sam build

init-local-dynamodb:
	./scripts/init-local-dynamodb.sh

serve: build
	sam local start-api --docker-network sam

deploy: build
	sam deploy

deploy-cicd: build
	sam deploy --no-confirm-changeset
