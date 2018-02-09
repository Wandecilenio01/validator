package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// TAG - will stay presence in struct to apply the validators
	TAG_NAME       = "struct-validator"
	PARAM_MIN_NAME = "min"
	PARAM_MAX_NAME = "max"
)

// Validate - will validate all structs with the tag "struct-validator" that you pass by argument
func Validate(st interface{}) []error {
	var errors []error
	stValue := reflect.ValueOf(st)
	for i := 0; i < stValue.NumField(); i++ {
		field := stValue.Type().Field(i)
		errors = append(errors, checkValidations(field.Tag.Get(TAG_NAME), field.Name, stValue.Field(i).Interface())...)
	}
	return errors
}

// checkValidations - will make a parse in string of validations
func checkValidations(validators string, field string, value interface{}) (errors []error) {
	if paramsTest := strings.Split(validators, ","); len(paramsTest) > 0 {
		typeName := paramsTest[0]
		if types[typeName] == nil {
			errors = append(errors, fmt.Errorf("Validator %s not registered", typeName))
		} else {
			paramsTestMap := make(map[string]interface{})
			for _, paramTest := range paramsTest[1:] {
				if parts := strings.Split(paramTest, "="); len(parts) == 2 {
					paramsTestMap[parts[0]] = parts[1]
				} else {
					paramsTestMap[parts[0]] = true
				}
			}
			errors = append(errors, types[typeName](value, paramsTestMap)...)
		}
	}
	return errors
}
