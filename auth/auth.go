package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-qbit/qerror"
)

type Auth struct {
	salt       string
	secure     bool
	cookieName string
}

func New(salt string, secure bool) *Auth {
	return &Auth{salt, secure, "a"}
}

func (a *Auth) GetSessionId(id uint32) string {
	strId := strconv.FormatUint(uint64(id), 16)
	strTime := strconv.FormatUint(uint64(time.Now().Unix()), 16)

	return strings.Join([]string{strId, strTime, a.getSessionSign(strId, strTime)}, "-")
}

func (a *Auth) GetUserId(sessionId string) (uint32, error) {
	splitted := strings.Split(sessionId, "-")
	if len(splitted) != 3 || splitted[2] != a.getSessionSign(splitted[0], splitted[1]) {
		return 0, qerror.Errorf("invalid session")
	}

	id, err := strconv.ParseUint(splitted[0], 16, 32)
	if err != nil {
		return 0, qerror.Errorf("cannot parse id: %s", err.Error())
	}

	return uint32(id), nil
}

func (a *Auth) SetToHttpWriter(userId uint32, w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    a.cookieName,
		Value:   a.GetSessionId(userId),
		Path:    "/",
		Expires: time.Now().Add(365 * 24 * time.Hour),
		Secure:  a.secure,
	})
}

func (a *Auth) DeleteFromHttpWriter(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:    a.cookieName,
		Value:   "",
		Path:    "/",
		Expires: time.Now().Add(-10 * 365 * 24 * time.Hour),
		Secure:  a.secure,
	})
}

func (a *Auth) GetUserIdFromHttpRequest(r *http.Request) (uint32, error) {
	sessionCookie, err := r.Cookie(a.cookieName)
	if err != nil {
		return 0, nil
	}

	return a.GetUserId(sessionCookie.Value)
}

func (a *Auth) getSessionSign(id, time string) string {
	salt := sha256.Sum256(append([]byte(time), []byte(a.salt)...))
	sign := sha256.Sum256(append([]byte(id), salt[:]...))

	return hex.EncodeToString(sign[:])
}
