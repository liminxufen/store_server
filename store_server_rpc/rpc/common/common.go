package common

import (
	"fmt"
)

var RpcErrMap = map[int]string{
	2: "参数错误",
	3: "签名错误",
	4: "内部服务器错误",
}

//common rpc response
type CommRpcRsp struct {
	Code   int         `json:"code"`
	ErrMsg string      `json:"errmsg"`
	Data   interface{} `json:"data"`
}

func WrapRpcRsp(code int, inf string, payload interface{}, rsp *CommRpcRsp) {
	var msg string = ""
	if ix, ok := RpcErrMap[code]; ok && len(ix) != 0 {
		msg = ix
	} else {
		msg = inf
	}
	rsp.Code, rsp.ErrMsg, rsp.Data = code, msg, payload
	return
}

func CheckParamsIsNil(req interface{}) error {
	if req == nil {
		err := fmt.Errorf("request params is nil...")
		return err
	}
	return nil
}

/*----------------------------------------------------------------------------------------------------------*/
