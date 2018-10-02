package form_test

import (
	"context"
	"io"
	"net/http"
	"testing"

	"github.com/go-qbit/qerror"
	"github.com/go-qbit/web/form"
	"github.com/go-qbit/web/form/field"
)

type testImpl struct{}

func (*testImpl) GetFields(ctx context.Context) []field.IField {
	return []field.IField{}
}

func (*testImpl) GetSubmitCaption(ctx context.Context) string {
	return "test"
}

func (*testImpl) RenderHTML(ctx context.Context, w io.Writer, f *form.Form) {
	panic("implement me")
}

func (*testImpl) OnSave(ctx context.Context, w http.ResponseWriter, f *form.Form) qerror.PublicError {
	panic("implement me")
}

func (*testImpl) OnComplete(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	panic("implement me")
}

func TestNew(t *testing.T) {
	f := form.New(context.Background(), &testImpl{})
	if f == nil {
		t.Fatal("Form is nil")
	}
}
