package fiber

import (
	"net/http"
	"strconv"

	"github.com/gofiber/fiber/v2"
	httpx "github.com/hopeio/gox/net/http"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

var Marshaler = gatewayx.DefaultMarshal

var HandleResponseMessage = func(ctx *fiber.Ctx, message proto.Message) {
	_ = gatewayx.HandleResponseMessage(newResponseWriter(ctx), fiberRequest(ctx), message, Marshaler)
}

func fiberRequest(ctx *fiber.Ctx) *http.Request {
	req, _ := http.NewRequestWithContext(ctx.UserContext(), ctx.Method(), ctx.OriginalURL(), nil)
	ctx.Request().Header.VisitAll(func(k, v []byte) {
		req.Header.Add(string(k), string(v))
	})
	return req
}

var HttpError = func(ctx *fiber.Ctx, err error) {
	s, ok := status.FromError(err)
	if !ok {
		grpclog.Warningf("Failed to convert error to status: %v", err)
	}
	errcodeHeader := strconv.Itoa(int(s.Code()))
	buf, contentType := gatewayx.DefaultMarshal(ctx.UserContext(), s)
	ctx.Set(httpx.HeaderContentType, contentType)
	ctx.Set(httpx.HeaderGrpcStatus, errcodeHeader)
	ctx.Set(httpx.HeaderErrorCode, errcodeHeader)
	if err := ctx.Send(buf); err != nil {
		grpclog.Infof("Failed to write response: %v", err)
	}
}
