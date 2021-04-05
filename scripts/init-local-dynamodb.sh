docker run -d -p 8000:8000 amazon/dynamodb-local
aws dynamodb create-table --cli-input-json file://local-dynamodb-value-table.json --endpoint-url http://localhost:8000
