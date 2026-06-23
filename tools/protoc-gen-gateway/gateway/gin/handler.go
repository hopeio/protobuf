package gin

import (
	"context"

	"github.com/gin-gonic/gin"
	ginx "github.com/hopeio/gox/net/http/gin"
	grpcx "github.com/hopeio/gox/net/http/grpc"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
)

var Bind = ginx.Bind
var Marshaler = gatewayx.DefaultMarshal

func withMetadataContext(ctx *gin.Context, stream interface {
	bindContext(context.Context)
}) context.Context {
	c := gatewayx.NewMetadataContext(ctx.Request.Context(), ctx.Writer.Header(), ctx.Request.Header)
	stream.bindContext(c)
	return c
}

func UnaryCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp]](
	handler func(context.Context, ReqPtr) (RespPtr, error),
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req

		if err := Bind(ctx, &req); err != nil {
			HttpError(ctx, err)
			return
		}

		stream := NewServerTransportStream[Req, Resp, ReqPtr, RespPtr](ctx)
		ctx.Request = ctx.Request.WithContext(withMetadataContext(ctx, stream))

		resp, err := handler(ctx.Request.Context(), &req)
		if err != nil {
			HttpError(ctx, err)
			return
		}

		HandleResponseMessage(ctx, resp)
	}
}

func ServerSideStreamCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp], S grpcx.ServerSideStream[Resp, RespPtr]](
	handler grpcx.ServerSideStreamHandler[Req, Resp, ReqPtr, RespPtr, S],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req Req
		var err error

		if err = Bind(ctx, &req); err != nil {
			HttpError(ctx, err)
			return
		}

		stream := NewServerStream[Req, Resp, ReqPtr, RespPtr](ctx)
		stream.forServerSendOnly()
		ctx.Request = ctx.Request.WithContext(withMetadataContext(ctx, stream))
		defer func() { stream.FinalizeTrailers(err) }()
		if err = handler(&req, any(stream).(S)); err != nil {
			HttpError(ctx, err)
			return
		}
	}
}

func ClientSideStreamCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp], S grpcx.ClientSideStream[Req, Resp, ReqPtr, RespPtr]](
	handler grpcx.ClientSideStreamHandler[Req, Resp, ReqPtr, RespPtr, S],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		stream := NewServerStream[Req, Resp, ReqPtr, RespPtr](ctx)
		stream.forClientRecv()
		ctx.Request = ctx.Request.WithContext(withMetadataContext(ctx, stream))

		if err := handler(any(stream).(S)); err != nil {
			HttpError(ctx, err)
			return
		}
	}
}

func BidiStreamCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp], S grpcx.BidiStream[Req, Resp, ReqPtr, RespPtr]](
	handler grpcx.BidiStreamHandler[Req, Resp, ReqPtr, RespPtr, S],
) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var err error

		stream := NewServerStream[Req, Resp, ReqPtr, RespPtr](ctx)
		ctx.Request = ctx.Request.WithContext(withMetadataContext(ctx, stream))
		defer func() { stream.FinalizeTrailers(err) }()
		if err = handler(any(stream).(S)); err != nil {
			HttpError(ctx, err)
			return
		}
	}
}
