/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package gateway

import (
	"io"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/hopeio/gox/errors"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/grpc/gateway"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var HttpError = func(ctx *gin.Context, err error) {
	s, _ := status.FromError(err)

	delete(ctx.Request.Header, httpx.HeaderTrailer)
	errcodeHeader := strconv.Itoa(int(s.Code()))
	ctx.Header(httpx.HeaderContentType, gateway.DefaultMarshaler.ContentType(s))
	ctx.Header(httpx.HeaderGrpcStatus, errcodeHeader)
	ctx.Header(httpx.HeaderErrorCode, errcodeHeader)
	se := &errors.ErrResp{Code: errors.ErrCode(s.Code()), Msg: s.Message()}
	buf, merr := gateway.DefaultMarshaler.Marshal(se)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", se, merr)
		ctx.Status(http.StatusInternalServerError)
		if _, err := io.WriteString(ctx.Writer, err.Error()); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	if _, err := ctx.Writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

}
