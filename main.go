package main

import (
	"log"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"errors"
	"fmt"
)

var
	ErrNameNotProvided = errors.New("no name was provided in the HTTP body")
type Item struct {
	MessageId string`json:"messageId"`
}


func Handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	log.Printf("Processing Lambda request %s\n", request.RequestContext.RequestID)
	log.Printf("Request", request)
	if len(request.Body) < 1 {
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2")},
	)
	if err != nil {
		log.Printf("Error to create the session")
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}
	svc := dynamodb.New(sess)

	result, err := svc.GetItem(&dynamodb.GetItemInput{
		TableName: aws.String("test"),
		Key: map[string]*dynamodb.AttributeValue{
			"messageId": {S: aws.String("1234")},
		},

	})

	if err != nil {
		log.Printf("Error to get the item", err)
		return events.APIGatewayProxyResponse{}, ErrNameNotProvided
	}

	item := Item{}

	err = dynamodbattribute.UnmarshalMap(result.Item, &item)

	if err != nil {
		panic(fmt.Sprintf("Failed to unmarshal Record, %v", err))
	}

	log.Printf("Result", item)

	return events.APIGatewayProxyResponse{
		Body:       item.MessageId + " " + request.Body,
		StatusCode: 200,
	}, nil
}

func main()  {
	lambda.Start(Handler)
}

