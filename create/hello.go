package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

type Item struct {
	UserID  int    `dynamodbav:"userid" json:userid`
	Address string `dynamodbav:"address" json:address`
	Email   string `dynamodbav:"email" json:email`
	Gender  string `dynamodbav:"gender" json:gender`
	Name    string `dynamodbav:"name" json:name`
}

type Response struct {
	RequestMethod string `json:RequestMethod`
	Result        Item   `json:Result`
}

func handler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	method := request.HTTPMethod

	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: err.Error(),
			StatusCode: 500,
		}, err
	}

	db := dynamodb.New(sess)

	reqBody := request.Body
	resBodyJSONBytes := ([]byte)(reqBody)
	item := Item{}
	if err := json.Unmarshal(resBodyJSONBytes, &item); err != nil {
		return events.APIGatewayProxyResponse{
			Body: err.Error(),
			StatusCode: 500,
		}, err
	}
	inputAV, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
		}, err
}
	input := &dynamodb.PutItemInput{
		TableName: aws.String("user"),
		Item: inputAV,
	}

	log.Println(input)
	log.Println(reqBody)
	_, err = db.PutItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
		}, err
}

	res := Response {
		RequestMethod: method,
	}
	jsonBytes, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		StatusCode: 200,
	}, nil
}

func main () {
	lambda.Start(handler)
}
