package httpx

import (
	"encoding/json"
	"strings"
)

var bodyParsers = map[string]BodyParser{
	"application/json": JsonBodyParser,
}

type BodyParser func(data []byte, bodyPtr any) (err error)

func getBodyParser(contentType string) BodyParser {
	lowerType := strings.ToLower(contentType)
	for key, parser := range bodyParsers {
		if strings.HasPrefix(lowerType, key) {
			return parser
		}
	}
	return nil
}

func JsonBodyParser(data []byte, bodyPtr any) (err error) {
	// body should be pointer to a type
	err = json.Unmarshal(data, bodyPtr)
	return
}
