package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// TagName - will stay presence in struct to apply the validators
	TagName = "struct-validator"
)

var nativeValidators = map[string][]string{
	"numeric": []string{"min", "max"},
	"string":  []string{"min", "max"},
}

// Validate - will validate all structs with the tag "struct-validator" that you pass by argument
func Validate(st interface{}, messages map[string]map[string]string) (errors []error) {
	stValue := reflect.ValueOf(st)
	// for each field
	for i := 0; i < stValue.NumField(); i++ {
		field := stValue.Type().Field(i)
		errors = append(errors, checkValidations(strings.Replace(field.Tag.Get(TagName), " ", "", -1), MessageInput{
			FieldName:      field.Name,
			FieldValue:     stValue.Field(i).Interface(),
			CustomMessages: messages,
			FieldType:      field.Type,
		})...)
	}
	return errors
}

func getValidatorKeyType(typeName string) string {
	if parts := strings.Split(typeName, "[]"); len(parts) == 2 {
		switch parts[1] {
		case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8", "uintptr", "float32", "float64":
			return "array_numeric"
		case "bool":
			return "array_boolean"
		case "string":
			return "array_string"
		case "time.Time":
			return "array_timestamp"
		default:
			return ""
		}
	} else {
		switch typeName {
		case "int", "int64", "int32", "int16", "int8", "uint", "uint64", "uint32", "uint16", "uint8", "uintptr", "float32", "float64":
			return "numeric"
		case "string":
			return "string"
		case "time.Time":
			return "timestamp"
		default:
			return ""
		}
	}
}

// Will make a parse in string of validations
func checkValidations(tags string, messageInput MessageInput) (errors []error) {
	// get validator key
	if messageInput.ValidatorKeyType = getValidatorKeyType(messageInput.FieldType.String()); messageInput.ValidatorKeyType != "" && len(tags) > 0 {
		// get rules from field
		rules := strings.Split(tags, "|")
		for _, rule := range rules {
			//if rule has value
			parts := strings.Split(strings.Replace(rule, " ", "", -1), ":")
			messageInput.RuleName = parts[0]
			if len(parts) == 2 {
				messageInput.RuleValue = parts[1]
			}
			//get errors
			if types[messageInput.ValidatorKeyType][messageInput.RuleName] == nil {
				panic(fmt.Sprintf("The rule '%s' does not exists in %s validator", messageInput.RuleName, messageInput.ValidatorKeyType))
			}
			if err := types[messageInput.ValidatorKeyType][messageInput.RuleName](messageInput); err != nil {
				errors = append(errors, err)
			}
		}
	}
	return errors
}

func checkIfExistsNativeRuleName(typeName string, ruleName string) error {
	//foreach validator
	for k, v := range nativeValidators {
		if typeName == k {
			//foreach rule
			for _, r := range v {
				if ruleName == r {
					return fmt.Errorf("Error: The rule %s is a native rule of %s native type name, you cannot change this rule", ruleName, typeName)
				}
			}
		}
	}
	return nil
}

// AddCustomValidator - Will add one custom validator
func AddCustomValidator(typeName string, ruleName string, handler func(MessageInput) error) error {
	if err := checkIfExistsNativeRuleName(typeName, ruleName); err != nil {
		return err
	}
	//if is a new typeName
	if types[typeName] == nil {
		types[typeName] = make(map[string](func(MessageInput) error))
	}
	types[typeName][ruleName] = handler
	return nil
}

// DelCustomValidator - Will remove one custom validator
func DelCustomValidator(typeName string, ruleName string) error {
	if err := checkIfExistsNativeRuleName(typeName, ruleName); err != nil {
		return err
	}
	delete(types[typeName], ruleName)
	if len(types[typeName]) == 0 {
		delete(types, typeName)
	}
	return nil
}
