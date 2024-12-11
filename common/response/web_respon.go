package response

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type (
	FailedResponseMessage struct {
		Message string      `json:"message"`
		Status  string      `json:"status"`
		Errors  interface{} `json:"error"`
		Code    int         `json:"code"`
	}

	SuccessResponseMessage struct {
		Data    interface{} `json:"data"`
		Message string      `json:"message"`
		Status  string      `json:"status"`
		Code    int         `json:"code"`
	}

	ValidationJsonResponseMessage struct {
		FailedField string      `json:"failed_field"`
		Tag         string      `json:"tag"`
		Value       interface{} `json:"value"`
	}
)

func BuildSuccessResponseMessage(message string, code int, data interface{}) SuccessResponseMessage {
	return SuccessResponseMessage{
		Data:    data,
		Message: message,
		Status:  "success",
		Code:    code,
	}
}

func BuildFailedResponseMessage(message string, code int, errors interface{}) FailedResponseMessage {
	return FailedResponseMessage{
		Errors:  errors,
		Message: message,
		Status:  "failed",
		Code:    code,
	}
}

var validate = validator.New()

func (v *FailedResponseMessage) Error() string {
	return fmt.Sprintf("Code: %d, Message: %s, Data: %v", v.Code, v.Message, v.Errors)
}

func ValidateBodyRequest(input interface{}) []ValidationJsonResponseMessage {

	var errors []ValidationJsonResponseMessage

	if err := validate.Struct(input); err != nil {
		for _, err := range err.(validator.ValidationErrors) {
			var error ValidationJsonResponseMessage
			error.FailedField = err.Field()
			error.Tag = err.ActualTag()
			error.Value = err.Value()
			errors = append(errors, error)
		}
	}

	return errors
}
