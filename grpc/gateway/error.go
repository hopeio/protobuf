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
	const fallback = `{"code": 14, "msg": "failed to marshal error message"}`

	delete(ctx.Request.Header, httpx.HeaderTrailer)
	ctx.Header(httpx.HeaderContentType, gateway.Codec.ContentType(nil))
	ctx.Header(httpx.HeaderGrpcStatus, strconv.Itoa(int(s.Code())))
	se := &errors.ErrResp{Code: errors.ErrCode(s.Code()), Msg: s.Message()}
	buf, merr := gateway.Codec.Marshal(se)
	if merr != nil {
		grpclog.Infof("Failed to marshal error message %q: %v", se, merr)
		ctx.Status(http.StatusInternalServerError)
		ctx.Header(httpx.HeaderGrpcStatus, "14")
		ctx.Header(httpx.HeaderGrpcMessage, err.Error())
		if _, err := io.WriteString(ctx.Writer, fallback); err != nil {
			grpclog.Infof("Failed to write response: %v", err)
		}
		return
	}

	if _, err := ctx.Writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

}
