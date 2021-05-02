package main

import (
	"encoding/json"
	"fmt"

	"simple-information-store-app/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/google/uuid"
)

var infoCreator service.InfoCreator = service.NewInfoService()

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	body := request.Body

	// Generate an Id
	id := uuid.New().String()
	value := body
	info, err := infoCreator.CreateInfo(id, value)

	switch err := err.(type) {
	case nil:
		break
	case service.ValueTooLongError:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	default:
		fmt.Printf("Error when putting item: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	responseBody := map[string]string{
		"id": info.ID,
	}
	responseBodyBytes, _ := json.Marshal(responseBody)
	return events.APIGatewayProxyResponse{
		StatusCode: 201,
		Body:       string(responseBodyBytes),
	}, nil
}

func main() {
	lambda.Start(handler)
}
