package main

import (
	"log"

	"github.com/rookiecj/go-langext/mapper"
)

func main() {
	type common struct {
		Desc string
	}

	type src struct {
		Seq   int
		Label string
		Addr  string
		common
	}

	type dest struct {
		ID   int    `mapper:"Seq"`
		Name string `mapper:"Label"`
		Addr string
		common
	}

	srcValue := src{
		Seq:   13,
		Label: "Prime",
		Addr:  "Number",
		common: common{
			Desc: "Embedded",
		},
	}

	destValue := dest{}

	m := mapper.NewMapperWithTag("mapper")

	if err := m.Map(&destValue, srcValue); err != nil {
		panic(err)
	}
	log.Printf("srcValue: %v", srcValue)
	log.Printf("destValue: %v", destValue)
}
