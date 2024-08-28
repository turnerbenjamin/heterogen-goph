package cookies

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type CookieExpiration int

const Minute CookieExpiration = 60
const Hour CookieExpiration = 60 * Minute
const Day CookieExpiration = Hour * 24

var secret string = os.Getenv("COOKIE_SECRET")
var authCookieName string = "hg_auth"

func NewAuthCookie(userId string, maxAge CookieExpiration) *http.Cookie {
	if secret == "nil" {
		log.Fatal("Cookie secret has not been initialised")
	}

	encodedCreatedAt := encodeBase64([]byte(time.Now().String()))
	encodedAlgoritm := encodeBase64([]byte("HS256"))
	encodedId := encodeBase64([]byte(userId))
	signature := generateSignature(encodeBase64([]byte(userId)))

	payload := strings.Join([]string{encodedCreatedAt, encodedAlgoritm, encodedId, signature}, ".")

	return &http.Cookie{
		Name:        authCookieName,
		Value:       payload,
		RawExpires:  "MaxAge>0",
		MaxAge:      int(maxAge),
		Secure:      true,
		HttpOnly:    true,
		SameSite:    http.SameSiteLaxMode,
		Partitioned: true,
	}
}

func UnsetAuthCookie() *http.Cookie {
	return &http.Cookie{
		Name:        authCookieName,
		RawExpires:  "MaxAge>0",
		MaxAge:      0,
		Secure:      true,
		HttpOnly:    true,
		SameSite:    http.SameSiteLaxMode,
		Partitioned: true,
	}
}

func ParseAuthCookie(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(authCookieName)
	if err != nil {
		return "", false
	}

	return decodeAuthCookie(cookie)
}

func decodeAuthCookie(cookie *http.Cookie) (string, bool) {

	err := cookie.Valid()
	if err != nil {
		return "", false
	}

	cookieElements := strings.Split(cookie.Value, ".")
	if len(cookieElements) != 4 {
		return "", false
	}

	encodedId := cookieElements[2]
	signature := cookieElements[3]
	if string(signature) != generateSignature(encodedId) {
		return "", false
	}

	userId, err := decodeBase64(encodedId)
	if err != nil {
		return "", false
	}

	return string(userId), true
}

func generateSignature(value string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(value))
	return encodeBase64(hash.Sum(nil))
}

func encodeBase64(src []byte) string {
	return base64.RawURLEncoding.EncodeToString(src)
}

func decodeBase64(src string) ([]byte, error) {
	bs, err := base64.RawURLEncoding.DecodeString(src)
	if err != nil {
		return []byte{}, err
	}
	return bs, nil
}
