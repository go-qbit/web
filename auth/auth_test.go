package auth_test

import (
	"log"
	"testing"

	"github.com/go-qbit/web/auth"
)

func TestAuth_GetUserId(t *testing.T) {
	a := auth.New("salt", true)

	session := a.GetSessionId(100500)
	userId, err := a.GetUserId(session)
	if err != nil {
		log.Fatal(err)
	}
	if userId != 100500 {
		t.Fatalf("User Ids are not equal")
	}
}
