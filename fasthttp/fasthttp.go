package fasthttp

import (
	"context"
	"farm.e-pedion.com/repo/context/media"
	"farm.e-pedion.com/repo/logger"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

//ResponseWriter is a wrapper function to store status and body length of the request
type ResponseWriter interface {
	http.ResponseWriter
	http.Flusher
	// Status returns the status code of the response or 200 if the response has
	// not been written (as this is the default response code in net/http)
	Status() int
	// Written returns whether or not the ResponseWriter has been written.
	Written() bool
	// Size returns the size of the response body.
	Size() int
}

// NewResponseWriter creates a ResponseWriter that wraps an http.ResponseWriter
func NewResponseWriter(w http.ResponseWriter) ResponseWriter {
	return &responseWriter{
		ResponseWriter: w,
	}
}

type responseWriter struct {
	http.ResponseWriter
	status int
	size   int
}

func (w *responseWriter) WriteHeader(s int) {
	w.status = s
	w.ResponseWriter.WriteHeader(s)
}

func (w *responseWriter) Write(b []byte) (int, error) {
	if !w.Written() {
		// The status will be 200 if WriteHeader has not been called yet
		w.WriteHeader(http.StatusOK)
	}
	size, err := w.ResponseWriter.Write(b)
	w.size += size
	return size, err
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}

func (w *responseWriter) Written() bool {
	return w.status != 0
}

func (w *responseWriter) Flush() {
	flusher, ok := w.ResponseWriter.(http.Flusher)
	if ok {
		if !w.Written() {
			// The status will be 200 if WriteHeader has not been called yet
			w.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

type ErrorHandler func(container context.Context, ctx *fasthttp.RequestCtx) error

func (h ErrorHandler) HandleRequest(ctx *fasthttp.RequestCtx) {
	if err := h(ctx); err != nil {
		ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	}
}

type LogHandler func(container context.Context, ctx *fasthttp.RequestCtx) error

func (h LogHandler) HandleRequest(ctx *fasthttp.RequestCtx) {
	start := time.Now()
	logger.Info("contex.Request",
		logger.Bytes("method", ctx.Method()),
		logger.Bytes("path", ctx.Path()),
	)
	logger.Debug("context.Context",
		logger.String("ctx", context.String()),
	)
	if err := h(ctx); err != nil {
		logger.Error("contex.LogHandler.Error",
			logger.Bytes("method", ctx.Method()),
			logger.Bytes("path", ctx.Path()),
			logger.Err(err),
		)
	}
	response := w.(ResponseWriter)
	logger.Info("context.Response",
		logger.Bytes("method", ctx.Method()),
		logger.Bytes("path", ctx.Path()),
		logger.String("status", http.StatusText(ctx.Response.StatusCode())),
		logger.Int("size", ctx.Response.Header.ContentLength()),
		logger.Time("requestTime", time.Since(start)),
	)
}

func JSON(ctx *fasthttp.RequestCtx, status int, result media.JSON) error {
	jsonBytes, err := result.ToBytes()
	if err != nil {
		return err
	}
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
	return nil
}

func Status(ctx *fasthttp.RequestCtx, status int) error {
	ctx.SetStatusCode(status)
	return nil
}

func Err(ctx *fasthttp.RequestCtx, err error) error {
	//w.WriteHeader(http.StatusInternalServerError)
	ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	return err
}

type Handler struct {
}

func (h Handler) JSON(ctx *fasthttp.RequestCtx, status int, result media.JSON) error {
	return JSON(ctx, status, result)
}

func (h Handler) Status(ctx *fasthttp.RequestCtx, status int) error {
	return Status(ctx, status)
}

func (h Handler) Err(ctx *fasthttp.RequestCtx, err error) error {
	return Err(ctx, err)
}
