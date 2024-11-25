/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package errcode

import (
	"github.com/hopeio/utils/errors/errcode"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"net/http"
)

func (x ErrCode) Code() int {
	return int(x)
}

func (x ErrCode) ErrRep() *errcode.ErrRep {
	return &errcode.ErrRep{Code: errcode.ErrCode(x), Msg: x.String()}
}

// example 实现
func (x ErrCode) GRPCStatus() *status.Status {
	return status.New(codes.Code(x), x.String())
}

func (x ErrCode) Msg(msg string) *errcode.ErrRep {
	return &errcode.ErrRep{Code: errcode.ErrCode(x), Msg: msg}
}

func (x ErrCode) Wrap(err error) *errcode.ErrRep {
	return &errcode.ErrRep{Code: errcode.ErrCode(x), Msg: err.Error()}
}

func (x ErrCode) Error() string {
	return x.String()
}

func (x ErrCode) Origin() errcode.ErrCode {
	return errcode.ErrCode(x)
}

/*func (x ErrCode) MarshalJSON() ([]byte, error) {
	return stringsi.ToBytes(`{"code":` + strconv.Itoa(int(x)) + `,"message":"` + x.String() + `"}`), nil
}

*/

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
