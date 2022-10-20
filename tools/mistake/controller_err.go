package mistake

// ReqErr 请求参数异常错误
type ReqErr struct {
	Message string
}

func (e *ReqErr) Error() string {
	return e.Message
}

// NewReqErr 创建请求参数异常错误
func NewReqErr(message string) *ReqErr {
	return &ReqErr{Message: message}
}

// StatusUnauthorizedErr 用户鉴权失败，没有token或token过期 错误
type StatusUnauthorizedErr struct {
	Message string
}

func (e *StatusUnauthorizedErr) Error() string {
	return e.Message
}

// NewStatusUnauthorizedErr 创建用户鉴权失败错误
func NewStatusUnauthorizedErr(message string) *StatusUnauthorizedErr {
	return &StatusUnauthorizedErr{Message: "unauthorized,error:" + message}
}
