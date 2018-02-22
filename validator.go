package validator

import (
	"fmt"
	"reflect"
	"strings"
)

const (
	// TagName - will stay present in struct to apply the validators
	TagName = "struct-validator"
)

// relation between 'validator key type' and 'rules'
var nativeValidators = map[string][]string{
	"numeric": []string{"min", "max"},
	"string":  []string{"min", "max"},
}

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

// getValidatorKeyType - check the field type and returns the 'validator key type' associated to field type
func getValidatorKeyType(typeName string) string {
	if parts := strings.Split(typeName, "[]"); len(parts) == 2 {
		return nativeValidatorsKeyType["array"]
	}
	if nativeValidatorsKeyType[typeName] == "" {
		return "custom"
	}
	return nativeValidatorsKeyType[typeName]
}

// Will get the 'validator key type', get rules of field tag and get errors if they exist.
// A panic is throwed if the rule of 'messageInput' does not exists for the field 'validator key type'
func checkValidations(tags string, messageInput MessageInput) (errors []error) {
	// get validator key
	if messageInput.ValidatorKeyType = getValidatorKeyType(messageInput.FieldType.String()); messageInput.ValidatorKeyType == "" || len(tags) == 0 {
		return errors
	}
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
	return errors
}

// Will check if exists a native 'validator key type' and 'rule', and returns a error if exists
func checkIfExistsNativeValidadorKeyTypeAndRuleName(validatorKeyType string, ruleName string) error {
	//foreach validator
	for k, v := range nativeValidators {
		if validatorKeyType != k {
			continue
		}
		//foreach rule
		for _, r := range v {
			if ruleName == r {
				return fmt.Errorf("Error: The rule %s is a native rule of %s native type name, you cannot change this rule", ruleName, validatorKeyType)
			}
		}
	}
	return nil
}

// AddCustomValidator - Will add one custom validator, this method don't permit change native validators
// and returns a error when the 'typeName' and 'ruleName' are native validators
func AddCustomValidator(typeName string, ruleName string, handler func(MessageInput) error) error {
	if err := checkIfExistsNativeValidadorKeyTypeAndRuleName(typeName, ruleName); err != nil {
		return err
	}
	//if is a new typeName
	if types[typeName] == nil {
		types[typeName] = make(map[string](func(MessageInput) error))
	}
	types[typeName][ruleName] = handler
	return nil
}

// DelCustomValidator - Will remove one custom validator, this method don't permit change native validators
// and returns a error when the 'typeName' and 'ruleName' are native validators
func DelCustomValidator(typeName string, ruleName string) error {
	if err := checkIfExistsNativeValidadorKeyTypeAndRuleName(typeName, ruleName); err != nil {
		return err
	}
	delete(types[typeName], ruleName)
	if len(types[typeName]) == 0 {
		delete(types, typeName)
	}
	return nil
}
