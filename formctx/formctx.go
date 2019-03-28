package formctx

import "context"

type formCtxKeyT struct{}

var formCtxKey formCtxKeyT

type data struct {
	storage map[string]interface{}
}

func WithFormData(ctx context.Context) context.Context {
	return context.WithValue(ctx, formCtxKey, &data{
		storage: map[string]interface{}{},
	})
}

func SetValue(ctx context.Context, name string, value interface{}) {
	getData(ctx).storage[name] = value
}

func GetValue(ctx context.Context, name string) interface{} {
	return getData(ctx).storage[name]
}

func getData(ctx context.Context) *data {
	v := ctx.Value(formCtxKey)
	if v == nil {
		panic("form data wasn't initialized")
	}

	return v.(*data)
}
