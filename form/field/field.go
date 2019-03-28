//go:generate ttgen -package field *.gtt
package field

import (
	"context"
	"io"
	"net/url"

	"github.com/go-qbit/qerror"
)

type IField interface {
	GetName() string
	Init(ctx context.Context, form url.Values)
	GetValue() interface{}
	GetStringValue() string
	Check(ctx context.Context) qerror.PublicError
	ProcessHtml(ctx context.Context, w io.Writer)
	SetError(err qerror.PublicError) qerror.PublicError
}

var ErrMissedReqField = func(ctx context.Context) qerror.PublicError {
	return qerror.PublicErrorf("Field is required")
}
