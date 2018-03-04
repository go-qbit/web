//go:generate ttgen -package field *.gtt
package field

import (
	"context"
	"io"
	"net/url"

	"github.com/go-qbit/qerror"
)

type IField interface {
	Process(context.Context, url.Values) error
	GetName() string
	GetValue() interface{}
	GetStringValue() string
	ProcessField(context.Context, io.Writer)
	SetError(error)
}

type Field struct {
	Name      string
	Caption   string
	Required  bool
	LastError error
}

var ErrMissedReqField = qerror.Errorf("Обязательное поле")

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) SetError(err error) {
	f.LastError = err
}
