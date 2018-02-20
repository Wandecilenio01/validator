package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	// TimestampDefaultFormat -
	TimestampDefaultFormat = "2006-1-2 15:4:5"
	// TimestampDateDefaultFormat -
	TimestampDateDefaultFormat = "2006-1-2"
)

// MessageInput - Input struct used
type MessageInput struct {
	FieldName        string
	FieldType        reflect.Type
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
		if lessCompareNumbers(messageInput.FieldValue, messageInput.RuleValue, messageInput.FieldType.Kind()) {
			return nil
		}
		return generateErrorMessage(messageInput)
	}
	types["numeric"]["max"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		if moreCompareNumbers(messageInput.FieldValue, messageInput.RuleValue, messageInput.FieldType.Kind()) {
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
	//date INCOMPLETE
	// today
	// after, after_date, equal, equal_date, before, before_equal
	types["timestamp"] = make(map[string](func(MessageInput) error))
	types["timestamp"]["after"] = func(messageInput MessageInput) error {
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.After(getTimestampFromRuleString(messageInput.RuleValue)) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return generateErrorMessage(messageInput)
	}
	types["timestamp"]["after_date"] = func(messageInput MessageInput) error {
		fieldValueTime := truncateTime(messageInput.FieldValue.(time.Time))
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		if fieldValueTime.After(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return generateErrorMessage(messageInput)
	}
	types["timestamp"]["before"] = func(messageInput MessageInput) error {
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.Before(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return generateErrorMessage(messageInput)
	}
	types["timestamp"]["before_date"] = func(messageInput MessageInput) error {
		ruleValueTime := truncateTime(getTimestampFromRuleString(messageInput.RuleValue))
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.Before(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return generateErrorMessage(messageInput)
	}
	types["timestamp"]["equal_date"] = func(messageInput MessageInput) error {
		ruleValueTime := truncateTime(getTimestampFromRuleString(messageInput.RuleValue))
		fieldValueTime := truncateTime(messageInput.FieldValue.(time.Time))
		if fieldValueTime.Equal(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return generateErrorMessage(messageInput)
	}
	types["timestamp"]["equal"] = func(messageInput MessageInput) error {
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.Equal(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return generateErrorMessage(messageInput)
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

func getTimestampFromRuleString(value string) time.Time {
	parts := strings.Split(value, "+")
	if parts[0] != "today" {
		panic("Error: The rule value should be 'today' or 'today+1' ... ")
	}
	if len(parts) == 2 {
		days, err := strconv.Atoi(parts[1])
		if err == nil {
			return time.Now().AddDate(0, 0, days)
		}
		panic(err)
	}
	return time.Now()
}

func truncateTime(valueTime time.Time) time.Time {
	return new(time.Time).AddDate(valueTime.Year()-1, int(valueTime.Month())-1, valueTime.Day()-1)
}
