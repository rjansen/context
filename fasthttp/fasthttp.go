package fasthttp

import (
	"context"
	"farm.e-pedion.com/repo/context/media"
	"farm.e-pedion.com/repo/logger"
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

func errorHandle(handler HTTPHandlerFunc, c context.Context, fc *fasthttp.RequestCtx) error {
	if err := handler(c, fc); err != nil {
		fc.Error(err.Error(), fasthttp.StatusInternalServerError)
		return err
	}
	return nil
}

//ErrorHandler is a helper type to add exception control to other handlers
// type ErrorHandler func(context.Context, *fasthttp.RequestCtx) error

//HandleRequest is the HTTPHandler contract
// func (h ErrorHandler) HandleRequest(c context.Context, fc *fasthttp.RequestCtx) error {
// 	return errorHandle(HTTPHandlerFunc(h), c, fc)
// }

//Error wraps the provided HTTPHandlerFunc with exception control
func Error(handler HTTPHandlerFunc) HTTPHandlerFunc {
	return func(c context.Context, fc *fasthttp.RequestCtx) error {
		return errorHandle(handler, c, fc)
	}
}

func logHandle(handler HTTPHandlerFunc, c context.Context, fc *fasthttp.RequestCtx) error {
	start := time.Now()
	logger.Info("contex.Request",
		logger.Bytes("method", fc.Method()),
		logger.Bytes("path", fc.Path()),
	)
	logger.Debug("context.Context",
		logger.Bool("ctxIsNil", fc == nil),
		logger.Bool("containerIsNil", c == nil),
	)
	var err error
	if err = errorHandle(handler, c, fc); err != nil {
		logger.Error("contex.LogHandler.Error",
			logger.Bytes("method", fc.Method()),
			logger.Bytes("path", fc.Path()),
			logger.Err(err),
		)
	}
	logger.Info("context.Response",
		logger.Bytes("method", fc.Method()),
		logger.Bytes("path", fc.Path()),
		logger.String("status", http.StatusText(fc.Response.StatusCode())),
		logger.Int("size", fc.Response.Header.ContentLength()),
		logger.Duration("requestTime", time.Since(start)),
	)
	return err
}

//LogHandler is a helper type to add access logging control to other handlers
// type LogHandler func(context.Context, *fasthttp.RequestCtx) error

//HandleRequest is the HTTPHandler contract
// func (h LogHandler) HandleRequest(c context.Context, fc *fasthttp.RequestCtx) error {
// 	return logHandle(HTTPHandlerFunc(h), c, fc)
// }

//Log wraps the provided HTTPHandlerFunc with access logging control
func Log(handler HTTPHandlerFunc) HTTPHandlerFunc {
	return func(c context.Context, fc *fasthttp.RequestCtx) error {
		return logHandle(handler, c, fc)
	}
}

//JSON writes the provided json media to the response
func JSON(ctx *fasthttp.RequestCtx, status int, result media.JSON) error {
	jsonBytes, err := result.ToBytes()
	if err != nil {
		return err
	}
	ctx.SetBody(jsonBytes)
	ctx.SetContentType("application/json")
	ctx.SetStatusCode(status)
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

//Handler is a struct to add response helper function to other handlers
type Handler struct {
}

//JSON writes a json media to response
func (h Handler) JSON(ctx *fasthttp.RequestCtx, status int, result media.JSON) error {
	return JSON(ctx, status, result)
}

//Status writes the provided status to response
func (h Handler) Status(ctx *fasthttp.RequestCtx, status int) error {
	return Status(ctx, status)
}

//Err writes a error to response
func (h Handler) Err(ctx *fasthttp.RequestCtx, err error) error {
	return Err(ctx, err)
}
