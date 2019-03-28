//go:generate ttgen -package form form.gtt
package form

import (
	"context"
	"io"
	"net/http"

	"github.com/go-qbit/qerror"
	"github.com/go-qbit/web/form/field"
	"github.com/go-qbit/web/formctx"
)

type IFormImplementation interface {
	GetFields(ctx context.Context) []field.IField
	GetCaption(ctx context.Context) string
	GetSubmitCaption(ctx context.Context) string
	RenderHTML(ctx context.Context, w io.Writer, f *Form)
	OnSave(ctx context.Context, w http.ResponseWriter, f *Form) qerror.PublicError
	OnComplete(ctx context.Context, w http.ResponseWriter, r *http.Request)
}

type Form struct {
	UsePost   bool
	LastError qerror.PublicError
	impl      IFormImplementation
	fields    []field.IField
	fieldsMap map[string]field.IField
}

func New(ctx context.Context, impl IFormImplementation) *Form {
	return &Form{
		impl:   impl,
		fields: impl.GetFields(ctx),
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

	ctx := formctx.WithFormData(r.Context())

	for _, field := range f.fields {
		field.Init(ctx, r.Form)
	}

	if inputValue := r.Form.Get(AntiCSRFInputName); inputValue != "" {
		if err := checkToken(ctx, inputValue); err != nil {
			f.LastError = qerror.ToPublic(err, getAntiCSRFErrorText())
			f.impl.RenderHTML(ctx, w, f)

			return
		}

		f.LastError = nil

		hasFieldsErrors := false
		for _, field := range f.fields {
			if field.Check(ctx) != nil {
				hasFieldsErrors = true
			}
		}

		if !hasFieldsErrors {
			f.LastError = f.impl.OnSave(ctx, w, f)
			if f.LastError == nil {
				f.impl.OnComplete(ctx, w, r)
				return
			}
		}
	}

	f.impl.RenderHTML(ctx, w, f)
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

func (f *Form) GetCaption(ctx context.Context) string {
	return f.impl.GetCaption(ctx)
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

func (f *Form) GetAntiCSRFToken(ctx context.Context) string {
	if config == nil || config.GetTokenSalt() == "" || config.GetCtxParser() == nil {
		panic(errNotInitialized)
	}

	userID, formPath := config.GetCtxParser()(ctx)

	return GenerateToken(userID, formPath)
}
