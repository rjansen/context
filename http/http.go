package http

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"time"
    "farm.e-pedion.com/repo/context/media"
	"farm.e-pedion.com/repo/logger"
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
			rw.WriteHeader(http.StatusOK)
		}
		flusher.Flush()
	}
}

type ErrorHandler func(http.ResponseWriter, *http.Request) error

func (h ErrorHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type LogHandler func(http.ResponseWriter, *http.Request) error

func (h LogHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()
    rw := newResponseWriter(w)
	logger.Info("contex.Request",
		logger.String("method", r.Method),
		logger.String("path", r.URL.Path),
	)
    logger.Debug("context.Context", 
        logger.String("ctx", context.String()),
    )
	if err := h(w, r); err != nil {
		logger.Error("contex.LogHandler.Error",
			logger.String("method", r.Method),
			logger.String("path", r.URL.Path),
			logger.Err(err),
		)
	}
	response := w.(ResponseWriter)
	logger.Info("context.Response",
		logger.String("method", w.r.Method),
		logger.String("path", r.URL.Path),
        logger.String("status", http.StatusText(response.Status())), 
        logger.Int("size", response.Sixe()), 
        logger.Time("requestTime", time.Since(start),
    )
}

func JSON(w http.ResponseWriter, status int, result media.JSON) error {
    if err := result.Marshal(w); err != nil {
        return err
    }
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return nil
}

func Status(w http.ResponseWriter, status int) error {
	w.WriteHeader(status)
	return nil
}

func Err(w http.ResponseWriter, err error) error {
	//w.WriteHeader(http.StatusInternalServerError)
    http.Error(w, err.Error(), http.StatusInternalServerError)
	return err
}

type Handler struct {
}

func (h Handler) JSON(w http.ResponseWriter, status int, result media.JSON) error {
	return JSON(w, status, result)
}

func (h Handler) Status(w http.ResponseWriter, status int) error {
	return Status(w, status)
}

func (h Handler) Err(w http.ResponseWriter, err error) error {
	return Err(w, err)
}
