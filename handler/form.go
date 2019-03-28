package handler

import (
	"context"
	"net/http"

	"github.com/go-qbit/rbac"
	"github.com/go-qbit/web/form"
)

type formHandler struct {
	*Handler
	formImplementation form.IFormImplementation
	caption            string
	perm               *rbac.Permission
}

func (h *formHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	form.New(r.Context(), h.formImplementation).ServeHTTP(w, r)
}

func (h *formHandler) RequiredPermission() (*rbac.Permission) {
	return h.perm
}

func (h *formHandler) Caption(ctx context.Context) string {
	return h.caption
}
