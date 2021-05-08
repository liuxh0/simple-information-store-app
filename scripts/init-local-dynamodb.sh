docker run --name dynamodb --network sam -p 8000:8000 -d amazon/dynamodb-local
aws dynamodb create-table --cli-input-json file://local-dynamodb-value-table.json --endpoint-url http://localhost:8000 --no-cli-pager
