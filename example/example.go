package main

import (
	"fmt"

	"github.com/wandecilenio/validator"
)

// MyModel - just to make one test
type MyModel struct {
	ID   int64  `json:"id" struct-validator:"integer"`
	Name string `json:"name" struct-validator:"string,min=3,max=20"`
	Age  int64  `json:"age" struct-validator:"integer"`
}

func main() {
	fmt.Println("Hello there!!!")

	onePerson := MyModel{1, "Wandecilenio", 21}

	errors := validator.Validate(onePerson)
	if len(errors) > 0 {
		for eindex, err := range errors {
			fmt.Println("Error Nº", eindex, "Error ->", err.Error())
		}
	}
}
