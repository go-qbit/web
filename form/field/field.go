//go:generate ttgen -package field *.gtt
package field

import (
	"context"
	"io"
	"net/url"

	"github.com/go-qbit/qerror"
)

type IField interface {
	Process(context.Context, url.Values) qerror.PublicError
	GetName() string
	GetValue() interface{}
	GetStringValue() string
	ProcessField(context.Context, io.Writer)
	SetError(qerror.PublicError) qerror.PublicError
}

type Field struct {
	Name      string
	Caption   string
	Required  bool
	LastError qerror.PublicError
}

func (f *Field) GetName() string {
	return f.Name
}

func (f *Field) SetError(err qerror.PublicError) qerror.PublicError {
	f.LastError = err

	return f.LastError
}
