package rest

const (
	HeaderContentType    = "Content-Type"
	HeaderUserAgent      = "User-Agent"
	HeaderAcceptLanguage = "Accept-Language"
	HeaderAuthorization  = "Authorization" // jwt token
	HeaderCookie         = "Cookie"        // Cookie
	HeaderRequestId      = "X-request-id"
)

const (
	// ContentTypeJSON json类型
	ContentTypeJSON = "application/json; charset=UTF-8"
	// ContentTypeXML xml类型
	ContentTypeXML = "application/xml; charset=UTF-8"
	// ContentTypeForm form表单 application/x-www-form-urlencoded;charset=utf-8
	ContentTypeForm = "application/x-www-form-urlencoded; charset=UTF-8"
	// application/octet-stream
	ContentTypeStream = "application/octet-stream"
)

const (
	_DefaultUserAgent = "JOOX"
)
