package qhttp_test

import (
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/go-qbit/qerror"
	"github.com/go-qbit/web/qhttp"
)

type nullLogger struct{}

func (nullLogger) Print(v ...interface{}) {}

func TestError(t *testing.T) {
	err := errors.New("private error text")
	pubErr := qerror.ToPublic(err, "public message")

	qhttp.Logger = nullLogger{}

	w := httptest.NewRecorder()
	qhttp.Error(w, pubErr)

	if w.Body.String() != pubErr.PublicError()+"\n" {
		t.Fatalf("Invalid message, got '%s' expected '%s'", w.Body.String()+"\n", pubErr.PublicError())
	}
}
