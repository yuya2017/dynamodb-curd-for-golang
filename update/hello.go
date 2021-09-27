package main

import (
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
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

	reqBody := request.Body
	reqBodyJSONBytes := ([]byte)(reqBody)
	item := Item{}
	if err := json.Unmarshal(reqBodyJSONBytes, &item); err != nil {
		return events.APIGatewayProxyResponse{
			Body:       err.Error(),
			StatusCode: 500,
		}, err
	}

	log.Println(reqBody)
	log.Println(item)

	update := expression.UpdateBuilder{}
	if address := item.Address; address != "" {
		update = update.Set(expression.Name("address"), expression.Value(address))
	}
	if email := item.Email; email != "" {
		update = update.Set(expression.Name("email"), expression.Value(email))
	}
	if gender := item.Gender; gender != "" {
		update = update.Set(expression.Name("gender"), expression.Value(gender))
	}
	if name := item.Name; name != "" {
		update = update.Set(expression.Name("name"), expression.Value(name))
	}
	expr, err := expression.NewBuilder().WithUpdate(update).Build()
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
		}, err
	}

	input := &dynamodb.UpdateItemInput{
		TableName: aws.String("user"),
		Key: map[string]*dynamodb.AttributeValue{
				"userid": {
						N: aws.String(pathParam),
				},
		},
		ExpressionAttributeNames: expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		UpdateExpression: expr.Update(),
	}

	_, err = db.UpdateItem(input)
	if err != nil {
		return events.APIGatewayProxyResponse{
				Body:       err.Error(),
				StatusCode: 500,
		}, err
	}

	res := Response {
		RequestMethod: method,
	}

	log.Println(res)

	jsonBytes, _ := json.Marshal(res)

	return events.APIGatewayProxyResponse{
		Body: string(jsonBytes),
		StatusCode: 200,
	}, nil

}

func main () {
	lambda.Start(handler)
}
