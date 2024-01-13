package httpx

import (
	"net/http"
	"reflect"
)

type Response struct {
	res         *http.Response
	bodyParsers map[string]BodyParser
	Header      map[string][]string
	Status      string
	StatusCode  int
	Data        []byte
}

func (resp *Response) Unmarshal(ptrType any) (err error) {
	// parse body
	if contentTypes, ok := resp.Header["Content-Type"]; ok {
		if bodyParser := getBodyParser(contentTypes[0]); bodyParser != nil {
			err = bodyParser(resp.Data, ptrType)
			if err == nil {
				return nil
			}
		}
	}

	// string
	bodyType := reflect.TypeOf(ptrType).Elem()
	if bodyType.Kind() == reflect.String {
		valType := reflect.ValueOf(ptrType)
		valRef := reflect.Indirect(valType)
		valRef.SetString(string(resp.Data))
		return nil
	}

	return
}
