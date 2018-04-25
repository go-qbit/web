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
	config IAntiCSRFConfig

	errNotInitialized = errors.New("Please call InitAntiCSRF() first to get your web forms protected and working well")
	errInvalidToken   = errors.New("Invalid anti-CSRF token")
	errEmptyParams    = errors.New("Empty init params must be not empty")
)

func InitAntiCSRF(c IAntiCSRFConfig) error {
	if config != nil && c.GetTokenSalt() != "" {
		log.Println("AntiCSRF seems already initialized and you called InitAntiCSRF() twice")
		return nil
	}
	if c == nil || c.GetTokenSalt() == "" || c.GetTokenLength() == 0 || c.GetCtxParser() == nil {
		return errEmptyParams
	}

	config = c

	return nil
}

func GenerateToken(userID uint32, formPath string, prevHour bool) string {
	if config == nil || config.GetTokenSalt() == "" || config.GetCtxParser() == nil {
		return ""
	}

	dt := time.Now().UTC()
	if prevHour {
		dt = dt.Add(-1 * lifetime)
	}

	row := fmt.Sprintf("%s%s%s%d", config.GetTokenSalt(), dt.Format(timeFormat), formPath, userID)

	hasher := md5.New()
	hasher.Write([]byte(row))

	return hex.EncodeToString(hasher.Sum(nil)[:config.GetTokenLength()])
}

func checkToken(ctx context.Context, token string) error {
	if config == nil || config.GetTokenSalt() == "" || config.GetCtxParser() == nil {
		return errNotInitialized
	}
	if token == "" {
		return errInvalidToken
	}

	userID, formPath := config.GetCtxParser()(ctx)
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

func getAntiCSRFErrorText() string {
	if config == nil {
		return "Form token is invalid"
	}

	return config.GetErrorText()
}
