package db

import (
	"2dal/core"
	"strings"

	"github.com/getsentry/sentry-go"
	"github.com/go-playground/validator/v10"
	"github.com/kamva/mgm/v3"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var DB_CONNECTION_STRING = "mongodb://localhost:27017"
var DB_NAME = "2dal"

func GetConnection() error{
	connectionString := core.Env("DB_CONNECTION_STRING", DB_CONNECTION_STRING)
	dbName := core.Env("DB_NAME", DB_NAME)

	err := mgm.SetDefaultConfig(nil, dbName, options.Client().ApplyURI(connectionString))
	if err != nil {
		return err
	}

	return nil
}

func IsUniqueValidator() validator.Func {
	return func(fl validator.FieldLevel) bool {
		GetConnection()

		value, _ := fl.Field().Interface().(string)
		field := strings.Split(fl.Param(), "@") // field@document
		
		result := []map[string]interface{}{}
		
		document := mgm.CollectionByName(field[1] + "s")
		err := document.SimpleFind(&result, bson.D{{Key: field[0], Value: value}})

		if err != nil{
			sentry.CaptureException(err)
			return false
		}

		return len(result) == 0;
	}
}