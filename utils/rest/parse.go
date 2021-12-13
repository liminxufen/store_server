package rest

import (
	"encoding/json"
	"github.com/store_server/utils/errors"
)

func ParseResultJSON(body []byte) (*APIJSONResult, error) {
	res := APIJSONResult{}
	err := json.Unmarshal(body, &res)
	if err != nil {
		return nil, errors.Errorf(err, "parse")
	}
	return &res, nil
}
