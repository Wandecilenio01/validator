package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
)

type MessageInput struct {
	FieldName        string
	FieldKind        reflect.Kind
	ValidatorKeyType string
	FieldValue       interface{}
	RuleName         string
	RuleValue        string
	CustomMessages   map[string]map[string]string
}

// relation of type : rule, for example: int64 : min
var types map[string]map[string](func(MessageInput) error)

func init() {
	types = make(map[string]map[string](func(MessageInput) error))
	defineTypes()
}

//PanicOnEmptyRuleValue - Check if the rule value is empty and panic if true
func PanicOnEmptyRuleValue(rule string, ruleValue string) {
	if ruleValue == "" {
		panic(fmt.Sprintf("The rule %s cannot be empty, pass a value, like %v:value", rule, rule))
	}
}

func lessCompareNumbers(num1 interface{}, num2 string, kind reflect.Kind) bool {
	if intValue, err := getIntFromInterface(num1, kind); err == nil {
		return intValue < getIntFromString(num2)
	}
	if uintValue, err := getUintFromInterface(num1, kind); err == nil {
		return uintValue < getUintFromString(num2)
	}
	if floatValue, err := getFloatFromInterface(num1, kind); err == nil {
		return floatValue < getFloatFromString(num2)
	}
	panic("Error: Unknwon type")
}

func moreCompareNumbers(num1 interface{}, num2 string, kind reflect.Kind) bool {
	if intValue, err := getIntFromInterface(num1, kind); err == nil {
		return intValue > getIntFromString(num2)
	}
	if uintValue, err := getUintFromInterface(num1, kind); err == nil {
		return uintValue > getUintFromString(num2)
	}
	if floatValue, err := getFloatFromInterface(num1, kind); err == nil {
		return floatValue > getFloatFromString(num2)
	}
	panic("Error: Unknwon type")
}

func getFloatFromInterface(inter interface{}, kind reflect.Kind) (float64, error) {
	switch kind {
	case reflect.Float32:
		return float64(inter.(float32)), nil
	case reflect.Float64:
		return inter.(float64), nil
	}
	return 0, errors.New("Error: The interface passed is not a Float")
}

func getUintFromInterface(inter interface{}, kind reflect.Kind) (uint64, error) {
	switch kind {
	case reflect.Uint:
		return uint64(inter.(uint)), nil
	case reflect.Uint8:
		return uint64(inter.(uint8)), nil
	case reflect.Uint16:
		return uint64(inter.(uint16)), nil
	case reflect.Uint32:
		return uint64(inter.(uint32)), nil
	case reflect.Uint64:
		return inter.(uint64), nil
	}
	return 0, errors.New("Error: The interface passed is not a Uint")
}

func getIntFromInterface(inter interface{}, kind reflect.Kind) (int64, error) {
	switch kind {
	case reflect.Int:
		return int64(inter.(int)), nil
	case reflect.Int8:
		return int64(inter.(int8)), nil
	case reflect.Int16:
		return int64(inter.(int16)), nil
	case reflect.Int32:
		return int64(inter.(int32)), nil
	case reflect.Int64:
		return inter.(int64), nil
	}
	return 0, errors.New("Error: The interface passed is not a Int")
}

// Define native types
func defineTypes() {
	// numeric
	types["numeric"] = make(map[string](func(MessageInput) error))
	types["numeric"]["min"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("min", messageInput.RuleValue)
		if lessCompareNumbers(messageInput.FieldValue, messageInput.RuleValue, messageInput.FieldKind) {
			return nil
		}
		return generateErrorMessage(messageInput)
	}
	types["numeric"]["max"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		if moreCompareNumbers(messageInput.FieldValue, messageInput.RuleValue, messageInput.FieldKind) {
			return nil
		}
		return generateErrorMessage(messageInput)
	}
	// string
	types["string"] = make(map[string](func(MessageInput) error))
	types["string"]["max"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		if int64(len(messageInput.FieldValue.(string))) > getIntFromString(messageInput.RuleValue) {
			return generateErrorMessage(messageInput)
		}
		return nil
	}
}

func getFloatFromString(value string) float64 {
	if convertedValue, err := strconv.ParseFloat(value, 64); err == nil {
		return convertedValue
	}
	panic("Error: " + value + "is not a valid Float")
}

func getIntFromString(value string) int64 {
	if convertedValue, err := strconv.ParseInt(value, 10, 64); err == nil {
		return convertedValue
	}
	panic("Error: " + value + "is not a valid Int")
}

func getUintFromString(value string) uint64 {
	if convertedValue, err := strconv.ParseUint(value, 10, 64); err == nil {
		return convertedValue
	}
	panic("Error: " + value + "is not a valid Uint")
}
