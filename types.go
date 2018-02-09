package validator

import (
	"fmt"
)

// Types - all types registered to validate the struct
var types map[string](func(interface{}, map[string]interface{}) []error)

func init() {
	types = make(map[string](func(interface{}, map[string]interface{}) []error))
	defineTypes()
}

func defineTypes() {
	types["float32"] = func(value interface{}, params map[string]interface{}) (errors []error) {
		if min := params[PARAM_MIN_NAME]; min != nil && value.(float32) < min.(float32) {
			errors = append(errors, fmt.Errorf("The value cannot be less than %d", min))
		}
		if max := params[PARAM_MAX_NAME]; max != nil && value.(float32) > max.(float32) {
			errors = append(errors, fmt.Errorf("The value cannot be greater than %d", max))
		}
		return errors
	}
	types["float64"] = func(value interface{}, params map[string]interface{}) (errors []error) {
		if min := params[PARAM_MIN_NAME]; min != nil && value.(float64) < min.(float64) {
			errors = append(errors, fmt.Errorf("The value cannot be less than %d", min))
		}
		if max := params[PARAM_MAX_NAME]; max != nil && value.(float64) > max.(float64) {
			errors = append(errors, fmt.Errorf("The value cannot be greater than %d", max))
		}
		return errors
	}
}

// AddCustomValidator - will add one custom validator
func AddCustomValidator(name string, validateFn func(interface{}, map[string]interface{}) []error) {
	types[name] = validateFn
}

// DelCustomValidator - will remove one custom validator
func DelCustomValidator(name string) {
	delete(types, name)
}
