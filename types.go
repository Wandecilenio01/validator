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
	// TimestampDefaultFormat - Timestamp format used as output error messages
	TimestampDefaultFormat = "2006-1-2 15:4:5"
	// TimestampDateDefaultFormat - Date format used as output error messages
	TimestampDateDefaultFormat = "2006-1-2"
	// EmailRegex - Email regular expression used for validade email string
	EmailRegex = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	// URLRegex - URL regular expression used for validade url string
	URLRegex = "^((http[s]?|ftp):\\/)?\\/?([^:\\/\\s]+)((\\/\\w+)*\\/)([\\w\\-\\.]+[^#?\\s]+)(.*)?(#[\\w\\-]+)?$"
	// IPv4Regex - IPv4 regular expression used for validade ipv4 string
	IPv4Regex = "((25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})[.](25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})[.](25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2})[.](25[0-5]|2[0-4][0-9]|[0-1]?[0-9]{1,2}))"
	// IPv6Regex -
	// IPv6Regex = ""
	// AlphabeticRegex - Latin regular expression used for validade alphabetic string without spaces
	AlphabeticRegex = "^[A-Za-z\u00C0-\u00D6\u00D8-\u00f6\u00f8-\u00ff]*$"
	// AlphabeticSpacesRegex - Latin regular expression used for validade alphabetic string with spaces
	AlphabeticSpacesRegex = "^[A-Za-z\u00C0-\u00D6\u00D8-\u00f6\u00f8-\u00ff\\s]*$"
	// AlphaNumericDashRegex - Latin regular expression used for validade string with alphabet, numbers,
	// slash and underscore char's without spaces
	AlphaNumericDashRegex = "^[a-zA-Z\u00C0-\u00D6\u00D8-\u00f6\u00f8-\u00ff0-9-_]*$"
	// AlphaNumericDashSpacesRegex - Latin regular expression used for validade string with alphabet, numbers,
	// slash and underscore char's with spaces
	AlphaNumericDashSpacesRegex = "^[a-zA-Z\u00C0-\u00D6\u00D8-\u00f6\u00f8-\u00ff0-9-_\\s]*$"
	// AlphaNumericRegex - Latin regular expression used for validade string with alphabet, numbers
	// char's without spaces
	AlphaNumericRegex = "^[a-zA-Z0-9]*$"
	// AlphaNumericSpacesRegex - Latin regular expression used for validade string with alphabet, numbers
	// char's with spaces
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

