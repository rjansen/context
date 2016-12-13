package fast

import (
	"bytes"
	"context"
	"errors"
	"github.com/rjansen/haki/media/proto"
	"github.com/rjansen/l"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
	"os"
	"strings"
	"testing"
)

func init() {
	os.Args = append(os.Args, "-ecf", "../test/etc/haki/haki.yaml")
	l.Info("context.fast_test.init")
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

func TestLogWrapperErr(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_log_wrapper"}`)
	uri := "http://loghandle/"
	mockErr := errors.New("MockErr")

	handler := Log(func(c context.Context, fc *fasthttp.RequestCtx) error {
		assert.NotEmpty(t, fc.PostBody())
		assert.True(t, bytes.Contains(fc.PostBody(), clientMsg))
		assert.NotEmpty(t, fc.URI())
		assert.True(t, strings.Contains(fc.URI().String(), uri))

		fc.SetStatusCode(fasthttp.StatusInternalServerError)
		return mockErr
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

	assert.NotNil(t, resultErr)
	assert.Equal(t, mockErr, resultErr)
	assert.Empty(t, ctx.Response.Body())
	assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
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

func TestErrorWrapperErr(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_error_wrapper"}`)
	uri := "http://errorhadle"
	mockErr := errors.New("MockErr")

	handler := Error(func(c context.Context, fc *fasthttp.RequestCtx) error {
		assert.NotEmpty(t, fc.PostBody())
		assert.True(t, bytes.Contains(fc.PostBody(), clientMsg))
		assert.NotEmpty(t, fc.URI())
		assert.True(t, strings.Contains(fc.URI().String(), uri))

		fc.SetStatusCode(fasthttp.StatusInternalServerError)
		return mockErr
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

	assert.NotNil(t, resultErr)
	assert.Equal(t, mockErr, resultErr)
	assert.NotEmpty(t, ctx.Response.Body())
	assert.Equal(t, fasthttp.StatusInternalServerError, ctx.Response.StatusCode())
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

func TestJSONByContentType(t *testing.T) {
	media := &proto.Store{
		Id:   1,
		Name: "Proto Buffer Store",
		Data: []*proto.Store_Data{
			&proto.Store_Data{
				Id:    1,
				Name:  "Proto Data Name",
				Email: "Proto Data Email",
			},
		},
	}

	uri := "http://resultjsontype/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.Header.Set("Accept", "application/json")
	ctx.Init(&req, nil, nil)

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = WriteByAccept(&ctx, fasthttp.StatusOK, media)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.NotEmpty(t, ctx.Response.Body())
	//TODO: Check better if body content is the correct protocol buffer of the message
	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Name)), "Response.Body does not contain media.Name")
	assert.True(t, bytes.Contains(ctx.Response.Header.ContentType(), []byte("application/json")), "Response.ContenType is not application/json")
}

func TestProtoResult(t *testing.T) {
	media := &proto.Store{
		Id:   1,
		Name: "Proto Buffer Store",
		Data: []*proto.Store_Data{
			&proto.Store_Data{
				Id:    1,
				Name:  "Proto Data Name",
				Email: "Proto Data Email",
			},
		},
	}

	uri := "http://resultprotobuf/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	ctx.Init(&req, nil, nil)

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = ProtoBuff(&ctx, fasthttp.StatusOK, media)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.NotEmpty(t, ctx.Response.Body())
	//TODO: Check better if body content is the correct protocol buffer of the message
	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Name)), "Response.Body does not contain media.Name")
	assert.True(t, bytes.Contains(ctx.Response.Header.ContentType(), []byte("application/octet-stream")), "Response.ContenType is not application/octet-stream")
}

func TestProtoByContentType(t *testing.T) {
	media := &proto.Store{
		Id:   1,
		Name: "Proto Buffer Store",
		Data: []*proto.Store_Data{
			&proto.Store_Data{
				Id:    1,
				Name:  "Proto Data Name",
				Email: "Proto Data Email",
			},
		},
	}

	uri := "http://resultprototype/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.Header.Set("Accept", "application/octet-stream")
	ctx.Init(&req, nil, nil)

	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = WriteByAccept(&ctx, fasthttp.StatusOK, media)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
	assert.NotEmpty(t, ctx.Response.Body())
	//TODO: Check better if body content is the correct protocol buffer of the message
	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Name)), "Response.Body does not contain media.Name")
	assert.True(t, bytes.Contains(ctx.Response.Header.ContentType(), []byte("application/octet-stream")), "Response.ContenType is not application/octect-stream")
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

func TestJSONRead(t *testing.T) {
	type mockJSON struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	rawMedia := `
	{
		"username": "mock-raw.json",
		"name": "Mock Raw Json",
		"age": 35
	}
	`
	uri := "http://contentjson/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.SetBody([]byte(rawMedia))
	ctx.Init(&req, nil, nil)

	var readErr error
	var media mockJSON
	assert.NotPanics(t, func() {
		readErr = ReadJSON(&ctx, &media)
	})

	assert.Nil(t, readErr)
	assert.NotZero(t, media)
	assert.Equal(t, "mock-raw.json", media.Username)
	assert.Equal(t, "Mock Raw Json", media.Name)
	assert.Equal(t, 35, media.Age)
}

func TestJSONReadByContentType(t *testing.T) {
	type mockJSON struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	rawMedia := `
	{
		"username": "mock-raw.json",
		"name": "Mock Raw Json",
		"age": 35
	}
	`
	uri := "http://contentjson/"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.SetBody([]byte(rawMedia))
	req.Header.SetContentType("application/json")
	ctx.Init(&req, nil, nil)

	var readErr error
	var media mockJSON
	assert.NotPanics(t, func() {
		readErr = ReadByContentType(&ctx, &media)
	})

	assert.Nil(t, readErr)
	assert.NotZero(t, media)
	assert.Equal(t, "mock-raw.json", media.Username)
	assert.Equal(t, "Mock Raw Json", media.Name)
	assert.Equal(t, 35, media.Age)
}

func TestHandlerResult(t *testing.T) {
	type mockJSON struct {
		Username string `json:"username"`
		Name     string `json:"name"`
		Age      int    `json:"age"`
	}
	rawMedia := `
	{
		"username": "mock-raw.json",
		"name": "Mock Raw Json",
		"age": 35
	}
	`
	uri := "http://resultstatus/handler"

	var ctx fasthttp.RequestCtx
	var req fasthttp.Request
	req.SetRequestURI(uri)
	req.SetBody([]byte(rawMedia))
	req.Header.SetContentType("application/json")
	ctx.Init(&req, nil, nil)

	var resultErr error
	assert.NotPanics(t, func() {
		handler := Handler(
			func(c context.Context, ctx *fasthttp.RequestCtx) error {
				resultErr = Status(ctx, fasthttp.StatusFound)
				return resultErr
			},
		)
		handler(&ctx)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, fasthttp.StatusFound, ctx.Response.StatusCode())
	assert.Empty(t, ctx.Response.Body())
}
