package handler

import (
	"context"
	"fmt"
	"net/http"
	"path"
	"strings"

	"github.com/go-qbit/rbac"
	"github.com/go-qbit/web/form"
)

type IHandler interface {
	Caption(ctx context.Context) string
	SubHandlers() []string
	SubHandler(path string) IHandler
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Parent() IHandler
	SetParent(parent IHandler, path string)
	GetPath() string
	GetFullPath() string
	FirstAllowedSubHandler(ctx context.Context) IHandler
	RequiredPermission() *rbac.Permission
}

type Handler struct {
	children    []string
	childrenMap map[string]IHandler
	path        string
	parent      IHandler
}

type handlerCtxT int8

var handlerCtx handlerCtxT = 0

func (h *Handler) Caption(ctx context.Context) string {
	return ""
}

func (h *Handler) RequiredPermission() (perm *rbac.Permission) {
	return
}

func (h *Handler) SubHandlers() []string {
	return h.children
}

func (h *Handler) SubHandler(path string) IHandler {
	return h.childrenMap[path]
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.ServeSubHandlers(w, r)
}

func (h *Handler) Parent() IHandler {
	return h.parent
}

func (h *Handler) SetParent(parent IHandler, path string) {
	h.parent = parent
	h.path = path
}

func (h *Handler) GetPath() string {
	return h.path
}

func (h *Handler) GetFullPath() string {
	if h.parent == nil {
		return "/"
	}

	fullPath := h.parent.GetFullPath()
	if fullPath == "/" {
		return fullPath + h.GetPath()
	}

	return fullPath + "/" + h.GetPath()
}

func (h *Handler) ServeSubHandlers(w http.ResponseWriter, r *http.Request) {
	fullPath := h.GetFullPath()

	urlPath := cleanPath(r.URL.Path)

	var relativePath string
	if len(fullPath) > len(urlPath) && strings.HasPrefix(fullPath, urlPath) {
		relativePath = ""
	} else {
		relativePath = strings.TrimPrefix(urlPath, fullPath)
		relativePath = strings.TrimPrefix(relativePath, "/")
	}

	subPath := strings.Split(relativePath, "/")[0]

	subHandler := h.SubHandler(subPath)

	if subPath == "" && subHandler == nil {
		subHandler = h.FirstAllowedSubHandler(r.Context())
	}

	if subHandler == nil {
		http.NotFound(w, r)
		return
	}

	if !rbac.HasPermission(r.Context(), subHandler.RequiredPermission()) {
		http.NotFound(w, r)
		return
	}

	subHandler.ServeHTTP(w, r.WithContext(
		context.WithValue(r.Context(), handlerCtx, subHandler),
	))
}

func (h *Handler) Handle(path string, handler IHandler) {
	if h.childrenMap == nil {
		h.childrenMap = make(map[string]IHandler)
	}

	if h.childrenMap[path] != nil {
		panic(fmt.Sprintf("%s already exists", path))
	}

	h.children = append(h.children, path)
	h.childrenMap[path] = handler
	handler.SetParent(h, path)
}

func (h *Handler) HandleHTTP(path string, handler http.Handler, permission *rbac.Permission) {
	h.Handle(path, &httpHandler{
		&Handler{},
		handler,
		permission,
	})
}

func (h *Handler) HandleForm(path string, implementation form.IFormImplementation, caption string, permission *rbac.Permission) {
	h.Handle(path, &formHandler{
		&Handler{},
		implementation,
		caption,
		permission,
	})
}

func (h *Handler) FirstAllowedSubHandler(ctx context.Context) IHandler {
	for _, p := range h.children {
		if handler := h.childrenMap[p]; handler != nil {
			subHandlers := handler.SubHandlers()
			if (len(subHandlers) == 0 || handler.FirstAllowedSubHandler(ctx) != nil) &&
				rbac.HasPermission(ctx, handler.RequiredPermission()) {
				return handler
			}
		}
	}
	return nil
}

func GetCurHandler(ctx context.Context) IHandler {
	ctxData := ctx.Value(handlerCtx)
	if ctxData == nil {
		return nil
	}

	return ctxData.(IHandler)
}

func cleanPath(p string) string {
	if p == "" {
		return "/"
	}
	if p[0] != '/' {
		p = "/" + p
	}
	np := path.Clean(p)

	if p[len(p)-1] == '/' && np != "/" {
		np += "/"
	}

	return np
}
