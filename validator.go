package validator

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var (
	// TagName - will stay present in struct to apply the validators
	TagName string
)

var (
	// relation between 'validator key type' and 'rules'
	nativeValidators map[string][]string
	// relation between golang type names and 'validators key types'
	nativeValidatorsKeyType map[string]string
)

func init() {
	SetTag("struct-validator")
	nativeValidatorsKeyType = map[string]string{
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
	// fill nativeValidator using the 'type' relation
	nativeValidators = make(map[string][]string, 0)
	for validatorKeyType, ruleHandler := range types {
		nativeValidators[validatorKeyType] = make([]string, 0)
		for rule := range ruleHandler {
			nativeValidators[validatorKeyType] = append(nativeValidators[validatorKeyType], rule)
		}
	}
}

// Validate - will validate all structs with the tag "struct-validator" that you pass by argument
func Validate(st interface{}, messages map[string]map[string]string) (returnedErrors []error) {
	if st == nil {
		return append(returnedErrors, errors.New("The interface passed is nil"))
	}
	flagTag := true
	stValue := reflect.ValueOf(st)
	// mount message input list
	messagesInput := make([]MessageInput, 0)
	for i := 0; i < stValue.NumField(); i++ {
		field := stValue.Field(i)
		var interfaceValue interface{}
		if fieldKind := field.Type().Kind(); (reflect.Int <= fieldKind && fieldKind <= reflect.Int64) || fieldKind == reflect.Float32 || fieldKind == reflect.Float64 {
			if fieldKind == reflect.Float32 || fieldKind == reflect.Float64 {
				interfaceValue = field.Float()
			} else {
				//convert int type to float64
				interfaceValue = float64(field.Int())
			}
		} else if reflect.Uint <= fieldKind && fieldKind <= reflect.Uintptr {
			//convert uint type to uint64
			interfaceValue = field.Uint()
		} else {
			//anothers types
			interfaceValue = field.Interface()
		}
		messagesInput = append(messagesInput, MessageInput{
			FieldName:        stValue.Type().Field(i).Name,
			FieldValue:       interfaceValue,
			CustomMessages:   messages,
			FieldType:        field.Type(),
			ValidatorKeyType: getValidatorKeyType(field.Type().String()),
		})
	}
	//get errors
	for i := 0; i < stValue.NumField(); i++ {
		_, tagLookup := stValue.Type().Field(i).Tag.Lookup(TagName)
		if tagLookup {
			flagTag = false
		}
		messagesInput[i].OthersMessageInput = messagesInput
		tags := strings.Replace(stValue.Type().Field(i).Tag.Get(TagName), " ", "", -1)
		// get validator key
		if messagesInput[i].ValidatorKeyType == "" || len(tags) == 0 {
			continue
		}
		//get errors
		returnedErrors = append(returnedErrors, checkValidations(tags, messagesInput[i])...)
	}
	if flagTag {
		return append(returnedErrors, errors.New("Not found TAG: "+TagName))
	}
	return returnedErrors
}

// ValidateFields - will validate all structs with the tag "struct-validator" that you pass by argument
func ValidateFields(st interface{}, fields []string, messages map[string]map[string]string) (returnedErrors []error) {
	if st == nil {
		return append(returnedErrors, errors.New("The interface passed is nil"))
	}

	if len(fields) == 0 {
		return append(returnedErrors, errors.New("The field \"fields\" cannot be empty"))
	}
	flagTag := true
	// Let's convert the array of fields to validate, to map of string and bool
	namesMap := make(map[string]bool)
	for _, field := range fields {
		namesMap[strings.ToLower(field)] = true
	}

	stValue := reflect.ValueOf(st)
	// mount message input list
	messagesInput := make([]MessageInput, 0)
	for i := 0; i < stValue.NumField(); i++ {
		field := stValue.Field(i)
		var interfaceValue interface{}
		if fieldKind := field.Type().Kind(); (reflect.Int <= fieldKind && fieldKind <= reflect.Int64) || fieldKind == reflect.Float32 || fieldKind == reflect.Float64 {
			if fieldKind == reflect.Float32 || fieldKind == reflect.Float64 {
				interfaceValue = field.Float()
			} else {
				//convert int type to float64
				interfaceValue = float64(field.Int())
			}
		} else if reflect.Uint <= fieldKind && fieldKind <= reflect.Uintptr {
			//convert uint type to uint64
			interfaceValue = field.Uint()
		} else {
			//anothers types
			interfaceValue = field.Interface()
		}
		messagesInput = append(messagesInput, MessageInput{
			FieldName:        stValue.Type().Field(i).Name,
			FieldValue:       interfaceValue,
			CustomMessages:   messages,
			FieldType:        field.Type(),
			ValidatorKeyType: getValidatorKeyType(field.Type().String()),
		})
	}

	//get errors
	if len(messagesInput) > 0 {
		for i := 0; i < stValue.NumField(); i++ {
			_, tagLookup := stValue.Type().Field(i).Tag.Lookup(TagName)
			if tagLookup {
				flagTag = false
			}
			messagesInput[i].OthersMessageInput = messagesInput
			tags := strings.Replace(stValue.Type().Field(i).Tag.Get(TagName), " ", "", -1)
			valueTagJSON, validTagJSON := stValue.Type().Field(i).Tag.Lookup("json")
			if validTagJSON {
				// get validator key
				_, exist := namesMap[valueTagJSON]
				if exist {
					if messagesInput[i].ValidatorKeyType == "" || len(tags) == 0 || !exist {
						continue
					}
				} else {
					// get validator key
					_, exist := namesMap[strings.ToLower(stValue.Type().Field(i).Name)]
					if messagesInput[i].ValidatorKeyType == "" || len(tags) == 0 || !exist {
						continue
					}
				}
			} else {
				// get validator key
				_, exist := namesMap[strings.ToLower(stValue.Type().Field(i).Name)]
				if messagesInput[i].ValidatorKeyType == "" || len(tags) == 0 || !exist {
					continue
				}
			}
			//get errors
			returnedErrors = append(returnedErrors, checkValidations(tags, messagesInput[i])...)
		}
	}
	if flagTag {
		return append(returnedErrors, errors.New("Not found TAG: "+TagName))
	}
	return returnedErrors
}

// getValidatorKeyType - check the field type and returns the 'validator key type' associated to field type
func getValidatorKeyType(typeName string) string {
	if parts := strings.Split(typeName, "[]"); len(parts) == 2 {
		return nativeValidatorsKeyType["array"]
	}
	return nativeValidatorsKeyType[typeName]
}

// Will get the 'validator key type', get rules of field tag and get errors if they exist.
// A panic is throwed if the rule of 'messageInput' does not exists for the field 'validator key type'
func checkValidations(tags string, messageInput MessageInput) (errors []error) {
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

// SetTag - Seta o valor da tag que receber por parametro.
func SetTag(tag string) {
	TagName = tag
}
