package httpx

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"reflect"
)

var BodyParsers = map[string]BodyParser{
	"text/plain":       TextBodyParser,
	"application/json": JsonBodyParser,
	// application/octet-stream
}

type BodyParser func(buf io.Reader, bodyPtr any) (err error)

// TextBodyParser parses text/plain content type
func TextBodyParser(buf io.Reader, bodyPtr any) (err error) {

	bodyType := reflect.TypeOf(bodyPtr).Elem()
	if bodyType.Kind() != reflect.String {
		return errors.New("bodyPtr is not string")
	}
	valType := reflect.ValueOf(bodyPtr)
	valRef := reflect.Indirect(valType)

	var data []byte
	data, err = io.ReadAll(buf)
	if err != nil {
		return fmt.Errorf("error while reading: %s", err)
	}

	valRef.SetString(string(data))
	return nil
}

// JsonBodyParser parses application/json content type
func JsonBodyParser(buf io.Reader, bodyPtr any) (err error) {
	// body should be pointer to a type
	data, err := io.ReadAll(buf)
	if err != nil {
		return fmt.Errorf("error while reading: %s", err)
	}
	err = json.Unmarshal(data, bodyPtr)
	return
}
