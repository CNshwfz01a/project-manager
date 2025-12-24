package pkg

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	SystemErr    = 500
	MySqlErr     = 501
	LdapErr      = 505
	OperationErr = 506
	ValidatorErr = 412
)

type RspError struct {
	code int
	err  error
}

func (re *RspError) Error() string {
	return re.err.Error()
}

func (re *RspError) Code() int {
	return re.code
}

// NewRspError New
func NewRspError(code int, err error) *RspError {
	return &RspError{
		code: code,
		err:  err,
	}
}

// Success http 成功
func Success(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

// Err http 错误
func Err(c *gin.Context, err *RspError, data any) {
	c.JSON(err.Code(), gin.H{
		"code":  err.Code(),
		"error": err.Error(),
	})
}

// ReloadErr 重新加载错误
func ReloadErr(err any) *RspError {
	rspErr, ok := err.(*RspError)
	if !ok {
		rspError, ok := err.(error)
		if !ok {
			return &RspError{
				code: SystemErr,
				err:  fmt.Errorf("unknown error"),
			}
		}
		return &RspError{
			code: SystemErr,
			err:  rspError,
		}
	}
	return rspErr
}

// NewMySqlError mysql错误
func NewMySqlError(err error) *RspError {
	return NewRspError(MySqlErr, err)
}

// NewValidatorError 验证错误
func NewValidatorError(err error) *RspError {
	return NewRspError(ValidatorErr, err)
}

// 资源冲突 409
func NewConflictError(err error) *RspError {
	return NewRspError(http.StatusConflict, err)
}

// 无权操作 403
func NewUnauthorizedError(err error) *RspError {
	return NewRspError(http.StatusForbidden, err)
}

// 资源不存在 404
func NewNotFoundError() *RspError {
	return NewRspError(http.StatusNotFound, fmt.Errorf("资源不存在"))
}
