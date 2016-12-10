package http

import (
	"bytes"
	"errors"
	// "farm.e-pedion.com/repo/context/media/proto"
	haki "farm.e-pedion.com/repo/context"
	"farm.e-pedion.com/repo/context/media/json"
	"farm.e-pedion.com/repo/logger"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func init() {
	os.Args = append(os.Args, "-ecf", "../test/etc/context/context.yaml")
	logger.Info("context.fast_test.init")
}

func TestLogWrapper(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_log_wrapper"}`)
	serverMsg := []byte("context.fasthttp_test.TestLogWrapper")
	uri := "http://loghandle.com/log"

	handler := Log(func(w http.ResponseWriter, r *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.NotEmpty(t, bodyBytes)
		assert.True(t, bytes.Contains(bodyBytes, clientMsg))
		assert.NotEmpty(t, r.URL.Path)
		assert.True(t, strings.Contains(uri, r.URL.Path))

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(serverMsg)
		return err
	})
	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(clientMsg))
	assert.Nil(t, err)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler.ServeHTTP(rec, req)
	})

	assert.Nil(t, resultErr)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.True(t, bytes.Contains(rec.Body.Bytes(), serverMsg))
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func TestLogWrapperErr(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_log_wrapper"}`)
	uri := "http://loghandle/logerr"
	mockErr := errors.New("MockErr")

	handler := Log(func(w http.ResponseWriter, r *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.NotEmpty(t, bodyBytes)
		assert.True(t, bytes.Contains(bodyBytes, clientMsg))
		assert.NotEmpty(t, r.URL.Path)
		assert.True(t, strings.Contains(uri, r.URL.Path))

		w.WriteHeader(http.StatusInternalServerError)
		return mockErr
	})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(clientMsg))
	assert.Nil(t, err)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler.ServeHTTP(rec, req)
	})

	assert.NotNil(t, resultErr)
	assert.Equal(t, mockErr, resultErr)
	assert.Empty(t, rec.Body.Bytes())
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestErrorWrapper(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_error_wrapper"}`)
	serverMsg := []byte("context.fasthttp_test.TestErrorWrapper")
	uri := "http://errorhadle/noerror"

	handler := Log(func(w http.ResponseWriter, r *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.NotEmpty(t, bodyBytes)
		assert.True(t, bytes.Contains(bodyBytes, clientMsg))
		assert.NotEmpty(t, r.URL.Path)
		assert.True(t, strings.Contains(uri, r.URL.Path))

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(serverMsg)
		return err
	})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(clientMsg))
	assert.Nil(t, err)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler.ServeHTTP(rec, req)
	})

	assert.Nil(t, resultErr)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.True(t, bytes.Contains(rec.Body.Bytes(), serverMsg))
	assert.Equal(t, http.StatusAccepted, rec.Code)
}

func TestErrorWrapperErr(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_error_wrapper"}`)
	uri := "http://errorhadle/error"
	mockErr := errors.New("MockErr")

	handler := Log(func(w http.ResponseWriter, r *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.NotEmpty(t, bodyBytes)
		assert.True(t, bytes.Contains(bodyBytes, clientMsg))
		assert.NotEmpty(t, r.URL.Path)
		assert.True(t, strings.Contains(uri, r.URL.Path))

		w.WriteHeader(http.StatusInternalServerError)
		return mockErr
	})

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(clientMsg))
	assert.Nil(t, err)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler.ServeHTTP(rec, req)
	})

	assert.NotNil(t, resultErr)
	assert.Equal(t, mockErr, resultErr)
	assert.Empty(t, rec.Body.Bytes())
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestLogAndErrorWrapper(t *testing.T) {
	clientMsg := []byte(`{"username": "mock_log_error_wrapper"}`)
	serverMsg := []byte("context.fasthttp_test.TestLogAndErrorWrapper")
	uri := "http://loghandle/loganderror/"

	handler := Log(Error(func(w http.ResponseWriter, r *http.Request) error {
		bodyBytes, err := ioutil.ReadAll(r.Body)
		assert.Nil(t, err)
		assert.NotEmpty(t, bodyBytes)
		assert.True(t, bytes.Contains(bodyBytes, clientMsg))
		assert.NotEmpty(t, r.URL.Path)
		assert.True(t, strings.Contains(uri, r.URL.Path))

		w.WriteHeader(http.StatusAccepted)
		_, err = w.Write(serverMsg)
		return err
	}))

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(clientMsg))
	assert.Nil(t, err)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = handler.ServeHTTP(rec, req)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, http.StatusAccepted, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.True(t, bytes.Contains(rec.Body.Bytes(), serverMsg))
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

	rec := httptest.NewRecorder()
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = JSON(rec, http.StatusOK, media)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.True(t, bytes.Contains(rec.Body.Bytes(), []byte(media.Username)), "Response.Body does not contain media.Username")
	assert.True(t, bytes.Contains(rec.Body.Bytes(), []byte(media.Name)), "Response.Body does not contain media.Name")
	assert.True(t, strings.Contains(rec.Header().Get("Content-Type"), "application/json"), "Response.ContenType is not application/json")
}

