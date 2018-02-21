package validator

import (
	"encoding/json"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	// TimestampDefaultFormat -
	TimestampDefaultFormat = "2006-1-2 15:4:5"
	// TimestampDateDefaultFormat -
	TimestampDateDefaultFormat = "2006-1-2"
	// EmailRegex -
	EmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	// URLRegex -
	URLRegex = "^((http[s]?|ftp):\\/)?\\/?([^:\\/\\s]+)((\\/\\w+)*\\/)([\\w\\-\\.]+[^#?\\s]+)(.*)?(#[\\w\\-]+)?$"
	// IPv4Regex -
	IPv4Regex = "((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})[.](25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})[.](25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})[.](25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2}))"
	// IPv6Regex -
	// IPv6Regex = ""
	// AlphabeticRegex - Latin Regex
	AlphabeticRegex = "^[A-Za-z\u00C0-\u00D6\u00D8-\u00f6\u00f8-\u00ff]*$"
	// AlphabeticSpacesRegex - Latin Regex
	AlphabeticSpacesRegex = "^[A-Za-z\u00C0-\u00D6\u00D8-\u00f6\u00f8-\u00ff\\s]*$"
	// AlphaNumericDashRegex -
	AlphaNumericDashRegex = "^[a-zA-Z0-9-_]*$"
	// AlphaNumericDashSpacesRegex -
	AlphaNumericDashSpacesRegex = "^[a-zA-Z0-9-_\\s]*$"
	// AlphaNumericRegex -
	AlphaNumericRegex = "^[a-zA-Z0-9]*$"
	// AlphaNumericSpacesRegex -
	AlphaNumericSpacesRegex = "^[a-zA-Z0-9\\s]*$"
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

// Define native types
func defineTypes() {
	// numeric
	types["numeric"] = make(map[string](func(MessageInput) error))
	types["numeric"]["min"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("min", messageInput.RuleValue)
		fieldValueString := fmt.Sprintf("%v", messageInput.FieldValue)
		//try uint
		uintFieldValue, errFieldValue := getUintFromString(fieldValueString)
		uintRuleValue, errRuleValue := getUintFromString(messageInput.RuleValue)
		if errFieldValue != nil {
			//try float
			floatFieldValue, errFieldValue := getFloatFromString(fieldValueString)
			floatRuleValue, errRuleValue := getFloatFromString(messageInput.RuleValue)
			if errFieldValue != nil {
				return errFieldValue
			} else if errRuleValue != nil {
				panic(errRuleValue.Error())
			} else if floatFieldValue >= floatRuleValue {
				return nil
			}
		} else if errRuleValue != nil {
			panic(errRuleValue.Error())
		} else if uintFieldValue >= uintRuleValue {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["numeric"]["max"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		fieldValueString := fmt.Sprintf("%v", messageInput.FieldValue)
		//try uint
		uintFieldValue, errFieldValue := getUintFromString(fieldValueString)
		uintRuleValue, errRuleValue := getUintFromString(messageInput.RuleValue)
		if errFieldValue != nil {
			//try float
			floatFieldValue, errFieldValue := getFloatFromString(fieldValueString)
			floatRuleValue, errRuleValue := getFloatFromString(messageInput.RuleValue)
			if errFieldValue != nil {
				return errFieldValue
			} else if errRuleValue != nil {
				panic(errRuleValue.Error())
			} else if floatFieldValue <= floatRuleValue {
				return nil
			}
		} else if errRuleValue != nil {
			panic(errRuleValue.Error())
		} else if uintFieldValue <= uintRuleValue {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	// string
	types["string"] = make(map[string](func(MessageInput) error))
	types["string"]["max"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		if floatRuleValue, errRuleValue := getFloatFromString(messageInput.RuleValue); errRuleValue != nil {
			panic(errRuleValue.Error())
		} else if float64(len(messageInput.FieldValue.(string))) > floatRuleValue {
			return GenerateErrorMessage(messageInput)
		}
		return nil
	}
	types["string"]["min"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("min", messageInput.RuleValue)
		if floatRuleValue, errRuleValue := getFloatFromString(messageInput.RuleValue); errRuleValue != nil {
			panic(errRuleValue.Error())
		} else if float64(len(messageInput.FieldValue.(string))) < floatRuleValue {
			return GenerateErrorMessage(messageInput)
		}
		return nil
	}
	types["string"]["length"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("length", messageInput.RuleValue)
		if floatRuleValue, errRuleValue := getFloatFromString(messageInput.RuleValue); errRuleValue != nil {
			panic(errRuleValue.Error())
		} else if float64(len(messageInput.FieldValue.(string))) != floatRuleValue {
			return GenerateErrorMessage(messageInput)
		}
		return nil
	}
	types["string"]["email"] = func(messageInput MessageInput) error {
		if fieldValueString := messageInput.FieldValue.(string); fieldValueString == "" || regexp.MustCompile(EmailRegex).MatchString(fieldValueString) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["url"] = func(messageInput MessageInput) error {
		if fieldValueString := messageInput.FieldValue.(string); fieldValueString == "" || regexp.MustCompile(URLRegex).MatchString(fieldValueString) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["ipv4"] = func(messageInput MessageInput) error {
		if fieldValueString := messageInput.FieldValue.(string); fieldValueString == "" || regexp.MustCompile(IPv4Regex).MatchString(fieldValueString) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	// types["string"]["ipv6"] = func(messageInput MessageInput) error {
	// 	if regexp.MustCompile(IPv6Regex).MatchString(messageInput.FieldValue.(string)) {
	// 		return nil
	// 	}
	// 	return generateErrorMessage(messageInput)
	// }
	types["string"]["json"] = func(messageInput MessageInput) error {
		var temp interface{}
		if fieldValueString := messageInput.FieldValue.(string); fieldValueString == "" || json.Unmarshal([]byte(messageInput.FieldValue.(string)), &temp) == nil {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha"] = func(messageInput MessageInput) error {
		if regexp.MustCompile(AlphabeticRegex).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha_dash"] = func(messageInput MessageInput) error {
		if regexp.MustCompile(AlphaNumericDashRegex).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha_num"] = func(messageInput MessageInput) error {
		if regexp.MustCompile(AlphaNumericRegex).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha_space"] = func(messageInput MessageInput) error {
		if regexp.MustCompile(AlphabeticSpacesRegex).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha_dash_space"] = func(messageInput MessageInput) error {
		if regexp.MustCompile(AlphaNumericDashSpacesRegex).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha_num_space"] = func(messageInput MessageInput) error {
		if regexp.MustCompile(AlphaNumericSpacesRegex).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["regex"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("regex", messageInput.RuleValue)
		if regexp.MustCompile(messageInput.RuleValue).MatchString(messageInput.FieldValue.(string)) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	//timestamps
	types["timestamp"] = make(map[string](func(MessageInput) error))
	types["timestamp"]["after"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("after", messageInput.RuleValue)
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.After(getTimestampFromRuleString(messageInput.RuleValue)) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["after_date"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("after_date", messageInput.RuleValue)
		fieldValueTime := truncateTime(messageInput.FieldValue.(time.Time))
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		if fieldValueTime.After(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["before"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("before", messageInput.RuleValue)
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.Before(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["before_date"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("before_date", messageInput.RuleValue)
		ruleValueTime := truncateTime(getTimestampFromRuleString(messageInput.RuleValue))
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.Before(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["equal_date"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("equal_date", messageInput.RuleValue)
		ruleValueTime := truncateTime(getTimestampFromRuleString(messageInput.RuleValue))
		fieldValueTime := truncateTime(messageInput.FieldValue.(time.Time))
		if fieldValueTime.Equal(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["equal"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("equal", messageInput.RuleValue)
		ruleValueTime := getTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.Equal(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["after_or_equal"] = func(messageInput MessageInput) error {
		err := types["timestamp"]["after"](messageInput)
		if err != nil {
			return types["timestamp"]["equal"](messageInput)
		}
		return nil
	}
	types["timestamp"]["before_or_equal"] = func(messageInput MessageInput) error {
		err := types["timestamp"]["before"](messageInput)
		if err != nil {
			return types["timestamp"]["equal"](messageInput)
		}
		return nil
	}
	types["timestamp"]["after_or_equal_date"] = func(messageInput MessageInput) error {
		err := types["timestamp"]["after_date"](messageInput)
		if err != nil {
			return types["timestamp"]["equal_date"](messageInput)
		}
		return nil
	}
	types["timestamp"]["before_or_equal_date"] = func(messageInput MessageInput) error {
		err := types["timestamp"]["before_date"](messageInput)
		if err != nil {
			return types["timestamp"]["equal_date"](messageInput)
		}
		return nil
	}
	//arrays
	types["array"] = make(map[string](func(MessageInput) error))
	// types["array"]["max"] = func(messageInput MessageInput) error {
	// 	PanicOnEmptyRuleValue("max", messageInput.RuleValue)
	// 	//try uint
	// 	uintFieldValue, errFieldValue := getFloatArrayFromInterface(message.FieldValue)
	// 	uintRuleValue, errRuleValue := getUintFromString(messageInput.RuleValue)
	// 	if errFieldValue != nil {
	// 		//try float
	// 		floatFieldValue, errFieldValue := getFloatFromString(fieldValueString)
	// 		floatRuleValue, errRuleValue := getFloatFromString(messageInput.RuleValue)
	// 		if errFieldValue != nil {
	// 			return errFieldValue
	// 		} else if errRuleValue != nil {
	// 			panic(errRuleValue.Error())
	// 		} else if floatFieldValue <= floatRuleValue {
	// 			return nil
	// 		}
	// 	} else if errRuleValue != nil {
	// 		panic(errRuleValue.Error())
	// 	} else if uintFieldValue <= uintRuleValue {
	// 		return nil
	// 	}
	// 	return GenerateErrorMessage(messageInput)
	// }
}

//PanicOnEmptyRuleValue - Check if the rule value is empty and panic if true
func PanicOnEmptyRuleValue(rule string, ruleValue string) {
	if ruleValue == "" {
		panic(fmt.Sprintf("The rule %s cannot be empty, pass a value, like %v:value", rule, rule))
	}
}

func getFloatArrayFromInterface(inter interface{}) ([]float64, error) {
	var float64Array []float64
	//marsh to get json string
	fieldValueJSONString, err := json.Marshal(inter)
	if err != nil {
		return nil, err
	}
	//unmarsh to get []float64
	if err = json.Unmarshal([]byte(fieldValueJSONString), &float64Array); err != nil {
		return nil, err
	}
	return float64Array, nil
}

func getStringArrayFromInterface(inter interface{}) ([]string, error) {
	var stringArray []string
	//marsh to get json string
	fieldValueJSONString, err := json.Marshal(inter)
	if err != nil {
		return nil, err
	}
	//unmarsh to get []uint64
	if err = json.Unmarshal([]byte(fieldValueJSONString), &stringArray); err != nil {
		return nil, err
	}
	return stringArray, nil
}

func getUintArrayFromInterface(inter interface{}) ([]uint64, error) {
	var uint64Array []uint64
	//marsh to get json string
	fieldValueJSONString, err := json.Marshal(inter)
	if err != nil {
		return nil, err
	}
	//unmarsh to get []uint64
	if err = json.Unmarshal([]byte(fieldValueJSONString), &uint64Array); err != nil {
		return nil, err
	}
	return uint64Array, nil
}

func getFloatFromString(value string) (float64, error) {
	if convertedValue, err := strconv.ParseFloat(value, 64); err == nil {
		return convertedValue, nil
	}
	return 0, fmt.Errorf("Error: %v is not a valid float", value)
}

func getUintFromString(value string) (uint64, error) {
	if convertedValue, err := strconv.ParseUint(value, 10, 64); err == nil {
		return convertedValue, nil
	}
	return 0, fmt.Errorf("Error: %v is not a valid uint", value)
}

func getTimestampFromRuleString(value string) time.Time {
	parts := strings.Split(value, "+")
	if parts[0] != "today" {
		panic("Error: The rule value should to be today or today+1 or today+2 or ... ")
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
