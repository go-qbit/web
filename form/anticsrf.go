package form

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"time"
)

const AntiCSRFInputName = "anti_csrf_token"

const (
	lifetime   = time.Hour * 1
	timeFormat = "2006010215" // year+month+day+hour
)

var (
	salt      string
	length    int
	ctxParser func(ctx context.Context) (userID uint32, formPath string) = nil

	errNotInitialized = errors.New("Please call InitAntiCSRF() first to get your web forms protected and working well")
	errInvalidToken   = errors.New("Invalid anti-CSRF token")
	errEmptyParams    = errors.New("Empty init params must be not empty")

	ErrTextInvalidToken = "Ключ формы устарел, попробуйте обновить страницу"
)

func InitAntiCSRF(tokenSalt string, tokenLength int, f func(ctx context.Context) (userID uint32, formPath string)) error {
	if tokenSalt == "" || tokenLength == 0 || f == nil {
		return errEmptyParams
	}
	if salt != "" {
		log.Println("AntiCSRF seems already initialized and you called Init() twice")
		return nil
	}

	salt = tokenSalt
	length = tokenLength
	ctxParser = f

	return nil
}

func GenerateToken(userID uint32, formPath string, prevHour bool) string {
	dt := time.Now().UTC()
	if prevHour {
		dt = dt.Add(-1 * lifetime)
	}

	row := fmt.Sprintf("%s%s%s%d", salt, dt.Format(timeFormat), formPath, userID)

	hasher := md5.New()
	hasher.Write([]byte(row))

	return hex.EncodeToString(hasher.Sum(nil)[:length])
}

func checkToken(ctx context.Context, token string) error {
	if salt == "" || ctxParser == nil {
		return errNotInitialized
	}
	if token == "" {
		return errInvalidToken
	}

	userID, formPath := ctxParser(ctx)
	t := GenerateToken(userID, formPath, false)
	if t == token {
		return nil
	}
	t = GenerateToken(userID, formPath, true)
	if t == token {
		return nil
	}

	return errInvalidToken
}