func TestJSONResultByAccept(t *testing.T) {
	type mockJSONData struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	type mockJSON struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
		Age  int    `json:"age"`
		Data []mockJSONData
	}
	media := mockJSON{
		ID:   1,
		Name: "mockJSON",
		Age:  24,
		Data: []mockJSONData{
			mockJSONData{
				ID:    1,
				Name:  "Proto Data Name",
				Email: "Proto Data Email",
			},
		},
	}

	uri := "http://resultjsontype/json"

	rec := httptest.NewRecorder()
	req, err := http.NewRequest("GET", uri, nil)
	req.Header.Set(haki.AcceptHeader, json.ContentType)
	assert.Nil(t, err)
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = WriteByAccept(rec, req, http.StatusOK, media)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.True(t, bytes.Contains(rec.Body.Bytes(), []byte(media.Name)), "Response.Body does not contain media.Name")
	assert.True(t, strings.Contains(rec.Header().Get(haki.ContentTypeHeader), "application/json"), "Response.ContenType is not application/json")
}

// func TestProtoResult(t *testing.T) {
// 	media := &proto.Store{
// 		Id:   1,
// 		Name: "Proto Buffer Store",
// 		Data: []*proto.Store_Data{
// 			&proto.Store_Data{
// 				Id:    1,
// 				Name:  "Proto Data Name",
// 				Email: "Proto Data Email",
// 			},
// 		},
// 	}

// 	uri := "http://resultprotobuf/"

// 	var ctx fasthttp.RequestCtx
// 	var req fasthttp.Request
// 	req.SetRequestURI(uri)
// 	ctx.Init(&req, nil, nil)

// 	var resultErr error
// 	assert.NotPanics(t, func() {
// 		resultErr = ProtoBuff(&ctx, fasthttp.StatusOK, media)
// 	})

// 	assert.Nil(t, resultErr)
// 	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
// 	assert.NotEmpty(t, ctx.Response.Body())
// 	//TODO: Check better if body content is the correct protocol buffer of the message
// 	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Name)), "Response.Body does not contain media.Name")
// 	assert.True(t, bytes.Contains(ctx.Response.Header.ContentType(), []byte("application/octet-stream")), "Response.ContenType is not application/octet-stream")
// }

// func TestProtoByContentType(t *testing.T) {
// 	media := &proto.Store{
// 		Id:   1,
// 		Name: "Proto Buffer Store",
// 		Data: []*proto.Store_Data{
// 			&proto.Store_Data{
// 				Id:    1,
// 				Name:  "Proto Data Name",
// 				Email: "Proto Data Email",
// 			},
// 		},
// 	}

// 	uri := "http://resultprototype/"

// 	var ctx fasthttp.RequestCtx
// 	var req fasthttp.Request
// 	req.SetRequestURI(uri)
// 	req.Header.Set("Accept", "application/octet-stream")
// 	ctx.Init(&req, nil, nil)

// 	var resultErr error
// 	assert.NotPanics(t, func() {
// 		resultErr = WriteByAccept(&ctx, fasthttp.StatusOK, media)
// 	})

// 	assert.Nil(t, resultErr)
// 	assert.Equal(t, fasthttp.StatusOK, ctx.Response.StatusCode())
// 	assert.NotEmpty(t, ctx.Response.Body())
// 	//TODO: Check better if body content is the correct protocol buffer of the message
// 	assert.True(t, bytes.Contains(ctx.Response.Body(), []byte(media.Name)), "Response.Body does not contain media.Name")
// 	assert.True(t, bytes.Contains(ctx.Response.Header.ContentType(), []byte("application/octet-stream")), "Response.ContenType is not application/octect-stream")
// }

func TestBytesResult(t *testing.T) {
	serverMsg := []byte("this is a mock server message")
	rec := httptest.NewRecorder()
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = Bytes(rec, http.StatusFound, serverMsg)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.Equal(t, rec.Body.Bytes(), serverMsg)
}

func TestStatusResult(t *testing.T) {
	rec := httptest.NewRecorder()
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = Status(rec, http.StatusFound)
	})

	assert.Nil(t, resultErr)
	assert.Equal(t, http.StatusFound, rec.Code)
	assert.Empty(t, rec.Body.Bytes())
}

func TestErrResult(t *testing.T) {
	mockErrorMsg := "TestErrResult.Mock"
	mockError := errors.New(mockErrorMsg)
	rec := httptest.NewRecorder()
	var resultErr error
	assert.NotPanics(t, func() {
		resultErr = Err(rec, mockError)
	})

	assert.NotNil(t, resultErr)
	assert.Equal(t, mockError, resultErr)
	assert.Equal(t, http.StatusInternalServerError, rec.Code)
	assert.NotEmpty(t, rec.Body.Bytes())
	assert.True(t, bytes.Contains(rec.Body.Bytes(), []byte(mockErrorMsg)), "Response.Body does not contains the error message")
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
	uri := "http://contentjson/read"

	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(rawMedia))
	assert.Nil(t, err)
	var readErr error
	var media mockJSON
	assert.NotPanics(t, func() {
		readErr = ReadJSON(req, &media)
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
	uri := "http://contentjson/json"
	req, err := http.NewRequest("POST", uri, bytes.NewBufferString(rawMedia))
	req.Header.Set(haki.ContentTypeHeader, json.ContentType)
	assert.Nil(t, err)

	var readErr error
	var media mockJSON
	assert.NotPanics(t, func() {
		readErr = ReadByContentType(req, &media)
	})

	assert.Nil(t, readErr)
	assert.NotZero(t, media)
	assert.Equal(t, "mock-raw.json", media.Username)
	assert.Equal(t, "Mock Raw Json", media.Name)
	assert.Equal(t, 35, media.Age)
}
