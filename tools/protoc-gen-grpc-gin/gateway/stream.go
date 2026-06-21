package gateway

import (
	"context"
	"encoding/base64"
	"encoding/binary"
	"strings"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/proto"
)

// ServerStream 将 gin.Context 包装为 gRPC server stream，
// 通过 HTTP/2（或 HTTP/1.1 chunked）逐帧写出消息。
// 编码方式由 Marshaler 决定（与一元 RPC 一致），可被 scaffold 自定义覆盖。
//
// 帧格式与 gRPC-over-HTTP2 一致：
//
//	[compressed-flag: 1 byte] [payload-length: 4 bytes big-endian] [payload: N bytes]
type ServerStream[T proto.Message] struct {
	ctx         *gin.Context
	header      metadata.MD
	trailers    metadata.MD
	started     bool
	contentType string
}

// NewServerStream 创建一个 gRPC server stream 适配器。
// 首次 Send 时自动设置 Content-Type 和 HTTP 200 状态码。
func NewServerStream[T proto.Message](ctx *gin.Context) *ServerStream[T] {
	return &ServerStream[T]{
		ctx:    ctx,
		header: metadata.MD{},
	}
}

// Send 将一条消息通过 Marshaler 编码后写入流式响应帧。
// 首次调用时自动声明 gRPC trailer key、设置 Content-Type 并写入 HTTP 200。
func (s *ServerStream[T]) Send(msg T) error {
	data, contentType := Marshaler(s.ctx, msg)
	if !s.started {
		s.started = true
		s.contentType = contentType
		// 声明 trailer key，必须在写 body 前完成
		s.ctx.Writer.Header().Add("Trailer", "Grpc-Status")
		s.ctx.Writer.Header().Add("Trailer", "Grpc-Message")
		s.ctx.Header("Content-Type", contentType)
		s.ctx.Writer.WriteHeader(200)
	}
	// gRPC length-prefixed frame: [compressed(1)][length(4)][payload(N)]
	frame := make([]byte, 5+len(data))
	frame[0] = 0 // not compressed
	binary.BigEndian.PutUint32(frame[1:5], uint32(len(data)))
	copy(frame[5:], data)
	if _, err := s.ctx.Writer.Write(frame); err != nil {
		return err
	}
	s.ctx.Writer.Flush()
	return nil
}

// Status 返回流式响应是否已成功启动（至少发送过一帧数据）。
// 用于在 defer 中决定 trailer 的 grpc-status 值。
func (s *ServerStream[T]) Status() bool {
	return s.started
}

// Context 返回底层 context（满足 grpc.ServerStream 接口）。
func (s *ServerStream[T]) Context() context.Context {
	return s.ctx.Request.Context()
}

// SetHeader 设置响应元数据（转换为 HTTP 响应头）。
func (s *ServerStream[T]) SetHeader(md metadata.MD) error {
	for k, vs := range md {
		for _, v := range vs {
			if strings.HasSuffix(k, "-bin") {
				s.ctx.Header(k, base64.StdEncoding.EncodeToString([]byte(v)))
			} else {
				s.ctx.Header(k, v)
			}
		}
	}
	return nil
}

// SendHeader 立即发送已设置的响应头。
// Gin 在首次 Write 时自动发送 header，此处为显式触发。
func (s *ServerStream[T]) SendHeader(md metadata.MD) error {
	_ = s.SetHeader(md)
	return nil
}

// SetTrailer 存储 trailer 元数据。
// 实际值通过 HTTP/2 trailer 在 handler 返回后写出。
// 生成的代码应在 defer 中调用 stream.Context().Writer.Header().Set() 来设置最终值。
func (s *ServerStream[T]) SetTrailer(md metadata.MD) {
	s.trailers = metadata.Join(s.trailers, md)
}

// SendMsg 实现 grpc.ServerStream 接口，将消息编码后写入流式响应帧。
func (s *ServerStream[T]) SendMsg(m any) error {
	if msg, ok := m.(T); ok {
		return s.Send(msg)
	}
	return status.Errorf(codes.Internal, "SendMsg: unexpected message type %T, expected %T", m, *new(T))
}

// RecvMsg 实现 grpc.ServerStream 接口（server streaming 不需要接收客户端消息）。
func (s *ServerStream[T]) RecvMsg(any) error {
	return status.Error(codes.Internal, "RecvMsg not supported on server streaming")
}
