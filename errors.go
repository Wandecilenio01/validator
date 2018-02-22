package validator

import (
	"bytes"
	"errors"
	"text/template"
)

var (
	nativeMessages map[string]map[string]string
)

func init() {
	nativeMessages = map[string]map[string]string{
		// int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64
		"numeric": map[string]string{
			"min": "The {{.fieldName}} cannot be less than {{.ruleValue}}, the value informed was {{.value}}.",
			"max": "The {{.fieldName}} cannot be greater than {{.ruleValue}}, the value informed was {{.value}}.",
		},
		// array
		"array": map[string]string{
			"min":      "The {{.fieldName}} cannot have length less than {{.ruleValue}}, the value informed was {{.value}}.",
			"max":      "The {{.fieldName}} cannot have length greater than {{.ruleValue}}, the value informed was {{.value}}.",
			"distinct": "The {{.fieldName}} cannot have to be {{.ruleName}} and cannot have repeated itens, the value informed was {{.value}}.",
		},
		// string
		"string": map[string]string{
			"min":   "The {{.fieldName}} cannot have length less than {{.ruleValue}}, the length of informed value was \"{{.value}}\".",
			"max":   "The {{.fieldName}} cannot have length greater than {{.ruleValue}}, the length of informed value was \"{{.value}}\".",
			"email": "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"url":   "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"ipv4":  "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			// "ipv6":  "The {{.fieldName}} is not a valid {{.ruleValue}}, the informed value was {{.value}}.",
			"json":             "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha":            "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_dash":       "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_num":        "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_space":      "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_dash_space": "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"alpha_num_space":  "The {{.fieldName}} is not a valid {{.ruleName}}, the informed value was \"{{.value}}\".",
			"length":           "The {{.fieldName}} cannot have length different than {{.ruleValue}}, the length of informed value was \"{{.value}}\".",
			"regex":            "The {{.fieldName}} is not a valid {{.ruleName}}:{{.ruleValue}} , the informed value was {{.value}}.",
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

// GenerateErrorMessage -
func GenerateErrorMessage(messageInput MessageInput) error {
	//theres some message for the rule
	if messageInput.CustomMessages["all"] != nil && messageInput.CustomMessages["all"][messageInput.RuleName] != "" {
		messageInput.ValidatorKeyType = "all"
		return templateErrorMessage(messageInput)
	} else if messageInput.CustomMessages[messageInput.FieldName] != nil && messageInput.CustomMessages[messageInput.FieldName][messageInput.RuleName] != "" {
		messageInput.ValidatorKeyType = messageInput.FieldName
		return templateErrorMessage(messageInput)
	}
	messageInput.CustomMessages = nativeMessages
	return templateErrorMessage(messageInput)
}

func templateErrorMessage(messageInput MessageInput) error {
	var errorMessage bytes.Buffer
	if err := template.Must(template.New("ErrorMessageTemplate").Parse(messageInput.CustomMessages[messageInput.ValidatorKeyType][messageInput.RuleName])).Execute(&errorMessage, map[string]interface{}{"fieldName": messageInput.FieldName, "value": messageInput.FieldValue, "ruleValue": messageInput.RuleValue, "ruleName": messageInput.RuleName}); err != nil {
		panic(err)
	}
	return errors.New(errorMessage.String())
}
