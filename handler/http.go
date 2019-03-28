package handler

import (
	"net/http"

	"github.com/go-qbit/rbac"
)

type httpHandler struct {
	*Handler
	httpHandler http.Handler
	perm *rbac.Permission
}

func (h *httpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.httpHandler.ServeHTTP(w, r)
}

func (h *httpHandler) RequiredPermission() (*rbac.Permission) {
	return h.perm
}
