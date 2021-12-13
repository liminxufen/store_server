package rest

import (
	"bytes"
	"encoding/json"
	"github.com/store_server/utils/errors"
)

func MarshalAPIJSON(code, message string, success bool, _struct interface{}) ([]byte, error) {
	x := &APIJSONFormat{}
	x.APIError.Code = code
	x.APIError.Message = message
	x.Success = success
	x.Result = _struct

	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	err := enc.Encode(x)
	if err != nil {
		return nil, errors.Errorf(err, "json encode err")
	}

	return buf.Bytes(), nil
}
