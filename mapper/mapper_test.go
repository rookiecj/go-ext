package mapper

import (
	"log"
	"testing"
)

func TestMapper_Map(t *testing.T) {

	// panic: reflect: reflect.Value.Set using value obtained using unexported field [recovered]
	//type src struct {
	//	id   int
	//	name string
	//}
	//
	//type dest struct {
	//	ID   int    `structMapper:"id"`
	//	Name string `structMapper:"name"`
	//}

	//srcValue := src{
	//	id:   1,
	//	name: "name",
	//}

	type src struct {
		Seq   int
		Label string
		Addr  string
	}

	type dest struct {
		ID   int    `mapper:"Seq"`
		Name string `mapper:"Label"`
		Addr string
	}

	srcValue := src{
		Seq:   13,
		Label: "Prime",
		Addr:  "Number",
	}

	destValue := dest{}

	m := NewMapper("mapper")

	if err := m.Map(&destValue, srcValue); err != nil {
		t.Fatal(err)
	}
	log.Printf("srcValue: %v", srcValue)
	log.Printf("destValue: %v", destValue)

	if destValue.ID != srcValue.Seq {
		t.Errorf("destValue.ID is not equal to srcValue.Seq")
	}

	if destValue.Name != srcValue.Label {
		t.Errorf("destValue.Name is not equal to srcValue.Label")
	}
}

func TestMapper_Map_Embedded(t *testing.T) {

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

	m := NewMapper("mapper")

	if err := m.Map(&destValue, srcValue); err != nil {
		t.Fatal(err)
	}
	log.Printf("srcValue: %v", srcValue)
	log.Printf("destValue: %v", destValue)

	if destValue.ID != srcValue.Seq {
		t.Errorf("destValue.ID(%v) is not equal to srcValue.Seq(%v", destValue.ID, srcValue.Seq)
	}

	if destValue.Name != srcValue.Label {
		t.Errorf("destValue.Name(%s) is not equal to srcValue.Label(%s)", destValue.Name, srcValue.Label)
	}

	if destValue.Desc != srcValue.Desc {
		t.Errorf("destValue.Desc want %s got %s", srcValue.Desc, destValue.Desc)
	}
}
