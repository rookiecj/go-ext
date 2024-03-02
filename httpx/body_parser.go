package httpx

import (
	"encoding/csv"
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"reflect"
	"slices"
	"strconv"
	"strings"
)

var BodyParsers = map[string]BodyParser{
	"text/plain":       TextBodyParser,
	"application/json": JsonBodyParser,
	// application/octet-stream
	"text/xml":        XmlBodyParser,
	"application/xml": XmlBodyParser,
	"text/csv":        NewCsvBodyParser(false, ',', '#'),
}

type BodyParser func(buf io.Reader, bodyPtr any) (err error)
type BodyParser2 func(buf io.Reader, bodyPtr any) (body any, err error)

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

func XmlBodyParser(buf io.Reader, bodyPtr any) (err error) {
	// body should be pointer to a type
	data, err := io.ReadAll(buf)
	if err != nil {
		return fmt.Errorf("error while reading: %s", err)
	}
	err = xml.Unmarshal(data, bodyPtr)
	return
}

// CsvBodyParser parses text/csv content type
// bodyPtr should be slice of a type, []myStruct
func CsvBodyParser(buf io.Reader, bodyPtr any) (err error) {
	bodyType := reflect.TypeOf(bodyPtr)
	if bodyType.Kind() != reflect.Slice {
		return errors.New("bodyPtr is not slice")
	}
	eleType := bodyType.Elem()
	if eleType.Kind() != reflect.Struct {
		return errors.New("element is not struct")
	}
	valType := reflect.ValueOf(eleType)
	for fidx := 0; fidx < valType.NumField(); fidx++ {

	}

	r := csv.NewReader(buf)
	r.Comment = '#'

	//// first line is header
	//var headers []string
	//if headers, err = r.Read(); err != nil {
	//	return errors.New("empty csv")
	//}

	for err == nil {

	}
	return
}

func NewCsvBodyParser(hasHeader bool, delim rune, comment rune) BodyParser {
	type fieldInfo struct {
		fname  string
		ftype  reflect.Type
		column string
		idx    int
	}

	// bodyPtr should be a pointer to a slice of a struct type
	csvBodyParser := func(buf io.Reader, bodyPtr any) (err error) {
		bodyType := reflect.TypeOf(bodyPtr)
		if bodyType.Kind() != reflect.Ptr {
			return errors.New("bodyPtr is not a pointer")
		}
		sliceType := bodyType.Elem()
		if sliceType.Kind() != reflect.Slice {
			return errors.New("bodyPtr is not a pointer to slice")
		}
		structType := sliceType.Elem()
		if structType.Kind() != reflect.Struct {
			return errors.New("element is not struct")
		}
		slice := reflect.MakeSlice(sliceType, 0, 16) // 초기 길이 0, 용량 16
		//valType := reflect.ValueOf(structType)
		var fields []fieldInfo
		// csv는 struct embedding를 표현할수없으므로 고려하지 않는다.
		for fidx := 0; fidx < structType.NumField(); fidx++ {
			field := structType.Field(fidx)

			columnName := field.Name
			if len(field.Tag) > 0 {
				if tagName := field.Tag.Get("csv"); len(tagName) > 0 {
					columnName = tagName
				} else if tagName = field.Tag.Get("json"); len(tagName) > 0 {
					columnName = tagName
				}
			}
			fields = append(fields, fieldInfo{
				fname:  field.Name,
				ftype:  field.Type,
				column: columnName,
				idx:    fidx,
			})
		}

		r := csv.NewReader(buf)
		r.Comma = delim
		r.Comment = comment
		r.TrimLeadingSpace = true
		r.LazyQuotes = true

		// first line is header
		var record []string
		if hasHeader {
			if record, err = r.Read(); err != nil {
				return errors.New("empty csv")
			}
			for fidx := 0; fidx < len(fields); fidx++ {
				if idx := slices.Index(record, fields[fidx].column); idx != -1 {
					fields[fidx].idx = idx
				}
			}
		}

		for record, err = r.Read(); err == nil; record, err = r.Read() {
			// struct type을 생성해서
			newStructValue := reflect.New(structType).Elem()
			// 각 필드값을 설정한다.
			for fidx, field := range fields {
				fieldValue := newStructValue.FieldByName(field.fname)
				if !fieldValue.IsValid() || !fieldValue.CanSet() {
					continue
				}
				switch field.ftype.Kind() {
				case reflect.Int8, reflect.Int16, reflect.Int, reflect.Int32, reflect.Int64:
					if val, err := strconv.Atoi(strings.TrimSpace(record[fidx])); err == nil {
						fieldValue.SetInt(int64(val))
					}
				case reflect.Uint8, reflect.Uint16, reflect.Uint, reflect.Uint32, reflect.Uint64:
					if val, err := strconv.ParseUint(strings.TrimSpace(record[fidx]), 10, 64); err == nil {
						fieldValue.SetUint(val)
					}
				case reflect.Float32, reflect.Float64:
					if val, err := strconv.ParseFloat(strings.TrimSpace(record[fidx]), 64); err == nil {
						fieldValue.SetFloat(val)
					}
				default: // string
					//fieldValue.Set(reflect.ValueOf(record[fidx]))
					fieldValue.SetString(strings.TrimSpace(record[fidx]))
				}
			}
			// slice에 추가한다.
			slice = reflect.Append(slice, newStructValue)
		}
		if err == io.EOF {
			err = nil
		}
		// slice를 리턴
		bodyValue := reflect.ValueOf(bodyPtr)
		bodyValue.Elem().Set(slice)
		return
	}
	return csvBodyParser
}
