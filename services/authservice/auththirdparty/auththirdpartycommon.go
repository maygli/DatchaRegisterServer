package auththirdparty

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"golang.org/x/oauth2"
)

const (
	STATE_LENGTH    = 24
	VERIFIER_LENGHT = 24
	// Life time of state cookie in seconds
	COOKIE_STATE_LIFE_TIME = 600
	STATE_COOKIE_NAME      = "oauthstate"
	VERIFIER_COOKIE_NAME   = "verifier"
	STATE_NAME             = "state"
	CODE_NAME              = "code"
	AUTH_HEADER            = "Authorization"
)

func GenerateRandonString(size int) string {
	b := make([]byte, size)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	return state
}

func GenerateStateCockie(state string) *http.Cookie {
	expiration := time.Now().Add(COOKIE_STATE_LIFE_TIME * time.Second)
	cookie := http.Cookie{
		Name:     STATE_COOKIE_NAME,
		Value:    state,
		Expires:  expiration,
		HttpOnly: true,
	}
	return &cookie
}

func GenerateVerifierCockie(verifier string) *http.Cookie {
	expiration := time.Now().Add(COOKIE_STATE_LIFE_TIME * time.Second)
	cookie := http.Cookie{
		Name:     VERIFIER_COOKIE_NAME,
		Value:    verifier,
		Expires:  expiration,
		HttpOnly: true,
	}
	return &cookie
}

func ProcessLogin(w http.ResponseWriter, r *http.Request, config oauth2.Config) {
	state := GenerateRandonString(STATE_LENGTH)
	url := config.AuthCodeURL(state)
	http.SetCookie(w, GenerateStateCockie(state))
	http.Redirect(w, r, url, http.StatusSeeOther)
}

func VerifyCookieValue(r *http.Request, coockieName string, expValue string) error {
	cookieValue, err := r.Cookie(coockieName)
	if err != nil {
		return err
	}
	if cookieValue.Value != expValue {
		return errors.New("invalid cookie value")
	}
	return nil
}

func GetToken(r *http.Request, config oauth2.Config) (*oauth2.Token, error) {
	rsvState := r.FormValue(STATE_NAME)
	if rsvState == "" {
		return nil, errors.New("empty state received")
	}
	err := VerifyCookieValue(r, STATE_COOKIE_NAME, rsvState)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Incorrect state value. Error: %s", err.Error()))
	}
	code := r.FormValue(CODE_NAME)
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("can't get token. Error: %s", err.Error()))
	}
	return token, nil
}
