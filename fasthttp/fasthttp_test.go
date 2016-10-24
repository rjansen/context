package fasthttp

import (
	"bytes"
	"context"
	"errors"
	"farm.e-pedion.com/repo/logger"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
	"testing"
)

func init() {
	os.Args = append(os.Args, "-ecf", "../test/etc/context/context.yaml")
	logger.Info("context.fasthttp_test.init")
}

func TestLogWrapper(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_log_wrapper"}`)
	serverMsg := []byte("context.fasthttp_test.TestLogWrapper")
	uri := "http://loghandle/"

	handler := Log(func(c context.Context, fc *fasthttp.RequestCtx) error {
		assert.NotEmpty(t, fc.PostBody())
		assert.True(t, bytes.Contains(fc.PostBody(), clientMsg))
		assert.NotEmpty(t, fc.URI())
		assert.True(t, strings.Contains(fc.URI().String(), uri))

		fc.Write(serverMsg)
		fc.SetStatusCode(fasthttp.StatusAccepted)
		return nil
	})
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.SetBody(clientMsg)
	ctx.Init(&req, nil, nil)
	c := context.Background()

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler(c, &ctx)
	})

	assert.Nil(t, resultErr)
	assert.NotEmpty(t, ctx.Response.Body())
	assert.True(t, bytes.Contains(ctx.Response.Body(), serverMsg))
	assert.Equal(t, fasthttp.StatusAccepted, ctx.Response.StatusCode())
}

func TestErrorWrapper(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_error_wrapper"}`)
	serverMsg := []byte("context.fasthttp_test.TestErrorWrapper")
	uri := "http://errorhadle"

	handler := Error(func(c context.Context, fc *fasthttp.RequestCtx) error {
		assert.NotEmpty(t, fc.PostBody())
		assert.True(t, bytes.Contains(fc.PostBody(), clientMsg))
		assert.NotEmpty(t, fc.URI())
		assert.True(t, strings.Contains(fc.URI().String(), uri))

		fc.Write(serverMsg)
		fc.SetStatusCode(fasthttp.StatusAccepted)
		return nil
	})
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.SetBody(clientMsg)
	ctx.Init(&req, nil, nil)
	c := context.Background()

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler(c, &ctx)
	})

	assert.Nil(t, resultErr)
	assert.NotEmpty(t, ctx.Response.Body())
	assert.True(t, bytes.Contains(ctx.Response.Body(), serverMsg))
	assert.Equal(t, fasthttp.StatusAccepted, ctx.Response.StatusCode())
}

func TestLogAndErrorWrapper(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_log_error_wrapper"}`)
	serverMsg := []byte("context.fasthttp_test.TestLogAndErrorWrapper")
	uri := "http://loghandle/errorhadle/"

	handler := Log(Error(func(c context.Context, fc *fasthttp.RequestCtx) error {
		assert.NotEmpty(t, fc.PostBody())
		assert.True(t, bytes.Contains(fc.PostBody(), clientMsg))
		assert.NotEmpty(t, fc.URI())
		assert.True(t, strings.Contains(fc.URI().String(), uri))

		fc.Write(serverMsg)
		fc.SetStatusCode(fasthttp.StatusAccepted)
		return nil
	}))
	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.SetBody(clientMsg)
	ctx.Init(&req, nil, nil)
	c := context.Background()

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler(c, &ctx)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusAccepted, ctx.Response.StatusCode())
	assert.NotEmpty(t, ctx.Response.Body())
	assert.True(t, bytes.Contains(ctx.Response.Body(), serverMsg))
}

func TestJSONResult(t *testing.T) {
	type mockJSON struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	media := &mockJSON{
		Username: "TestJSONResult",
		Name:     "Test JSON Result",
		Age:      15,
	}
	uri := "http://resultjson/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	ctx.Init(&req, nil, nil)

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = JSON(&ctx, fasthttp.StatusOK, media)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.NotEmpty(t, ctx.Response.Body())
	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Username)), "Response.Body does not contain media.Username")
	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Name)), "Response.Body does not contain media.Name")
	assert.True(t, bytes.Contains(ctx.Response.Header.ContentType(), []byte("application/json")), "Response.ContenType is not application/json")
}

func TestStatusResult(t *testing.T) {
	uri := "http://resultstatus/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	ctx.Init(&req, nil, nil)

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = Status(&ctx, fasthttp.StatusFound)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusFound, ctx.Response.StatusCode())
	assert.Empty(t, ctx.Response.Body())
}

func TestErrResult(t *testing.T) {
	uri := "http://resulterr/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	ctx.Init(&req, nil, nil)

	mockErrorMsg := "TestErrResult.Mock"
	mockError := errors.New(mockErrorMsg)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = Err(&ctx, mockError)
	})

	assert.NotNil(t, resultErr)
	assert.Equal(t, mockError, resultErr)
	assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
	assert.NotEmpty(t, ctx.Response.Body())
	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(mockErrorMsg)), "Response.Body does not contains the error message")
}
