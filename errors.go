package validator

import (
	"bytes"
	"errors"
	"text/template"
)

var (
	// nativeMessages - relations between 'validator key type' and 'rule' and 'message'
	nativeMessages map[string]map[string]string
)

// definition of nativeMessages attr
func init() {
	nativeMessages = map[string]map[string]string{
		// int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64
		"numeric": map[string]string{
			"min": "The {{.fieldName}} cannot be less than {{.ruleValue}}, the value informed was {{.value}}.",
			"max": "The {{.fieldName}} cannot be greater than {{.ruleValue}}, the value informed was {{.value}}.",
		},
		// array's in general
		"array": map[string]string{
			"min":                  "The {{.fieldName}} cannot have length less than {{.ruleValue}}, the value informed was {{.value}}.",
			"max":                  "The {{.fieldName}} cannot have length greater than {{.ruleValue}}, the value informed was {{.value}}.",
			"distinct":             "The {{.fieldName}} cannot have to be {{.ruleName}} and cannot have repeated itens, the value informed was {{.value}}.",
			"required_with":        "The {{.fieldName}} is not a valid {{.ruleName}}, because if at leat one of that fields: ({{.ruleValue}}) is filled, then {{.fieldName}} needs to be filled too.",
			"required_with_all":    "The {{.fieldName}} is not a valid {{.ruleName}}, because if all fields: ({{.ruleValue}}) are filled, then {{.fieldName}} needs to be filled too.",
			"required_without":     "The {{.fieldName}} is not a valid {{.ruleName}}, because if at least one that fields: ({{.ruleValue}}) are not filled, then {{.fieldName}} needs to be filled.",
			"required_without_all": "The {{.fieldName}} is not a valid {{.ruleName}}, because if all fields: ({{.ruleValue}}) are not filled, then {{.fieldName}} needs to be filled.",
		},
		// string
		"string": map[string]string{
			"min":   "The {{.fieldName}} cannot have length less than {{.ruleValue}}, the informed value was \"{{.value}}\".",
			"max":   "The {{.fieldName}} cannot have length greater than {{.ruleValue}}, the informed value was \"{{.value}}\".",
			"email": "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"url":   "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"ipv4":  "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			// "ipv6":  "The {{.fieldName}} is not a valid {{.ruleValue}}, the informed value was {{.value}}.",
			"json":                 "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha":                "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_dash":           "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_num":            "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_space":          "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_dash_space":     "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_num_space":      "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"length":               "The {{.fieldName}} cannot have length different than {{.ruleValue}}, the length of informed value was \"{{.value}}\".",
			"regex":                "The {{.fieldName}} is not a valid {{.ruleName}}:{{.ruleValue}} , the informed value was {{.value}}.",
			"required_with":        "The {{.fieldName}} is not a valid {{.ruleName}}, because if at leat one of that fields: ({{.ruleValue}}) is filled, then {{.fieldName}} needs to be filled too.",
			"required_with_all":    "The {{.fieldName}} is not a valid {{.ruleName}}, because if all fields: ({{.ruleValue}}) are filled, then {{.fieldName}} needs to be filled too.",
			"required_without":     "The {{.fieldName}} is not a valid {{.ruleName}}, because if at least one that fields: ({{.ruleValue}}) are not filled, then {{.fieldName}} needs to be filled.",
			"required_without_all": "The {{.fieldName}} is not a valid {{.ruleName}}, because if all fields: ({{.ruleValue}}) are not filled, then {{.fieldName}} needs to be filled.",
		},
		// timestamp
		"timestamp": map[string]string{
			"equal":                "The {{.fieldName}} have to be equals to {{.ruleValue}}, the timestamp informed was {{.value}}.",
			"after":                "The {{.fieldName}} have to be after {{.ruleValue}}, the timestamp informed was {{.value}}.",
			"before":               "The {{.fieldName}} have to be before {{.ruleValue}}, the timestamp informed was {{.value}}.",
			"equal_date":           "The {{.fieldName}} have to be equals to {{.ruleValue}}, the date informed was {{.value}}.",
			"after_date":           "The {{.fieldName}} have to be after {{.ruleValue}}, the date informed was {{.value}}.",
			"before_date":          "The {{.fieldName}} have to be before {{.ruleValue}}, the date informed was {{.value}}.",
			"after_or_equal":       "The {{.fieldName}} have to be after or equals to {{.ruleValue}}, the timestamp informed was {{.value}}.",
			"before_or_equal":      "The {{.fieldName}} have to be before or equals to {{.ruleValue}}, the timestamp informed was {{.value}}.",
			"after_or_equal_date":  "The {{.fieldName}} have to be after or equals to {{.ruleValue}}, the date informed was {{.value}}.",
			"before_or_equal_date": "The {{.fieldName}} have to be before or equals to {{.ruleValue}}, the date informed was {{.value}}.",
		},
	}
}

// GenerateErrorMessage - Generate an error using the messageInput.CustomMessages or the nativeMessages
func GenerateErrorMessage(messageInput MessageInput) error {
	if messageInput.CustomMessages["*"] != nil && messageInput.CustomMessages["*"][messageInput.RuleName] != "" {
		//there's some custom message for every field and that especific rule
		messageInput.ValidatorKeyType = "*"
		return TemplateErrorMessage(messageInput)
	} else if messageInput.CustomMessages[messageInput.FieldName] != nil && messageInput.CustomMessages[messageInput.FieldName][messageInput.RuleName] != "" {
		//there's some custom message for that especific field and rule
		messageInput.ValidatorKeyType = messageInput.FieldName
		return TemplateErrorMessage(messageInput)
	}
	//there's no custom message for that field and rule
	messageInput.CustomMessages = nativeMessages
	return TemplateErrorMessage(messageInput)
}

// TemplateErrorMessage - Returns an error with a templated string using attributes of messageInput parameter
func TemplateErrorMessage(messageInput MessageInput) error {
	var errorMessage bytes.Buffer
	if err := template.Must(template.New("ErrorMessageTemplate").Parse(messageInput.CustomMessages[messageInput.ValidatorKeyType][messageInput.RuleName])).Execute(&errorMessage, map[string]interface{}{"fieldName": messageInput.FieldName, "value": messageInput.FieldValue, "ruleValue": messageInput.RuleValue, "ruleName": messageInput.RuleName}); err != nil {
		panic(err)
	}
	return errors.New(errorMessage.String())
}

// SetNativeMessages - Sets a custom native message
func SetNativeMessages(NewNativeMessages map[string]map[string]string) {
	nativeMessages = NewNativeMessages
}
