package validator

import (
	"fmt"
	"strconv"
)

// Types - all types registered to validate the struct
var Types map[string](func(interface{}, map[string]interface{}) []error)

func init() {
	Types = make(map[string](func(interface{}, map[string]interface{}) []error))
	Types["numeric_float32"] = func(value interface{}, params map[string]interface{}) (errors []error) {
		if params["min"] != nil {
			min, _ := strconv.Atoi(params["min"].(string))
			if value.(float32) < float32(min) {
				errors = append(errors, fmt.Errorf("The value cannot be less than %d", min))
			}
		}

		if params["max"] != nil {
			max, _ := strconv.Atoi(params["max"].(string))
			if value.(float32) > float32(max) {
				errors = append(errors, fmt.Errorf("The value cannot be greater than %d", max))
			}
		}
		return errors
	}

	Types["numeric_float64"] = func(value interface{}, params map[string]interface{}) (errors []error) {
		if params["min"] != nil {
			min, _ := strconv.Atoi(params["min"].(string))
			if value.(float64) < float64(min) {
				errors = append(errors, fmt.Errorf("The value cannot be less than %d", min))
			}
		}

		if params["max"] != nil {
			max, _ := strconv.Atoi(params["max"].(string))
			if value.(float64) > float64(max) {
				errors = append(errors, fmt.Errorf("The value cannot be greater than %d", max))
			}
		}
		return errors
	}
}

// AddCustomValidator - will add one custom validator
func AddCustomValidator(name string, validateFn func(interface{}, map[string]interface{}) []error) {
	Types["CUSTOM_"+name] = validateFn
}

// DelCustomValidator - will remove one custom validator
func DelCustomValidator(name string) {
	delete(Types, "CUSTOM_"+name)
}
