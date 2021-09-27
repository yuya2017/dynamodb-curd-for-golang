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
	pathParam := request.PathParameters["pathparam"]

	log.Println(pathParam)
	log.Println(request)
	log.Println(request.PathParameters["pathparam"])
	log.Println(request.Body)

	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
		}, err
	}

	db := dynamodb.New(sess)

	getparam := &dynamodb.GetItemInput{
		TableName: aws.String("user"),
		Key: map[string]*dynamodb.AttributeValue{
			"userid": {
				N: aws.String(pathParam),
			},
		},
	}

	result, err := db.GetItem(getparam)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 404,
		}, err
	}

	item := Item{}
	err = dynamodbattribute.UnmarshalMap(result.Item, &item)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: err.Error(),
			StatusCode: 500,
		}, err
	}

	res := Response {
		RequestMethod: method,
		Result: item,
	}
	jsonBytesm, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body: string(jsonBytesm),
		StatusCode: 200,
	}, nil

}

func main () {
	lambda.Start(handler)
}
