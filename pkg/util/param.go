package util

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"
	"reflect"
)

// TODO: tag这一块可能要添加反序列化名和反序列化列表的功能
func ReadParamsFromQuery[T any](queryParams url.Values) (*T, error) {
	params := new(T)
	missing := make([]string, 0)
	val := reflect.ValueOf(params).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		paramName := CamelToSnake(field.Name)
		paramValue := queryParams.Get(paramName)
		if paramValue == "" {
			if field.Tag.Get("query") != "" {
				paramValue = field.Tag.Get("query")
			} else {
				missing = append(missing, paramName)
			}
		}
		val.Field(i).SetString(paramValue)
	}

	if len(missing) != 0 {
		return nil, fmt.Errorf("missing parameters: %v", missing)
	}

	return params, nil
}

func ReadParamsFromBody[T any](bodyReader io.ReadCloser) (*T, error) {
	body, err := io.ReadAll(bodyReader)
	if err != nil {
		return nil, err
	}

	params := new(T)
	err = json.Unmarshal(body, params)
	if err != nil {
		return nil, err
	}
	return params, nil
}
