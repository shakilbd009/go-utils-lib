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
	message string        `json:"message"`
	error   string        `json:"error"`
	status  int           `json:"status"`
	causes  []interface{} `json:"causes"`
}

func (e restErr) Error() string {
	return fmt.Sprintf("message: %s - status: %d - error: %s - causes: [%v]",
		e.message, e.status, e.error, e.causes)
}

func (e restErr) Message() string       { return e.message }
func (e restErr) Status() int           { return e.status }
func (e restErr) Causes() []interface{} { return e.causes }

func NewRestError(msg string, status int, err string, causes []interface{}) RestErr {
	return restErr{
		message: msg,
		status:  status,
		error:   err,
		causes:  causes,
	}
}

func NewBadRequestError(msg string) RestErr {
	return restErr{
		message: msg,
		error:   "bad_request",
		status:  http.StatusBadRequest,
	}
}

func NewNotFoundError(msg string) RestErr {
	return restErr{
		message: msg,
		error:   "not_found",
		status:  http.StatusNotFound,
	}
}

func NewInternalServerError(msg string, err error) RestErr {
	result := restErr{
		message: msg,
		error:   "internal_server_error",
		status:  http.StatusInternalServerError,
	}
	if err != nil {
		result.causes = append(result.causes, err.Error())
	}
	return result
}

func NewUnathorizedError(msg string) RestErr {
	return restErr{
		//message: "unable to retrive user infomation with given access_token",
		message: msg,
		status:  http.StatusUnauthorized,
		error:   "unauthorized",
	}
}

func NewRestErrorFromBytes(bytes []byte) (RestErr, error) {
	var apiErr restErr
	if err := json.Unmarshal(bytes, &apiErr); err != nil {
		return nil, fmt.Errorf("invalid json")
	}
	return apiErr, nil
}
