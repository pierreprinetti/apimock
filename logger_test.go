package main

import (
	"net/http"
	"testing"
)

type testrw struct {
	headerCalled      bool
	writeCalled       bool
	writeHeaderCalled bool
}

func (trw *testrw) Header() http.Header {
	trw.headerCalled = true
	return http.Header{}
}
func (trw *testrw) Write(b []byte) (int, error) {
	trw.writeCalled = true
	return len(b), nil
}
func (trw *testrw) WriteHeader(int) {
	trw.writeHeaderCalled = true
}

func TestResponseWriterRecorderWrite(t *testing.T) {
	t.Run("calls the underlying rw.Header", func(t *testing.T) {
		rr := &responseWriterRecorder{ResponseWriter: &testrw{}}
		rr.Header()
		if !rr.ResponseWriter.(*testrw).headerCalled {
			t.Error("rw.Header has not been called")
		}
	})

	t.Run("calls the underlying rw.Write", func(t *testing.T) {
		rr := &responseWriterRecorder{ResponseWriter: &testrw{}}
		rr.Write([]byte("hey"))
		if !rr.ResponseWriter.(*testrw).writeCalled {
			t.Error("rw.Write has not been called")
		}
	})

	t.Run("calls the underlying rw.WriteHeader", func(t *testing.T) {
		rr := &responseWriterRecorder{ResponseWriter: &testrw{}}
		rr.WriteHeader(418)
		if !rr.ResponseWriter.(*testrw).writeHeaderCalled {
			t.Error("rw.WriteHeader has not been called")
		}
	})

	t.Run("sets the status when WriteHeader is called", func(t *testing.T) {
		rr := &responseWriterRecorder{ResponseWriter: &testrw{}}
		rr.WriteHeader(418)
		if want, have := 418, rr.status; want != have {
			t.Errorf("expected status %d, found %d", want, have)
		}
	})

	t.Run("sets the status to 200 when Write is called first", func(t *testing.T) {
		rr := &responseWriterRecorder{ResponseWriter: &testrw{}}
		rr.Write([]byte(""))
		if want, have := 200, rr.status; want != have {
			t.Errorf("expected status %d, found %d", want, have)
		}
	})
}
