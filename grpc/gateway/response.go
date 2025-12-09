/*
 * Copyright 2024 hopeio. All rights reserved.
 * Licensed under the MIT License that can be found in the LICENSE file.
 * @Created by jyb
 */

package gateway

import (
	"github.com/gin-gonic/gin"
	httpx "github.com/hopeio/gox/net/http"
	ginx "github.com/hopeio/gox/net/http/gin"
	"github.com/hopeio/gox/net/http/grpc"
	"github.com/hopeio/gox/net/http/grpc/gateway"
	"google.golang.org/protobuf/proto"
)

var ForwardResponseMessage = func(ctx *gin.Context, md grpc.ServerMetadata, message proto.Message) {
	if !message.ProtoReflect().IsValid() {
		ginx.Respond(ctx, &httpx.ErrResp{})
		return
	}
	err := gateway.ForwardResponseMessage(ctx.Writer, ctx.Request, md, message, gateway.Codec)
	if err != nil {
		HttpError(ctx, err)
		return
	}
}
