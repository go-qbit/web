package qhttp

import (
	"net/http"
	"strings"
)

type ServeMux struct {
	*http.ServeMux
}

func NewServeMux() *ServeMux {
	return &ServeMux{
		http.NewServeMux(),
	}
}

func (m *ServeMux) HandleRootWithAlias(alias string, h http.Handler) {
	m.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, alias) {
			h.ServeHTTP(w, r)
		} else {
			http.NotFound(w, r)
		}
	})
}
