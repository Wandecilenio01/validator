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
	CreateAt time.Time `json:"createAt" struct-validator:"after_or_equal_date:today+3"`
	Email    string    `json:"email" struct-validator:"email"`
	Site     string    `json:"site" struct-validator:"url"`
	IPv4     string    `json:"ipv4" struct-validator:"ipv4"`
	JSON     string    `json:"json" struct-validator:"json"`
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
	onePerson := MyModel{1, "Wan", 21, time.Now().AddDate(0, 0, +3), "myem@emai.coms", "regextestercom/20", "123.123.123.123", "[\"sddsfsd\"]"}
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
