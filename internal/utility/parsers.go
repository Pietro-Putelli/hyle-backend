package utility

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"reflect"
	"strconv"

	"github.com/aws/aws-lambda-go/events"
)

// convertParams converts the string parameters to the appropriate type
func convertParams(params map[string]string) map[string]interface{} {
	convertedParams := make(map[string]interface{})
	for key, value := range params {
		if uintVal, err := strconv.ParseInt(value, 10, 64); err == nil {
			convertedParams[key] = int(uintVal)
		} else {
			convertedParams[key] = value
		}
	}
	return convertedParams
}

// validateParams validates the parameters of the target struct
func validateParams(item interface{}) error {
	validate := validator.New(validator.WithRequiredStructEnabled())
	if err := validate.Struct(item); err != nil {
		return err
	}
	return nil
}

// ParseRequestBody parses the request body and sets it to the target struct
func ParseRequestBody(request events.APIGatewayProxyRequest, target interface{}) error {
	if err := json.Unmarshal([]byte(request.Body), target); err != nil {
		return err
	}
	if err := validateParams(target); err != nil {
		return err
	}
	return nil
}

// ParsePathParams parses the path parameters from the request and sets them to the target struct
func ParsePathParams(request events.APIGatewayProxyRequest, target interface{}) error {
	pathParams := request.PathParameters

	convertedParams := convertParams(pathParams)

	v := reflect.ValueOf(target).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Tag.Get("path")
		if fieldName == "" {
			fieldName = v.Type().Field(i).Name
		}
		if val, ok := convertedParams[fieldName]; ok {
			field.Set(reflect.ValueOf(val))
		}
	}

	if err := validateParams(target); err != nil {
		return err
	}
	return nil
}

// ParseQueryParams parses the query parameters from the request and sets them to the target struct
func ParseQueryParams(request events.APIGatewayProxyRequest, target interface{}) error {
	queryParams := request.QueryStringParameters

	convertedParams := convertParams(queryParams)

	v := reflect.ValueOf(target).Elem()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Tag.Get("json")
		if fieldName == "" {
			fieldName = v.Type().Field(i).Name
		}
		if val, ok := convertedParams[fieldName]; ok {
			field.Set(reflect.ValueOf(val))
		}
	}

	if err := validateParams(target); err != nil {
		return err
	}
	return nil
}
