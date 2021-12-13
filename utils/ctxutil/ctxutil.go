package ctxutil

import "context"

const (
	KeyRequestID = 1 // 请求的id

	KeyLanguage = 2 // 请求的语言

	KeyClientIP = 3 // 请求的ip
)

func ContextDump(ctx context.Context) context.Context {
	if ctx == nil {
		return nil
	}

	ctxD := context.Background()
	ctxD = context.WithValue(ctxD, KeyRequestID, ctx.Value(KeyRequestID))
	ctxD = context.WithValue(ctxD, KeyLanguage, ctx.Value(KeyLanguage))
	ctxD = context.WithValue(ctxD, KeyClientIP, ctx.Value(KeyClientIP))
	return ctxD
}

// NewContextWithRequestID 返回带有request id的context对象
func NewContextWithRequestID(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, KeyRequestID, reqId)
}

// RequestIDFromContext 从ctx中得到消息的request id
func RequestIDFromContext(ctx context.Context) string {
	id, ok := ctx.Value(KeyRequestID).(string)
	if ok {
		return id
	}
	return ""
}

// NewContextWithLanguage 返回带有language的context对象
func NewContextWithLanguage(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, KeyLanguage, reqId)
}

// LanguageFromContext 从ctx中得到消息的language
func LanguageFromContext(ctx context.Context) string {
	id, ok := ctx.Value(KeyLanguage).(string)
	if ok {
		return id
	}
	return ""
}

func NewContextWithClientIP(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, KeyClientIP, reqId)
}

// LanguageFromContext 从ctx中得到消息的language
func ClientIPFromContext(ctx context.Context) string {
	id, ok := ctx.Value(KeyClientIP).(string)
	if ok {
		return id
	}
	return ""
}
