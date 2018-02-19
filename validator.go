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
			FieldKind:      field.Type.Kind(),
		})...)
	}
	return errors
}

// Will make a parse in string of validations
func checkValidations(tags string, messageInput MessageInput) (errors []error) {
	//get key type
	//if is a number
	if fieldKindType := messageInput.FieldKind; fieldKindType >= reflect.Int && fieldKindType <= reflect.Float64 {
		messageInput.ValidatorKeyType = "numeric"
	} else if fieldKindType == reflect.String {
		messageInput.ValidatorKeyType = "string"
	} else {
		return errors
	}
	if types[messageInput.ValidatorKeyType] == nil {
		panic(fmt.Sprintf("The validator '%s' does not exists", messageInput.ValidatorKeyType))
	}
	// get rules from field
	if len(tags) > 0 {
		rules := strings.Split(tags, "|")
		for _, rule := range rules {
			//if rule has value
			parts := strings.Split(rule, ":")
			messageInput.RuleName = parts[0]
			if len(parts) == 2 {
				messageInput.RuleValue = parts[1]
			} else {
				messageInput.RuleValue = "1"
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
