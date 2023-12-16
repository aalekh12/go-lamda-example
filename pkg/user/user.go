package user

import (
	"aws-lambda-in-go-lang/pkg/validators"
	"encoding/json"
	"errors"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbiface"
)

var (
	ErrorFailedToUnmarshalRecord = "failed to unmarshal record"
	ErrorFailedToFetchRecord     = "failed to fetch record"
	ErrorInvalidUserData         = "invalid user data"
	ErrorInvalidEmail            = "invalid email"
	ErrorCouldNotMarshalItem     = "could not marshal item"
	ErrorCouldNotDeleteItem      = "could not delete item"
	ErrorCouldNotDynamoPutItem   = "could not dynamo put item error"
	ErrorUserAlreadyExists       = "User already exists"
	ErrorUserDoesNotExists       = "User does not exist"
)

type User struct {
	Email     string `json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

func Fetchuser(emai, tablename string, dyanaclient dynamodbiface.DynamoDBAPI) (*User, error) {
	input := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(emai),
			},
		},
		TableName: aws.String(tablename),
	}

	res, err := dyanaclient.GetItem(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	item := new(User)

	err = dynamodbattribute.UnmarshalMap(res.Item, item)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}

	return item, nil
}

func FetchUsers(tablename string, dynaclient dynamodbiface.DynamoDBAPI) (*[]User, error) {
	input := &dynamodb.ScanInput{
		TableName: aws.String(tablename),
	}

	result, err := dynaclient.Scan(input)
	if err != nil {
		return nil, errors.New(ErrorFailedToFetchRecord)
	}

	items := new([]User)
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, items)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	return items, nil
}

func CreateUser(req events.APIGatewayProxyRequest, tablename string, dyanaclient dynamodbiface.DynamoDBAPI) (
	*User,
	error,
) {
	var u User

	if err := json.Unmarshal([]byte(req.Body), &u); err != nil {
		return nil, errors.New(ErrorInvalidUserData)
	}

	if !validators.IsEmailValid(u.Email) {
		return nil, errors.New(ErrorInvalidEmail)
	}

	currentuser, _ := Fetchuser(u.Email, tablename, dyanaclient)
	if currentuser != nil && len(currentuser.Email) != 0 {
		return nil, errors.New(ErrorUserAlreadyExists)
	}

	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = dyanaclient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil
}

func UpdateUser(req events.APIGatewayProxyRequest, tablename string, dyanaclient dynamodbiface.DynamoDBAPI) (
	*User, error,
) {
	var u User
	err := json.Unmarshal([]byte(req.Body), &u)
	if err != nil {
		return nil, errors.New(ErrorFailedToUnmarshalRecord)
	}
	currentuser, _ := Fetchuser(u.Email, tablename, dyanaclient)
	if currentuser != nil && len(currentuser.Email) == 0 {
		return nil, errors.New(ErrorUserDoesNotExists)
	}

	//SAVE USER
	av, err := dynamodbattribute.MarshalMap(u)
	if err != nil {
		return nil, errors.New(ErrorCouldNotMarshalItem)
	}

	input := &dynamodb.PutItemInput{
		Item:      av,
		TableName: aws.String(tablename),
	}

	_, err = dyanaclient.PutItem(input)
	if err != nil {
		return nil, errors.New(ErrorCouldNotDynamoPutItem)
	}

	return &u, nil
}

func DeleteUser(req events.APIGatewayProxyRequest, tablename string, dynaclient dynamodbiface.DynamoDBAPI) error {
	email := req.QueryStringParameters["email"]
	input := &dynamodb.DeleteItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"email": {
				S: aws.String(email),
			},
		},
		TableName: aws.String(tablename),
	}

	_, err := dynaclient.DeleteItem(input)
	if err != nil {
		return errors.New(ErrorCouldNotDeleteItem)
	}
	return nil

}
