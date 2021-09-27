package main

import (
	"encoding/json"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
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

	sess, err := session.NewSession()
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
		}, err
	}

	db := dynamodb.New(sess)

	deleteParam := &dynamodb.DeleteItemInput{
		TableName: aws.String("user"),
		Key: map[string]*dynamodb.AttributeValue{
				"userid": {
						// Nはnumber型の意味
						N: aws.String(pathParam),
				},
		},
	}
	_, err = db.DeleteItem(deleteParam)
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
