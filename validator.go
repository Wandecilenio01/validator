package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// TAG - will stay presence in struct to apply the validators
	TAG = "struct-validator"
)

// Validate - will validate all structs with the tag "struct-validator" that you pass by argument
func Validate(st interface{}) (errors []error) {
	stValue := reflect.ValueOf(st)
	for i := 0; i < stValue.NumField(); i++ {
		if tag := stValue.Type().Field(i).Tag.Get(TAG); tag == "" || tag == "-" {
			continue
		} else {
			e := checkValidations(tag, stValue.Type().Field(i).Name, stValue.Field(i).Interface())
			if len(e) > 0 {
				for _, a := range e {
					errors = append(errors, a)
				}
			}
		}
	}

	return errors
}

// checkValidations - will make a parse in string of validations
func checkValidations(validators string, field string, value interface{}) (errors []error) {
	for _, kindTest := range strings.Split(validators, ";") {
		testParams := strings.Split(kindTest, ",")
		test := testParams[0]
		p := testParams[1:len(testParams)]
		params := make(map[string]interface{})
		for _, p1 := range p {
			p2 := strings.Split(p1, "=")
			if len(p2) > 0 {
				params[p2[0]] = p2[1]
			}
		}

		nameFn := test
		if Types[nameFn] == nil {
			nameFn = "CUSTOM_" + nameFn
		}

		if Types[nameFn] == nil {
			errors = append(errors, fmt.Errorf("Validator %s not registered", test))
			continue
		}

		e := Types[nameFn](value, params)
		if len(e) > 0 {
			for _, a := range e {
				errors = append(errors, a)
			}
		}
	}
	return errors
}
