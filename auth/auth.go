package auth

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type authContextKey string

var authUserKey = authContextKey("user")

type IUserInfo interface {
	GetUserById(context.Context, string) (interface{}, error)
}

type Auth struct {
	user IUserInfo
	salt string
}

func New(user IUserInfo, salt string) *Auth {
	return &Auth{user, salt}
}

func (a *Auth) GetSessionId(id uint32) string {
	strId := strconv.FormatUint(uint64(id), 16)
	strTime := strconv.FormatUint(uint64(time.Now().Unix()), 16)

	return strings.Join([]string{strId, strTime, a.GetSessionSign(strId, strTime)}, "-")
}

func (a *Auth) GetSessionSign(id, time string) string {
	salt := sha256.Sum256(append([]byte(time), []byte(a.salt)...))
	sign := sha256.Sum256(append([]byte(id), salt[:]...))

	return hex.EncodeToString(sign[:])
}

func (a *Auth) ToCtx(ctx context.Context, r *http.Request) context.Context {
	sessionCookie, err := r.Cookie("a")
	if err != nil {
		return ctx
	}

	splitted := strings.Split(sessionCookie.Value, "-")
	if len(splitted) != 3 || splitted[2] != a.GetSessionSign(splitted[0], splitted[1]) {
		return ctx
	}

	user, err := a.user.GetUserById(ctx, splitted[0])
	if err != nil {
		return ctx
	}

	return context.WithValue(ctx, authUserKey, user)
}

func (a *Auth) Handler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		handler.ServeHTTP(w, r.WithContext(a.ToCtx(r.Context(), r)))
	})
}

func UserFromCtx(ctx context.Context) interface{} {
	return ctx.Value(authUserKey)
}
