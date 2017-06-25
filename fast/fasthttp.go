package fast

import (
	"bytes"
	"context"
	"github.com/rjansen/haki"
	"github.com/rjansen/haki/media/json"
	"github.com/rjansen/haki/media/proto"
	"github.com/rjansen/l"
	"github.com/valyala/fasthttp"
	"net/http"
	"time"
)

//SimpleHTTPHandler is a contract for fast http handlers
type SimpleHTTPHandler interface {
	HandleRequest(*fasthttp.RequestCtx)
}

//HTTPHandlerFunc is a function to handle fasthttp requrests
type HTTPHandlerFunc func(context.Context, *fasthttp.RequestCtx) error

//HandleRequest is the contract with HTTPHandler interface
func (h HTTPHandlerFunc) HandleRequest(c context.Context, fc *fasthttp.RequestCtx) error {
	return h(c, fc)
}

//HTTPHandler is a contract for fast http handlers
type HTTPHandler interface {
	HandleRequest(context.Context, *fasthttp.RequestCtx) error
}

//Handler wraps a library handler func nto a fasthttp handler func
func Handler(handler HTTPHandlerFunc) fasthttp.RequestHandler {
	return func(ctx *fasthttp.RequestCtx) {
		c, cancel := context.WithCancel(context.Background())
		defer cancel()
		handler(c, ctx)
	}
}

func errorHandle(handler HTTPHandlerFunc, c context.Context, fc *fasthttp.RequestCtx) error {
	if err := handler(c, fc); err != nil {
		fc.Error(err.Error(), fasthttp.StatusInternalServerError)
		return err
	}
	return nil
}

//ErrorHandler is a helper type to add exception control to other handlers
type ErrorHandler func(context.Context, *fasthttp.RequestCtx) error

//HandleRequest is the HTTPHandler contract
func (h ErrorHandler) HandleRequest(c context.Context, fc *fasthttp.RequestCtx) error {
	return errorHandle(HTTPHandlerFunc(h), c, fc)
}

//Error wraps the provided HTTPHandlerFunc with exception control
func Error(handler HTTPHandlerFunc) HTTPHandlerFunc {
	return func(c context.Context, fc *fasthttp.RequestCtx) error {
		return errorHandle(handler, c, fc)
	}
}

func logHandle(handler HTTPHandlerFunc, c context.Context, fc *fasthttp.RequestCtx) error {
	start := time.Now()
	logger := l.WithFields(
		l.Int64("tid", int64(fc.ConnID())),
		l.Int64("rid", int64(fc.ConnRequestNum())),
		l.Bytes("method", fc.Method()),
		l.Bytes("path", fc.Path()),
		l.String("auth", "anonymous"),
	)
	logger.Info("contex.Request")
	logger.Debug("context.Context",
		l.Bool("ctxIsNil", fc == nil),
		l.Bool("containerIsNil", c == nil),
	)
	c = context.WithValue(c, "log", logger)
	var err error
	if err = handler(c, fc); err != nil {
		logger.Error("contex.LogHandler.Error",
			l.Err(err),
		)
	}
	logger.Info("context.Response",
		l.String("status", http.StatusText(fc.Response.StatusCode())),
		l.Int("size", fc.Response.Header.ContentLength()),
		l.Duration("requestTime", time.Since(start)),
	)
	return err
}

//LogHandler is a helper type to add access logging control to other handlers
type LogHandler func(context.Context, *fasthttp.RequestCtx) error

//HandleRequest is the HTTPHandler contract
func (h LogHandler) HandleRequest(c context.Context, fc *fasthttp.RequestCtx) error {
	return logHandle(HTTPHandlerFunc(h), c, fc)
}

//Log wraps the provided HTTPHandlerFunc with access logging control
func Log(handler HTTPHandlerFunc) HTTPHandlerFunc {
	return func(c context.Context, fc *fasthttp.RequestCtx) error {
		return logHandle(handler, c, fc)
	}
}

//ReadByContentType reads data from context using the Content-Type header to define the media type
func ReadByContentType(ctx *fasthttp.RequestCtx, data interface{}) error {
	contentType := ctx.Request.Header.ContentType()
	switch {
	case bytes.Contains(contentType, []byte(json.ContentType)):
		return ReadJSON(ctx, data)
	case bytes.Contains(contentType, []byte(proto.ContentType)):
		return ReadProtoBuff(ctx, data)
	default:
		return haki.ErrInvalidContentType
	}
}

//WriteByAccept writes data to context using the Accept header to define the media type
func WriteByAccept(ctx *fasthttp.RequestCtx, status int, result interface{}) error {
	contentType := ctx.Request.Header.Peek(haki.AcceptHeader)
	switch {
	case bytes.Contains(contentType, []byte(json.ContentType)):
		return JSON(ctx, status, result)
	case bytes.Contains(contentType, []byte(proto.ContentType)):
		return ProtoBuff(ctx, status, result)
	default:
		return haki.ErrInvalidAccept
	}
}

//JSON writes the provided json media to the response
func JSON(ctx *fasthttp.RequestCtx, status int, result interface{}) error {
	jsonBytes, err := json.MarshalBytes(result)
	if err != nil {
		return err
	}
	ctx.SetBody(jsonBytes)
	ctx.SetContentType(json.ContentType)
	ctx.SetStatusCode(status)
	return nil
}

//ReadJSON unmarshals from provided context a json media into data
func ReadJSON(ctx *fasthttp.RequestCtx, data interface{}) error {
	if err := json.UnmarshalBytes(ctx.PostBody(), data); err != nil {
		return err
	}
	return nil
}

//ProtoBuff writes the provided protocol buffer media to the response
func ProtoBuff(ctx *fasthttp.RequestCtx, status int, result interface{}) error {
	protoBytes, err := proto.MarshalBytes(result)
	if err != nil {
		return err
	}
	ctx.SetBody(protoBytes)
	ctx.SetContentType(proto.ContentType)
	ctx.SetStatusCode(status)
	return nil
}

//ReadProtoBuff unmarshals from provided context a protocol buffer media into data
func ReadProtoBuff(ctx *fasthttp.RequestCtx, data interface{}) error {
	if err := proto.UnmarshalBytes(ctx.PostBody(), data); err != nil {
		return err
	}
	return nil
}

//Status writes the provided status to the response
func Status(ctx *fasthttp.RequestCtx, status int) error {
	ctx.SetStatusCode(status)
	return nil
}

//Err writes the provided  error to the response
func Err(ctx *fasthttp.RequestCtx, err error) error {
	//w.WriteHeader(http.StatusInternalServerError)
	ctx.Error(err.Error(), fasthttp.StatusInternalServerError)
	return err
}

//BaseHandler is a struct to add response helper function to other handlers
type BaseHandler struct {
}

//JSON writes a json media to response
func (h BaseHandler) JSON(ctx *fasthttp.RequestCtx, status int, result interface{}) error {
	return JSON(ctx, status, result)
}

//Status writes the provided status to response
func (h BaseHandler) Status(ctx *fasthttp.RequestCtx, status int) error {
	return Status(ctx, status)
}

//Err writes a error to response
func (h BaseHandler) Err(ctx *fasthttp.RequestCtx, err error) error {
	return Err(ctx, err)
}
