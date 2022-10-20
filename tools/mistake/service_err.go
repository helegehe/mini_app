package mistake

import (
	"fmt"
	"github.com/helegehe/mini_app/tools/logger"
	"net/http"
)

// ServiceErr 业务服务异常
type ServiceErr struct {
	HTTPCode int
	Message  string
	Err      error
	Stack    string
}

func (e *ServiceErr) Error() string {
	return e.Message
}

// New500ServiceErr 创建500
func New500ServiceErr(err error) *ServiceErr {
	return New500ServiceErrWithMessage(err, "InternalServerError")
}

// New500ServiceErrWithMessage 创建500和信息
func New500ServiceErrWithMessage(err error, message string) *ServiceErr {
	serviceErr := &ServiceErr{
		HTTPCode: http.StatusInternalServerError,
		Message:  message,
	}

	// 如果是数据库导致的异常，直接取用数据库调用栈信息，如果不是就记录当前调用栈信息
	if daoErr, ok := err.(*DaoErr); ok {
		serviceErr.Stack = daoErr.Stack
		serviceErr.Err = daoErr.Err
	} else if srvErr, ok := err.(*ServiceErr); ok {
		serviceErr.Stack = srvErr.Stack
		serviceErr.Err = srvErr.Err
	} else {
		serviceErr.Stack = logger.LogStack(2, 5)
		serviceErr.Err = err
	}

	return serviceErr
}

// New500ServiceErrWithAgreement 创建500和信息、message和error一致
func New500ServiceErrWithAgreement(message string) *ServiceErr {
	return New500ServiceErrWithMessage(fmt.Errorf(message), message)
}

// New400ServiceErr 构建 400 错误
func New400ServiceErr(err error, message string) *ServiceErr {
	return &ServiceErr{
		HTTPCode: http.StatusBadRequest,
		Message:  message,
		Err:      err,
		Stack:    logger.LogStack(2, 5),
	}
}

// New400ServiceErrOnlyMessage 构建 400 错误 只包含错误信息
func New400ServiceErrOnlyMessage(message string) *ServiceErr {
	return New400ServiceErr(nil, message)
}

// New404ServiceErr 构建 404 错误
func New404ServiceErr(message string) *ServiceErr {
	return &ServiceErr{
		HTTPCode: http.StatusNotFound,
		Message:  message,
		Stack:    logger.LogStack(2, 5),
	}
}

// New403ServiceErr 权限不足错误
func New403ServiceErr(message string) *ServiceErr {
	return &ServiceErr{
		HTTPCode: http.StatusForbidden,
		Message:  message,
		Stack:    logger.LogStack(2, 5),
	}
}

// NewServiceErr 构建 Diy 错误
func NewServiceErr(httpCode int, err error, message string) *ServiceErr {
	return &ServiceErr{
		HTTPCode: httpCode,
		Message:  message,
		Err:      err,
		Stack:    logger.LogStack(2, 5),
	}
}
