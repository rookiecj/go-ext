package encodingx

import (
	"bytes"
	"fmt"
	"io"
	"testing"
)

func TestBuffer_WriteAndRead(t *testing.T) {

	b := bytes.Buffer{}

	want := "Hello"
	b.WriteString(want)
	//b.WriteByte('9')
	fmt.Println(b.Len())

	got, err := b.ReadString(0)
	fmt.Println(b.Len())
	fmt.Println(err)
	fmt.Println(got)

	if want != got {
		t.Errorf("str want %s but %s", want, got)
	}

	if err != io.EOF {
		t.Errorf("err want %v but %v", io.EOF, err)
	}

	if b.Len() != 0 {
		t.Errorf("Len want %d but %d", 0, b.Len())
	}
}
