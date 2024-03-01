package pb

import (
	"fmt"
	"google.golang.org/protobuf/proto"
	"io"
)

func ProtoBufBodyParser(buf io.Reader, bodyPtr any) (err error) {
	//bodyType := reflect.TypeOf(bodyPtr).Elem()
	message := bodyPtr.(proto.Message)

	var data []byte
	data, err = io.ReadAll(buf)
	if err != nil {
		return fmt.Errorf("error while reading: %s", err)
	}

	err = proto.Unmarshal(data, message)
	return
}
