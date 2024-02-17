package mapper

import (
	"fmt"
	"reflect"
	"strings"
)

// Mapper is an interface to map a struct to another struct using field tags
// It is useful when you want to map a struct to another struct with different field names
type Mapper interface {
	// Map maps a struct to another struct its field tags
	Map(dest interface{}, src interface{}) error
}

// structMapper is a structMapper implementation
type structMapper struct {
	tagName string
}

// NewMapper returns a new Mapper
// tagName is the name of the struct tag to use
func NewMapper() Mapper {
	return NewMapperWithTag("json")
}

func NewMapperWithTag(tagName string) Mapper {
	return &structMapper{
		tagName: tagName,
	}
}

// Map maps a struct to another struct using field tags
func (m *structMapper) Map(destPtr interface{}, src interface{}) error {

	// 임베딩된 구조체를 지원하기 위해 reflect.ValueOf() 함수를 이용하여 srcValue를 가져온다.
	// srcValue가 포인터인 경우 포인터를 제거한다.
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	// destValue가 포인터가 아닌 경우 에러를 반환한다.
	// destValue가 포인터이면서 구조체가 아닌 경우 에러를 반환한다.
	destValue := reflect.ValueOf(destPtr)
	if destValue.Kind() != reflect.Ptr {
		return fmt.Errorf("dest must be a pointer")
	}

	if destValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("dest must be a pointer to a struct")
	}

	// destValue가 포인터이면서 구조체인 경우 포인터를 제거한다.
	destValue = destValue.Elem()

	err := m.mapByValue(destValue, srcValue)

	return err
}

func (m *structMapper) mapByValue(destValue reflect.Value, srcValue reflect.Value) error {

	for i := 0; i < destValue.NumField(); i++ {
		destField := destValue.Type().Field(i)
		tag := destField.Tag.Get(m.tagName)
		// tag가 없는경우 이름을 가져온다.
		if tag == "" {
			tag = destField.Name
		}

		srcFieldName := strings.Split(tag, ",")[0]
		if srcFieldName == "" {
			srcFieldName = destField.Name
		}

		srcFieldValue := srcValue.FieldByName(srcFieldName)
		if !srcFieldValue.IsValid() {
			continue
		}

		// base type이 다른 경우 다음 필드를 찾는다.
		// if destField.Type != srcFieldValue.Type() {
		// 	continue
		// }
		if destField.Type.Kind() != srcFieldValue.Type().Kind() {
			//continue
			return fmt.Errorf("expected type:  %v, but %v", destField.Type.Kind(), srcFieldValue.Type().Kind())
		}

		destFieldValue := destValue.FieldByName(destField.Name)
		if !destFieldValue.IsValid() {
			continue
		}

		if destField.Type.Kind() == reflect.Struct {
			if err := m.mapByValue(destFieldValue, srcFieldValue); err != nil {
				return err
			}
			continue
		}
		destFieldValue.Set(srcFieldValue)
	}

	return nil
}
