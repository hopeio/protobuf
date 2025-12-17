package gateway

import (
	"github.com/hopeio/gox/net/http/gin/binding"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
)

var Bind = binding.Bind
var Marshaler = gatewayx.DefaultMarshaler
