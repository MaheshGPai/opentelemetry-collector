// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package mw

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type MockTimeoutHandler struct {
}

func (*MockTimeoutHandler) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	time.Sleep(2 * time.Nanosecond)
	rw.Write([]byte("hello"))
}

type MockHandlerWithPanic struct {
}

func (*MockHandlerWithPanic) ServeHTTP(rw http.ResponseWriter, rq *http.Request) {
	panic("panic")
}

// implement http.ResponseWriter
type MockResponseWriter struct {
	resp   string
	status int
}

func (*MockResponseWriter) Header() http.Header {
	return make(http.Header)
}
func (mr *MockResponseWriter) Write(res []byte) (int, error) {
	mr.resp = string(res)
	return 200, nil
}
func (mr *MockResponseWriter) WriteHeader(statusCode int) {
	mr.status = statusCode
}
func TestTimeout(t *testing.T) {
	mock := &MockTimeoutHandler{}
	handler := TimeoutHandler(mock, 1*time.Nanosecond)
	mr := &MockResponseWriter{}
	handler.ServeHTTP(mr, &http.Request{})
	assert.Equal(t, mr.resp, "{\"error\":408,\"message\":\"request timed out\"}")
	assert.Equal(t, mr.status, 408)
}

func TestNoTimeout(t *testing.T) {
	mock := &MockTimeoutHandler{}
	handler := TimeoutHandler(mock, 1*time.Millisecond)
	mr := &MockResponseWriter{}
	handler.ServeHTTP(mr, &http.Request{})
	assert.Equal(t, mr.resp, "hello")
	assert.Equal(t, mr.status, 200)
}

func TestInternalServerError(t *testing.T) {
	mock := &MockHandlerWithPanic{}
	handler := TimeoutHandler(mock, 1*time.Millisecond)
	mr := &MockResponseWriter{}
	handler.ServeHTTP(mr, &http.Request{})
	assert.Equal(t, mr.resp, "{\"error\":500,\"message\":\"internal server error\"}")
	assert.Equal(t, mr.status, 500)
}
