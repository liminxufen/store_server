package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/store_server/utils/errors"
)

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (e APIError) Error() string {
	return fmt.Sprintf("Code:%s, Message:%s", e.Code, e.Message)
}

func (e APIError) GetCode() string {
	return e.Code
}

func (e APIError) GetMessage() string {
	return e.Message
}

type APIJSONFormat struct {
	APIError `json:"error,omitempty"`
	Success  bool        `json:"success"`
	Result   interface{} `json:"result"`
}

type APIJSONResult struct {
	APIError `json:"error,omitempty"`
	Success  bool            `json:"success"`
	Result   json.RawMessage `json:"result"`
}

type APIJSONResultV2 struct {
	Code   string          `json:"code,omitempty"`
	Result json.RawMessage `json:"result"`
}

func GetErrorJSONBody(code, message string) ([]byte, error) {
	return MarshalAPIJSON(code, message, false, struct{}{})
}

func GetSuccessJSONBody(_struct interface{}) ([]byte, error) {
	return MarshalAPIJSON("", "", true, _struct)
}

func SendSuccessJSON(w http.ResponseWriter, requestID string, _struct interface{}) error {
	body, err := GetSuccessJSONBody(_struct)
	if err != nil {
		return err
	}
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.Header().Add(HeaderRequestId, requestID)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func SendFailJSON(w http.ResponseWriter, requestID, code, message string) error {
	body, err := GetErrorJSONBody(code, message)
	if err != nil {
		return err
	}
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.Header().Add(HeaderRequestId, requestID)
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(body)
	if err != nil {
		return err
	}
	return nil
}

func SendBody(w http.ResponseWriter, requestID string, payload []byte) error {
	if payload == nil || len(payload) == 0 {
		return errors.Errorf(nil, "json body is nil")
	}
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.Header().Add(HeaderRequestId, requestID)
	w.WriteHeader(http.StatusOK)
	_, err := w.Write(payload)
	return err
}

func SendData(w http.ResponseWriter, requestID string, payload interface{}) error {
	if payload == nil {
		return errors.Errorf(nil, "json data is nil")
	}
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.Header().Add(HeaderRequestId, requestID)
	w.WriteHeader(http.StatusOK)
	b, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
}

func SendForbiddenJSON(w http.ResponseWriter, code, message string) error {
	w.Header().Set(HeaderContentType, ContentTypeJSON)
	w.WriteHeader(http.StatusForbidden)
	body, err := GetErrorJSONBody(code, message)
	if err != nil {
		_, _ = w.Write([]byte("forbidden"))
		return err
	}
	_, _ = w.Write(body)
	return nil
}
