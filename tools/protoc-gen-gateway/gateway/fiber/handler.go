package fiber

import (
	"context"

	"github.com/gofiber/fiber/v3"
	grpcx "github.com/hopeio/gox/net/http/grpc"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
)

func withMetadataContext(ctx fiber.Ctx, stream interface {
	bindContext(context.Context)
}) context.Context {
	c := gatewayx.NewMetadataContext(ctx.Context(), fiberReqHeader(ctx))
	stream.bindContext(c)
	return c
}

func UnaryCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp]](
	handler grpcx.GrpcHandler[Req, Resp, ReqPtr, RespPtr],
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var req Req

		if err := Bind(ctx, &req); err != nil {
			HttpError(ctx, err)
			return nil
		}

		stream := NewServerTransportStream[Req, Resp, ReqPtr, RespPtr](ctx)
		ctx.SetContext(withMetadataContext(ctx, stream))

		resp, err := handler(ctx.Context(), &req)
		if err != nil {
			HttpError(ctx, err)
			return nil
		}

		HandleResponseMessage(ctx, resp)
		return nil
	}
}

func ServerSideStreamCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp], S grpcx.ServerSideStream[Resp, RespPtr]](
	handler grpcx.ServerSideStreamHandler[Req, Resp, ReqPtr, RespPtr, S],
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var req Req
		var err error

		if err = Bind(ctx, &req); err != nil {
			HttpError(ctx, err)
			return nil
		}

		stream := NewServerStream[Req, Resp, ReqPtr, RespPtr](ctx)
		stream.forServerSendOnly()
		ctx.SetContext(withMetadataContext(ctx, stream))
		defer func() { stream.FinalizeTrailers(err) }()
		if err = handler(&req, any(stream).(S)); err != nil {
			HttpError(ctx, err)
			return nil
		}
		return nil
	}
}

func ClientSideStreamCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp], S grpcx.ClientSideStream[Req, Resp, ReqPtr, RespPtr]](
	handler grpcx.ClientSideStreamHandler[Req, Resp, ReqPtr, RespPtr, S],
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		stream := NewServerStream[Req, Resp, ReqPtr, RespPtr](ctx)
		stream.forClientRecv()
		ctx.SetContext(withMetadataContext(ctx, stream))

		if err := handler(any(stream).(S)); err != nil {
			HttpError(ctx, err)
			return nil
		}
		return nil
	}
}

func BidiStreamCall[Req, Resp any, ReqPtr grpcx.ProtoMessage[Req], RespPtr grpcx.ProtoMessage[Resp], S grpcx.BidiStream[Req, Resp, ReqPtr, RespPtr]](
	handler grpcx.BidiStreamHandler[Req, Resp, ReqPtr, RespPtr, S],
) fiber.Handler {
	return func(ctx fiber.Ctx) error {
		var err error

		stream := NewServerStream[Req, Resp, ReqPtr, RespPtr](ctx)
		ctx.SetContext(withMetadataContext(ctx, stream))
		defer func() { stream.FinalizeTrailers(err) }()
		if err = handler(any(stream).(S)); err != nil {
			HttpError(ctx, err)
			return nil
		}
		return nil
	}
}
