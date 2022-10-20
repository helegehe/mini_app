package mistake

import (
	"github.com/helegehe/mini_app/tools/logger"
)

// DaoErr 数据层通用异常错误
type DaoErr struct {
	Err   error
	Stack string
}

func (e *DaoErr) Error() string {
	return e.Err.Error()
}

// NewDaoErr 创建数据层通用异常错误
func NewDaoErr(err error) *DaoErr {
	return &DaoErr{
		Err:   err,
		Stack: logger.LogStack(2, 5),
	}
}

// DaoRecordNotFoundErr 数据层未查询到数据错误
type DaoRecordNotFoundErr struct{}

func (e *DaoRecordNotFoundErr) Error() string {
	return "Record Not Found"
}

// NewDaoRecordNotFoundErr 创建数据层未查询到数据错误
func NewDaoRecordNotFoundErr() *DaoRecordNotFoundErr {
	return &DaoRecordNotFoundErr{}
}

// IsDaoRecordNotFoundErr 返回是否是数据库记录不存在错误
func IsDaoRecordNotFoundErr(err error) bool {
	if err == nil {
		return false
	}
	_, ok := err.(*DaoRecordNotFoundErr)
	return ok
}
