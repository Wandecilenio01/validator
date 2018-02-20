package main

import (
	"fmt"
	"time"

	"github.com/Wandecilenio01/validator"
)

// MyModel - just to make one test
type MyModel struct {
	ID       int64     `json:"id" struct-validator:"min:3|max:20"`
	Name     string    `json:"name"`
	Age      int64     `json:"age" struct-validator:"min:3|max:20"`
	CreateAt time.Time `json:"createAt" struct-validator:"before_date:today"`
}

func main() {
	// validator.AddCustomValidator("string", "name", func(messageInput validator.MessageInput) error {
	// 	if len(messageInput.FieldValue.(string)) < 4 {
	// 		return fmt.Errorf("Erro: O campo nome ...")
	// 	}
	// 	return nil
	// })
	// validator.DelCustomValidator("string", "name")
	// messages := map[string]map[string]string{
	// 	"all": map[string]string{
	// 		"min": "The min value for {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
	// 	},
	// 	"Age": map[string]string{
	// 		"max": "The max {{.fieldName}} should be {{.ruleValue}}, and not {{.value}} ",
	// 	},
	// }
	onePerson := MyModel{1, "Wan", 21, time.Now().AddDate(0, 0, 0)}
	errors := validator.Validate(onePerson, nil)
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
