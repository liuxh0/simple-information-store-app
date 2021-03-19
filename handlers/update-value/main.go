package main

import (
	"fmt"

	"simple-information-store-app/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]
	value := request.Body

	_, err := service.UpdateInfo(id, value)
	switch err := err.(type) {
	case nil:
		break
	case service.ValueTooLongError:
		return events.APIGatewayProxyResponse{
			StatusCode: 400,
			Body:       err.Error(),
		}, nil
	case service.InfoNotFoundError:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	default:
		fmt.Printf("Error when updating item: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
	}, nil
}

func main() {
	lambda.Start(handler)
}
