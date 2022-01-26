package errors

import (
	"baby-fried-rice/internal/pkg/kit/constant"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error interface {
	Error() string
	Code() constant.Code
}

type CommonError struct {
	err  error
	code constant.Code
}

func NewCommonError(code constant.Code) Error {
	return CommonError{
		code: code,
	}
}

func NewCommonErrorWrapper(wrappedError error) Error {
	var code = constant.CodeInternalError
	if w, ok := As(wrappedError); ok {
		code = w.Code()
	} else {
		ew := convertRpcErrToEdgeX(wrappedError)
		if w, ok = As(ew); ok {
			code = w.Code()
		}
	}
	return &CommonError{
		err:  wrappedError,
		code: code,
	}
}

func As(err error) (*CommonError, bool) {
	var ce = new(CommonError)
	ok := errors.As(err, &ce)

	return ce, ok
}

func NewInternalErr(wrappedError error) error {
	return NewCommonErr(constant.CodeInternalError, wrappedError)
}

func NewCommonErr(code constant.Code, wrappedError error) error {
	_, ok := As(wrappedError)
	if ok {
		// 已经是自定义错误则不在封装
		return wrappedError
	}
	return errors.WithStack(&CommonError{
		code: code,
		err:  wrappedError,
	})
}

func ConvertEdgeXErrToRpc(err error) error {
	errw := NewCommonErrorWrapper(err)
	st := status.New(codes.Code(errw.Code()), errw.Error())
	return st.Err()
}

func convertRpcErrToEdgeX(err error) error {
	if err == nil {
		return nil
	}

	st := status.Convert(err)
	if st == nil {
		return err
	}

	if st.Code() == codes.Unknown {
		return NewCommonErr(constant.CodeInternalError, err)
	}
	return NewCommonErr(constant.Code(st.Code()), err)
}

func (ce CommonError) Error() string {
	if ce.err == nil {
		return "nil"
	}
	return ce.err.Error()
}

func (ce CommonError) Code() constant.Code {
	return ce.code
}
