//go:generate ttgen -package form form.gtt
package form

import (
	"context"
	"io"
	"net/http"

	"github.com/go-qbit/qerror"
	"github.com/go-qbit/web/form/field"
)

type IFormImplementation interface {
	GetFields() []field.IField
	GetSubmitCaption(context.Context) string
	RenderHTML(context.Context, io.Writer, *Form)
	OnSave(context.Context, http.ResponseWriter, *Form) qerror.PublicError
	OnComplete(context.Context, http.ResponseWriter, *http.Request)
}

type Form struct {
	UsePost   bool
	LastError qerror.PublicError
	impl      IFormImplementation
	fields    []field.IField
	fieldsMap map[string]field.IField
}

func New(impl IFormImplementation) *Form {
	return &Form{
		impl:   impl,
		fields: impl.GetFields(),
	}
}

func (f *Form) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if f.UsePost && r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Form.Get("save") != "" {
		f.LastError = nil

		hasFieldsErrors := false
		for _, field := range f.fields {
			if field.Process(r.Context(), r.Form) != nil {
				hasFieldsErrors = true
			}
		}

		if !hasFieldsErrors {
			f.LastError = f.impl.OnSave(r.Context(), w, f)
			if f.LastError == nil {
				f.impl.OnComplete(r.Context(), w, r)
				return
			}
		}
	}

	f.impl.RenderHTML(r.Context(), w, f)
}

func (f *Form) ProcessForm(ctx context.Context, w io.Writer) {
	f.impl.RenderHTML(ctx, w, f)
}

func (f *Form) GetFields() []field.IField {
	return f.fields
}

func (f *Form) GetField(name string) field.IField {
	if f.fieldsMap == nil {
		f.fieldsMap = make(map[string]field.IField, len(f.fields))
		for _, field := range f.fields {
			f.fieldsMap[field.GetName()] = field
		}
	}

	return f.fieldsMap[name]
}

func (f *Form) GetSubmitCaption(ctx context.Context) string {
	return f.impl.GetSubmitCaption(ctx)
}

func (f *Form) GetValue(name string) interface{} {
	field := f.GetField(name)
	if field != nil {
		return field.GetValue()
	}

	return nil
}

func (f *Form) GetStringValue(name string) string {
	field := f.GetField(name)
	if field != nil {
		return field.GetStringValue()
	}

	return ""
}
