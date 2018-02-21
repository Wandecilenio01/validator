package main

import (
	"fmt"
	"time"

	"github.com/Wandecilenio01/validator"
)

// MyModel - just to make one test
type MyModel struct {
	ID             int64     `json:"id" struct-validator:"min:3|max:20"`
	Name           string    `json:"name" struct-validator:"regex:^[0-9]*$|length:12"`
	Age            int64     `json:"age" struct-validator:"min:3|max:20"`
	CreateAt       time.Time `json:"createAt" struct-validator:"after_or_equal_date:today+3"`
	Email          string    `json:"email" struct-validator:"email"`
	Site           string    `json:"site" struct-validator:"url"`
	IPv4           string    `json:"ipv4" struct-validator:"ipv4"`
	JSON           string    `json:"json" struct-validator:"json"`
	AlphaDashField string    `json:"alphaDash" struct-validator:"alpha_dash_space"`
	AlphaNumField  string    `json:"alphaNUm" struct-validator:"alpha_num_space"`
	MyIntArray     []int
	MyInt64Array   []int64
	MyRune         []rune
	MyInt32Array   []int32
	MyInt16Array   []int16
	MyInt8Array    []int8
	MyUintArray    []uint
	MyUint64Array  []uint64
	MyUint32Array  []uint32
	MyUint16Array  []uint16
	MyUint8Array   []uint8
	MyInt          int
	MyBool         bool
}

func main() {
	// validator.AddCustomValidator("string", "name", func(messageInput validator.MessageInput) error {
	// 	if len(messageInput.FieldValue.(string)) < 4 {
	// 		return fmt.Errorf("Erro: O campo nome ...")
	// 	}
	// 	return nil
	// })
	// validator.DelCustomValidator("string", "name")
	messages := map[string]map[string]string{
		"all": map[string]string{
			"min": "The min value for {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
		},
		"Age": map[string]string{
			"max": "The max {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
		},
	}
	onePerson := MyModel{2, "123", 21, time.Now().AddDate(0, 0, +3), " ", "", "", "", "&&&**%%///\\\\%s", "&&&**%%///\\\\", []int{}, []int64{}, []rune{}, []int32{}, []int16{}, []int8{}, []uint{}, []uint64{}, []uint32{}, []uint16{}, []uint8{}, 0, true}
	errors := validator.Validate(onePerson, messages)
	if len(errors) > 0 {
		for eindex, err := range errors {
			fmt.Println("Error Nº", eindex, "Error ->", err.Error())
		}
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
