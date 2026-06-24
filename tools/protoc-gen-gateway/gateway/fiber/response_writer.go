package fiber

import (
	"net/http"

	"github.com/gofiber/fiber/v3"
)

// responseWriter 将 fiber.Ctx 适配为 http.ResponseWriter（含 Flush），供 gateway 流式写出复用。
type responseWriter struct {
	ctx fiber.Ctx
}

func (w *responseWriter) Header() http.Header {
	h := make(http.Header)
	w.ctx.Response().Header.VisitAll(func(key, value []byte) {
		h.Add(string(key), string(value))
	})
	return h
}

func (w *responseWriter) Write(p []byte) (int, error) {
	return w.ctx.Write(p)
}

func (w *responseWriter) WriteHeader(statusCode int) {
	w.ctx.Status(statusCode)
}

func (w *responseWriter) Flush() {
	w.ctx.Response().ImmediateHeaderFlush = true
}

var _ http.Flusher = (*responseWriter)(nil)

func newResponseWriter(ctx fiber.Ctx) *responseWriter {
	return &responseWriter{ctx: ctx}
}

func fasthttpRespHeader(ctx fiber.Ctx) http.Header {
	h := make(http.Header)
	ctx.Response().Header.VisitAll(func(key, value []byte) {
		h.Add(string(key), string(value))
	})
	return h
}

func fiberReqHeader(ctx fiber.Ctx) http.Header {
	h := make(http.Header)
	ctx.Request().Header.VisitAll(func(key, value []byte) {
		h.Add(string(key), string(value))
	})
	return h
}
