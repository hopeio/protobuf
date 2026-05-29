package gateway

import (
	"github.com/hopeio/gox/net/http/gin"
	gatewayx "github.com/hopeio/gox/net/http/grpc/gateway"
)

var Bind = gin.Bind
var Marshaler = gatewayx.DefaultMarshal
