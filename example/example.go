package main

import (
	"fmt"
	"time"
)

// MyModel - just to make one test
type MyModel struct {
	ID             int64          `json:"id" struct-validator:"min:3|max:20"`
	Name           string         `json:"name" struct-validator:"regex:^[0-9]*$|required"`
	Age            int64          `json:"age" struct-validator:"min:3|max:20"`
	CreateAt       time.Time      `json:"createAt" struct-validator:"after_or_equal_date:today+3"`
	Email          string         `json:"email" struct-validator:"email"`
	Site           string         `json:"site" struct-validator:"url"`
	IPv4           string         `json:"ipv4" struct-validator:"ipv4"`
	JSON           string         `json:"json" struct-validator:"json"`
	AlphaDashField string         `json:"alphaDash" struct-validator:"alpha_dash_space"`
	AlphaNumField  string         `json:"alphaNUm" struct-validator:"alpha_num_space"`
	MyIntArray     []int          `json:"MyIntArray" struct-validator:"required"`
	MyFloat32Array []float32      `json:"MyFloat32Array" struct-validator:"min:2|max:5|distinct"`
	MyAnotherModel MyAnotherModel `json:"MyAnotherModel" struct-validator:""`
	MyInt          int
	MyBool         bool
}

type MyAnotherModel struct {
	FieldOne string
	FieldTwo int
}

func main() {
	var nativeValidatorsKeyType = map[string]string{
		"int":       "numeric",
		"int64":     "numeric",
		"int32":     "numeric",
		"int16":     "numeric",
		"int8":      "numeric",
		"uint":      "numeric",
		"uint64":    "numeric",
		"uint32":    "numeric",
		"uint16":    "numeric",
		"uint8":     "numeric",
		"uintptr":   "numeric",
		"float32":   "numeric",
		"float64":   "numeric",
		"string":    "string",
		"time.Time": "timestamp",
		"array":     "array",
	}

	fmt.Println(nativeValidatorsKeyType)
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
	// onePerson := MyModel{2, "", 21, time.Now().AddDate(0, 0, +3), " ", "", "", "", "&&&**%%///\\\\%s", "&&&**%%///\\\\", nil, nil, *new(MyAnotherModel), 0, true}
	// // errors := validator.Validate(onePerson, messages)
	// errors := validator.Validate(onePerson, nil)
	// if len(errors) > 0 {
	// 	for eindex, err := range errors {
	// 		fmt.Println("Error Nº", eindex, "Error ->", err.Error())
	// 	}
	// }
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
