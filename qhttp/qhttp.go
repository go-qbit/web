package qhttp

import (
	"log"
	"net/http"
	"os"

	"github.com/go-qbit/qerror"
)

type logger interface {
	Print(v ...interface{})
}

const message = "Internal server error"

var Logger logger = log.New(os.Stderr, "", log.LstdFlags)

func Error(w http.ResponseWriter, err error) {
	Logger.Print(err.Error())

	switch err := err.(type) {
	case qerror.PublicError:
		http.Error(w, err.PublicError(), http.StatusInternalServerError)
	default:
		http.Error(w, message, http.StatusInternalServerError)
	}
}

func StdError(w http.ResponseWriter, code int) {
	http.Error(w, http.StatusText(code), code)
}
