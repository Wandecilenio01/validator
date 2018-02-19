package validator

import (
	"bytes"
	"errors"
	"html/template"
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
		//string
		"string": map[string]string{
			"min": "The {{.fieldName}} cannot have length less than {{.ruleValue}}, the length of informed value was {{.value}}.",
			"max": "The {{.fieldName}} cannot have length greater than {{.ruleValue}}, the length of informed value was {{.value}}.",
		},
	}
}

func generateErrorMessage(messageInput MessageInput) error {
	//theres some message for the rule
	if messageInput.CustomMessages["all"] != nil && messageInput.CustomMessages["all"][messageInput.RuleName] != "" {
		messageInput.ValidatorKeyType = "all"
		return templateErrorMessage(messageInput)
	} else if messageInput.CustomMessages[messageInput.FieldName] != nil && messageInput.CustomMessages[messageInput.FieldName][messageInput.RuleName] != "" {
		return templateErrorMessage(messageInput)
	}
	messageInput.CustomMessages = nativeMessages
	return templateErrorMessage(messageInput)
}

func templateErrorMessage(messageInput MessageInput) error {
	var errorMessage bytes.Buffer
	if err := template.Must(template.New("ErrorMessageTemplate").Parse(messageInput.CustomMessages[messageInput.ValidatorKeyType][messageInput.RuleName])).Execute(&errorMessage, map[string]interface{}{"fieldName": messageInput.FieldName, "value": messageInput.FieldValue, "ruleValue": messageInput.RuleValue}); err != nil {
		panic(err)
	}
	return errors.New(errorMessage.String())
}
