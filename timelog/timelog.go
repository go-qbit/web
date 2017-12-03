package timelog

import (
	"bytes"
	"net/http"
	"regexp"
	"strings"

	"github.com/go-qbit/timelog"
)

type tlHandler struct {
	handler http.Handler
}

type tlWriter struct {
	*bytes.Buffer
	httpWriter http.ResponseWriter
}

func (w *tlWriter) Header() http.Header {
	return w.httpWriter.Header()
}

func (w *tlWriter) WriteHeader(code int) {
	w.httpWriter.WriteHeader(code)
}

func HtmlInjection(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tlW := &tlWriter{
			Buffer:     &bytes.Buffer{},
			httpWriter: w,
		}

		ctx := timelog.Start(r.Context(), "ServeHTTP")
		handler.ServeHTTP(tlW, r.WithContext(ctx))
		timelog.Finish(ctx)

		htmlData := tlW.String()

		tlText := "<hr><pre>" + timelog.Get(ctx).Analyze().String() + "</pre>"

		if strings.Contains(w.Header().Get("Content-Type"), "text/html") {
			replaced := false
			htmlData = regexp.MustCompile("(?i:</body)").ReplaceAllStringFunc(htmlData, func(s string) string {
				replaced = true
				return tlText + s
			})

			if !replaced {
				htmlData += tlText
			}
		}

		w.Write([]byte(htmlData))
	})
}
