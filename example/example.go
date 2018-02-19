package main

import (
	"fmt"

	"github.com/Wandecilenio01/validator"
)

// MyModel - just to make one test
type MyModel struct {
	ID   int64  `json:"id" struct-validator:"min:3|max:20"`
	Name string `json:"name" struct-validator:"name"`
	Age  int64  `json:"age" struct-validator:"min:3|max:20"`
}

func main() {
	validator.AddCustomValidator("string", "name", func(messageInput validator.MessageInput) error {
		if len(messageInput.FieldValue.(string)) < 4 {
			return fmt.Errorf("Erro: O campo nome ...")
		}
		return nil
	})
	validator.DelCustomValidator("string", "name")
	messages := map[string]map[string]string{
		"all": map[string]string{
			"min": "The min value for {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
		},
		"Age": map[string]string{
			"max": "The max {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
		},
	}
	onePerson := MyModel{1, "Wan", 21}
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
