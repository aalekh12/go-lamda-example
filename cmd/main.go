package main

import (
	"aws-lambda-in-go-lang/pkg/handlers"
	"fmt"
	"log"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var(
	dynaclient dynamodbiface.DynamoDBAPI
)

func main(){

	fmt.Println("Hello")
	region:=os.Getenv("AWS_REGION")
	awssession,err:=session.NewSession(&aws.Config{
		Region: aws.String(region),
	})
	if err!=nil {
		log.Println("Error while session creation ",err)
	}

	dynaclient=dynamodb.New(awssession)
    lambda.Start(handler)
}

const tablename="LamdainGoUser"

func handler(req events.APIGatewayProxyRequest)(*events.APIGatewayProxyResponse,error){
   switch req.HTTPMethod{
   case "GET":
			return handlers.GetUser(req,tablename,dynaclient)
   case "POST":
			return handlers.CreateUser(req,tablename,dynaclient)
   case "PUT":
			return handlers.UpdateUser(req,tablename,dynaclient)
   case "DELETE":
			return handlers.DeleteUser(req,tablename,dynaclient)
   default:
			return handlers.UnhadledMethod()

   }
}