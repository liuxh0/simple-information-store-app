package main

import (
	"fmt"

	"simple-information-store-app/internal/service"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	id := request.PathParameters["id"]

	info, err := service.NewInfoService().GetInfo(id)
	switch err := err.(type) {
	case nil:
		break
	case service.InfoNotFoundError:
		return events.APIGatewayProxyResponse{
			StatusCode: 404,
		}, nil
	default:
		fmt.Printf("Error when retrieving item: %s\n", err.Error())
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       info.Value,
	}, nil
}

func main() {
	lambda.Start(handler)
}
