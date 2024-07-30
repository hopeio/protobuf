package errcode

import (
	"github.com/hopeio/utils/errors/errcode"
	stringsi "github.com/hopeio/utils/strings"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"net/http"
	"strconv"
)

func (x ErrCode) Code() int {
	return int(x)
}

func (x ErrCode) Rep() *ErrRep {
	return &ErrRep{Code: x, Message: x.String()}
}

// example 实现
func (x ErrCode) GrpcStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x ErrCode) Message(msg string) *ErrRep {
	return &ErrRep{Code: x, Message: msg}
}

func (x ErrCode) Wrap(err error) *ErrRep {
	return &ErrRep{Code: x, Message: err.Error()}
}

func (x ErrCode) Error() string {
	return x.String()
}

func ErrHandle(err interface{}) error {
	if e, ok := err.(*ErrRep); ok {
		return e
	}
	if e, ok := err.(ErrCode); ok {
		return e.Rep()
	}
	if e, ok := err.(*status.Status); ok {
		return e.Err()
	}
	if e, ok := err.(*errcode.ErrRep); ok {
		return e
	}
	if e, ok := err.(errcode.ErrCode); ok {
		return e.Rep()
	}
	if e, ok := err.(error); ok {
		return Unknown.Message(e.Error())
	}
	return Unknown.Rep()
}

func Code(err error) int {
	switch v := err.(type) {
	case *ErrRep:
		return int(v.Code)
	case ErrCode:
		return int(v)
	case *errcode.ErrRep:
		return int(v.Code)
	case errcode.ErrCode:
		return int(v)
	}
	return 0
}

func (x ErrCode) Origin() errcode.ErrCode {
	return errcode.ErrCode(x)
}

func (x *ErrRep) Error() string {
	return x.Message
}

func (x *ErrRep) GrpcStatus() *status.Status {
	return status.New(codes.Code(x.Code), x.Message)
}

func (x *ErrRep) MarshalJSON() ([]byte, error) {
	return stringsi.ToBytes(`{"code":` + strconv.Itoa(int(x.Code)) + `,"message":"` + x.Message + `"}`), nil
}

/*func (x ErrCode) MarshalJSON() ([]byte, error) {
	return stringsi.ToBytes(`{"code":` + strconv.Itoa(int(x)) + `,"message":"` + x.String() + `"}`), nil
}

*/

func FromError(err error) (s *ErrRep, ok bool) {
	if err == nil {
		return nil, true
	}
	if se, ok := err.(ErrRepInterface); ok {
		return se.ErrRep(), true
	}
	return NewErrRep(codes.Unknown, err.Error()), false
}

func SuccessRep() *ErrRep {
	return &ErrRep{Code: Success, Message: Success.Error()}
}

func NewErrRep[E ErrCodeGeneric](code E, msg string) *ErrRep {
	return &ErrRep{Code: ErrCode(code), Message: msg}
}

type ErrRepInterface interface {
	ErrRep() *ErrRep
}

type ErrCodeInterface interface {
}

type ErrCodeGeneric interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}

func HttpStatusFromCode(code ErrCode) int {
	switch code {
	case Success:
		return http.StatusOK
	case Canceled:
		return http.StatusRequestTimeout
	case Unknown:
		return http.StatusInternalServerError
	case InvalidArgument:
		return http.StatusBadRequest
	case DeadlineExceeded:
		return http.StatusGatewayTimeout
	case NotFound:
		return http.StatusNotFound
	case AlreadyExists:
		return http.StatusConflict
	case PermissionDenied:
		return http.StatusForbidden
	case Unauthenticated:
		return http.StatusUnauthorized
	case ResourceExhausted:
		return http.StatusTooManyRequests
	case FailedPrecondition:
		// Note, this deliberately doesn't translate to the similarly named '412 Precondition Failed' HTTP response status.
		return http.StatusBadRequest
	case Aborted:
		return http.StatusConflict
	case OutOfRange:
		return http.StatusBadRequest
	case Unimplemented:
		return http.StatusNotImplemented
	case Internal:
		return http.StatusInternalServerError
	case Unavailable:
		return http.StatusServiceUnavailable
	case DataLoss:
		return http.StatusInternalServerError
	}

	grpclog.Infof("Unknown gRPC error code: %v", code)
	return http.StatusInternalServerError
}
