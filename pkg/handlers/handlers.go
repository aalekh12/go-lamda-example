package handlers

import (
	"aws-lambda-in-go-lang/pkg/user"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var Errmethodnotallowed = "ErrorMethodNotAllowed"

type ErrorBody struct{
	ErrorMsg *string `json:"error,omitempty"`
}

func GetUser(req events.APIGatewayProxyRequest,tablename string,dyanclient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse,error){
		email:=req.QueryStringParameters["email"]
        if len(email)>0 {
			currentuser,err:= user.Fetchuser(email,tablename,dyanclient)
			if err!=nil {
				return apiresponse(http.StatusBadRequest,ErrorBody{ErrorMsg: aws.String(err.Error())})
			}

			return apiresponse(http.StatusOK,currentuser)

		}
		 
		result,err:=user.FetchUsers(tablename,dyanclient)
		if err!=nil {
			return apiresponse(http.StatusBadRequest,ErrorBody{ErrorMsg: aws.String(err.Error())})
		}

		return apiresponse(http.StatusOK,result)
	}

func CreateUser(req events.APIGatewayProxyRequest,tablename string,dyanacient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse,
	error,
){
	res,err:=user.CreateUser(req,tablename,dyanacient)
	if err!=nil {
		return apiresponse(http.StatusBadRequest,ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
    
	return apiresponse(http.StatusCreated,res)
}

func UpdateUser(req events.APIGatewayProxyRequest,tablename string,dynaclient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse,
	error,
){
	res,err:=user.UpdateUser(req,tablename,dynaclient)
	if err!=nil {
		return apiresponse(http.StatusBadRequest,ErrorBody{ErrorMsg: aws.String(err.Error())})
	}
	return apiresponse(http.StatusOK,res)
}

func DeleteUser(req events.APIGatewayProxyRequest,tablename string,dynaclient dynamodbiface.DynamoDBAPI)(
	*events.APIGatewayProxyResponse,
	error,
){
	user:=user.DeleteUser(req,tablename,dynaclient)
	if user!=nil {
		return apiresponse(http.StatusBadRequest,ErrorBody{ErrorMsg: aws.String(user.Error())})
	}

	return apiresponse(http.StatusOK,nil)
}

func UnhadledMethod()(*events.APIGatewayProxyResponse,error){
	return apiresponse(http.StatusMethodNotAllowed,Errmethodnotallowed)
}