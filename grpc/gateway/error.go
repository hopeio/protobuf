/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package gateway

import (
	"strconv"

	"github.com/gin-gonic/gin"
	httpx "github.com/hopeio/gox/net/http"
	"github.com/hopeio/gox/net/http/grpc/gateway"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
)

var HttpError = func(ctx *gin.Context, err error) {
	s, ok := status.FromError(err)
	if !ok {
		grpclog.Warningf("Failed to convert error to status: %v", err)
	}
	delete(ctx.Request.Header, httpx.HeaderTrailer)
	errcodeHeader := strconv.Itoa(int(s.Code()))
	buf, contentType := gateway.DefaultMarshal(ctx, s)
	ctx.Header(httpx.HeaderContentType, contentType)
	ctx.Header(httpx.HeaderGrpcStatus, errcodeHeader)
	ctx.Header(httpx.HeaderErrorCode, errcodeHeader)
	if _, err := ctx.Writer.Write(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}

}