// relation between 'validator key type' and 'rule' and 'handler'
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
		uintFieldValue, errFieldValue := GetUintFromString(fieldValueString)
		uintRuleValue, errRuleValue := GetUintFromString(messageInput.RuleValue)
		if errFieldValue != nil {
			//try float
			floatFieldValue, errFieldValue := GetFloatFromString(fieldValueString)
			floatRuleValue, errRuleValue := GetFloatFromString(messageInput.RuleValue)
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
		uintFieldValue, errFieldValue := GetUintFromString(fieldValueString)
		uintRuleValue, errRuleValue := GetUintFromString(messageInput.RuleValue)
		if errFieldValue != nil {
			//try float
			floatFieldValue, errFieldValue := GetFloatFromString(fieldValueString)
			floatRuleValue, errRuleValue := GetFloatFromString(messageInput.RuleValue)
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
		if uint64(len(messageInput.FieldValue.(string))) <= GetUintRuleValueOrPanic(messageInput.RuleName, messageInput.RuleValue) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["min"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("min", messageInput.RuleValue)
		if uint64(len(messageInput.FieldValue.(string))) >= GetUintRuleValueOrPanic(messageInput.RuleName, messageInput.RuleValue) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["length"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		if uint64(len(messageInput.FieldValue.(string))) == GetUintRuleValueOrPanic(messageInput.RuleName, messageInput.RuleValue) {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["email"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, EmailRegex)
	}
	types["string"]["url"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, URLRegex)
	}
	types["string"]["ipv4"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, IPv4Regex)
	}
	// types["string"]["ipv6"] = func(messageInput MessageInput) error {
	// return MatchRegex(messageInput, IPv6Regex)
	// }
	types["string"]["json"] = func(messageInput MessageInput) error {
		var temp interface{}
		if fieldValueString := messageInput.FieldValue.(string); fieldValueString == "" || json.Unmarshal([]byte(messageInput.FieldValue.(string)), &temp) == nil {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["string"]["alpha"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, AlphabeticRegex)
	}
	types["string"]["alpha_dash"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, AlphaNumericDashRegex)
	}
	types["string"]["alpha_num"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, AlphaNumericRegex)
	}
	types["string"]["alpha_space"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, AlphabeticSpacesRegex)
	}
	types["string"]["alpha_dash_space"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, AlphaNumericDashSpacesRegex)
	}
	types["string"]["alpha_num_space"] = func(messageInput MessageInput) error {
		return MatchRegex(messageInput, AlphaNumericSpacesRegex)
	}
	types["string"]["regex"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("regex", messageInput.RuleValue)
		return MatchRegex(messageInput, messageInput.RuleValue)
	}
	types["string"]["required"] = func(messageInput MessageInput) error {
		messageInput.RuleName = "min"
		messageInput.RuleValue = "1"
		return types["string"]["min"](messageInput)
	}
	//timestamps
	types["timestamp"] = make(map[string](func(MessageInput) error))
	types["timestamp"]["after"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("after", messageInput.RuleValue)
		ruleValueTime := GetTimestampFromRuleString(messageInput.RuleValue)
		fieldValueTime := messageInput.FieldValue.(time.Time)
		if fieldValueTime.After(GetTimestampFromRuleString(messageInput.RuleValue)) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["after_date"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("after_date", messageInput.RuleValue)
		fieldValueTime := TruncateTime(messageInput.FieldValue.(time.Time))
		ruleValueTime := GetTimestampFromRuleString(messageInput.RuleValue)
		if fieldValueTime.After(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["before"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("before", messageInput.RuleValue)
		ruleValueTime := GetTimestampFromRuleString(messageInput.RuleValue)
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
		ruleValueTime := TruncateTime(GetTimestampFromRuleString(messageInput.RuleValue))
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
		ruleValueTime := TruncateTime(GetTimestampFromRuleString(messageInput.RuleValue))
		fieldValueTime := TruncateTime(messageInput.FieldValue.(time.Time))
		if fieldValueTime.Equal(ruleValueTime) {
			return nil
		}
		messageInput.RuleValue = ruleValueTime.Format(TimestampDateDefaultFormat)
		messageInput.FieldValue = fieldValueTime.Format(TimestampDateDefaultFormat)
		return GenerateErrorMessage(messageInput)
	}
	types["timestamp"]["equal"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("equal", messageInput.RuleValue)
		ruleValueTime := GetTimestampFromRuleString(messageInput.RuleValue)
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
	types["array"]["max"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("max", messageInput.RuleValue)
		uintRuleValue := GetUintRuleValueOrPanic(messageInput.RuleName, messageInput.RuleValue)
		if interfaceArrayFieldValue, errFieldValue := GetInterfaceArrayFromInterface(messageInput.FieldValue); errFieldValue != nil {
			return errFieldValue
		} else if uint64(len(interfaceArrayFieldValue)) <= uintRuleValue {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["array"]["min"] = func(messageInput MessageInput) error {
		PanicOnEmptyRuleValue("min", messageInput.RuleValue)
		uintRuleValue := GetUintRuleValueOrPanic(messageInput.RuleName, messageInput.RuleValue)
		if interfaceArrayFieldValue, errFieldValue := GetInterfaceArrayFromInterface(messageInput.FieldValue); errFieldValue != nil {
			return errFieldValue
		} else if uint64(len(interfaceArrayFieldValue)) >= uintRuleValue {
			return nil
		}
		return GenerateErrorMessage(messageInput)
	}
	types["array"]["distinct"] = func(messageInput MessageInput) error {
		interfaceArrayFieldValue, errFieldValue := GetInterfaceArrayFromInterface(messageInput.FieldValue)
		if errFieldValue != nil {
			return errFieldValue
		}
		arrayLen := len(interfaceArrayFieldValue)
		for i := 0; i < arrayLen-1; i++ {
			for j := i + 1; j < arrayLen; j++ {
				if interfaceArrayFieldValue[i] == interfaceArrayFieldValue[j] {
					return GenerateErrorMessage(messageInput)
				}
			}
		}
		return nil
	}
	types["array"]["required"] = func(messageInput MessageInput) error {
		messageInput.RuleName = "min"
		messageInput.RuleValue = "1"
		return types["array"]["min"](messageInput)
	}
}

// MatchRegex - Check if a string regex match the messageInput.FieldValue, if not match then an error is
// returned, and return nil if not
func MatchRegex(messageInput MessageInput, regex string) error {
	if fieldValueString := messageInput.FieldValue.(string); fieldValueString == "" || regexp.MustCompile(regex).MatchString(fieldValueString) {
		return nil
	}
	return GenerateErrorMessage(messageInput)
}

// PanicOnEmptyRuleValue - Check if the rule value is empty and panic if true
func PanicOnEmptyRuleValue(rule string, ruleValue string) {
	if ruleValue == "" {
		panic(fmt.Sprintf("The rule %s cannot be empty, pass a value, like %v:value", rule, rule))
	}
}

// GetUintRuleValueOrPanic - Try to get an uint from ruleValue string, if possible then return uint64 value,
// and panic if not
func GetUintRuleValueOrPanic(rule string, ruleValue string) uint64 {
	if uintRuleValue, errRuleValue := GetUintFromString(ruleValue); errRuleValue == nil {
		return uintRuleValue
	}
	panic(fmt.Sprintf("The rule %s should to be a uint, pass a value, like %v:23", rule, rule))
}

// GetInterfaceArrayFromInterface - Try to get an array from a interface, if possible, then an interface array
// is returned, and if not an error is returned
func GetInterfaceArrayFromInterface(inter interface{}) ([]interface{}, error) {
	var interfaceArray []interface{}
	if fieldValueJSONString, err := json.Marshal(inter); err != nil {
		//marsh to get json string
		return nil, err
	} else if err = json.Unmarshal([]byte(fieldValueJSONString), &interfaceArray); err != nil {
		//unmarsh to get []interface{}
		return nil, err
	}
	return interfaceArray, nil
}

// GetFloatArrayFromInterface - Try to get a float64 array from a interface, if possible, then a float64  array
// is returned, and if not an error is returned
func GetFloatArrayFromInterface(inter interface{}) ([]float64, error) {
	var float64Array []float64
	if fieldValueJSONString, err := json.Marshal(inter); err != nil {
		//marsh to get json string
		return nil, err
	} else if err = json.Unmarshal([]byte(fieldValueJSONString), &float64Array); err != nil {
		//unmarsh to get []float64
		return nil, err
	}
	return float64Array, nil
}

// GetStringArrayFromInterface - Try to get a string array from a interface, if possible, then a string array
// is returned, and if not an error is returned
func GetStringArrayFromInterface(inter interface{}) ([]string, error) {
	var stringArray []string
	if fieldValueJSONString, err := json.Marshal(inter); err != nil {
		//marsh to get json string
		return nil, err
	} else if err = json.Unmarshal([]byte(fieldValueJSONString), &stringArray); err != nil {
		//unmarsh to get []uint64
		return nil, err
	}
	return stringArray, nil
}

// GetUintArrayFromInterface - Try to get an uint64 array from a interface, if possible, then a uint64 array
// is returned, and if not an error is returned
func GetUintArrayFromInterface(inter interface{}) ([]uint64, error) {
	var uint64Array []uint64
	if fieldValueJSONString, err := json.Marshal(inter); err != nil {
		//marsh to get json string
		return nil, err
	} else if err = json.Unmarshal([]byte(fieldValueJSONString), &uint64Array); err != nil {
		//unmarsh to get []uint64
		return nil, err
	}
	return uint64Array, nil
}

// GetFloatFromString - Try to get an float from a interface, if possible, then a float
// is returned, and if not an error is returned
func GetFloatFromString(value string) (float64, error) {
	if convertedValue, err := strconv.ParseFloat(value, 64); err == nil {
		return convertedValue, nil
	}
	return 0, fmt.Errorf("Error: %v is not a valid float", value)
}

// GetUintFromString - Try to get an uint from a interface, if possible, then a uint
// is returned, and if not an error is returned
func GetUintFromString(value string) (uint64, error) {
	if convertedValue, err := strconv.ParseUint(value, 10, 64); err == nil {
		return convertedValue, nil
	}
	return 0, fmt.Errorf("Error: %v is not a valid uint", value)
}

// GetTimestampFromRuleString - Try to get a time.Time from a rule value string, if possible, then a time.Time
// is returned, and if not a panic is returned
func GetTimestampFromRuleString(value string) time.Time {
	parts := strings.Split(value, "+")
	if parts[0] != "today" {
		panic("Error: The rule value should to be today or today+1 or today+2 or ... ")
	} else if len(parts) == 2 {
		days, err := strconv.Atoi(parts[1])
		if err == nil {
			return time.Now().AddDate(0, 0, days)
		}
		panic(err)
	}
	return time.Now()
}

// TruncateTime - Return a time.Time with the same day, month and year of valueTime parameter, but
// with hour = 0, minutes = 0 and seconds = 0
func TruncateTime(valueTime time.Time) time.Time {
	return new(time.Time).AddDate(valueTime.Year()-1, int(valueTime.Month())-1, valueTime.Day()-1)
}
