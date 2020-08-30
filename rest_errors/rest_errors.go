package rest_errors

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type RestErr interface {
	Message() string
	Error() string
	Status() int
	Causes() []interface{}
}

type restErr struct {
	AMessage string        `json:"message"`
	AnError  string        `json:"error"`
	AStatus  int           `json:"status"`
	ACauses  []interface{} `json:"causes"`
}

func (e *restErr) Error() string {
	return fmt.Sprintf("Message: %s - Status: %d - Error: %s - Causes: [%v]",
		e.AMessage, e.AStatus, e.AnError, e.ACauses)
}

func (e *restErr) Message() string       { return e.AMessage }
func (e *restErr) Status() int           { return e.AStatus }
func (e *restErr) Causes() []interface{} { return e.ACauses }

func NewRestError(msg string, status int, err string, causes []interface{}) RestErr {
	return &restErr{
		AMessage: msg,
		AStatus:  status,
		AnError:  err,
		ACauses:  causes,
	}
}

func NewBadRequestError(msg string) RestErr {
	return &restErr{
		AMessage: msg,
		AnError:  "bad_request",
		AStatus:  http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) RestErr {
	return &restErr{
		AMessage: msg,
		AnError:  "not_found",
		AStatus:  http.StatusNotFound,
	}
}

func NewInternalServerError(msg string, err error) RestErr {
	result := &restErr{
		AMessage: msg,
		AnError:  "internal_server_error",
		AStatus:  http.StatusInternalServerError,
	}
	if err != nil {
		result.ACauses = append(result.ACauses, err.Error())
	}
	return result
}

func NewUnathorizedError(msg string) RestErr {
	return &restErr{
		//AMessage: "unable to retrive user infomation with given access_token",
		AMessage: msg,
		AStatus:  http.StatusUnauthorized,
		AnError:  "unauthorized",
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, fmt.Errorf("invalid json")
	}
	return &apiErr, nil
}
