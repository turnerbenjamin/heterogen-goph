package jwt

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

var secret string = os.Getenv("JWT_SECRET")

var header string = encodeBase64([]byte("{\"alg\":\"HS256\",\"typ\":\"JWT\"}"))

type jwtPayload struct {
	Id  string
	Iat int64
	Exp int64
}

// Get JWT token. Payload includes id and expiration
func Sign(id string, exp time.Time) (string, error) {
	if secret == "nil" {
		log.Fatal("jwt secret has not been initialised")
	}

	payload := jwtPayload{
		Id:  id,
		Iat: time.Now().Unix(),
		Exp: exp.Unix(),
	}

	payloadBS, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	encodedPayload := encodeBase64(payloadBS)

	checkSumString := strings.Join([]string{string(header), string(encodedPayload)}, ".")
	signature := generateSignature(checkSumString)

	token := strings.Join([]string{checkSumString, signature}, ".")

	return token, nil
}

// Decode JWT token and return id
func Decode(token string) (string, error) {

	//Split JWT
	jwtElements := strings.Split(token, ".")
	if len(jwtElements) != 3 {
		return "", errors.New("invalid jwt")
	}

	//Validate signature
	tokenHeader := jwtElements[0]
	tokenPayload := jwtElements[1]
	tokenSignature := jwtElements[2]
	checkSum := generateSignature(fmt.Sprintf("%s.%s", tokenHeader, tokenPayload))
	if string(checkSum) != tokenSignature {
		return "", errors.New("invalid jwt")
	}

	//Decode JWT payload
	payloadBS, err := decodeBase64(tokenPayload)
	if err != nil {
		return "", errors.New("invalid jwt")
	}
	decodedBody := jwtPayload{}

	//Parse JSON
	err = json.Unmarshal(payloadBS, &decodedBody)
	if err != nil {
		return "", errors.New("invalid jwt")
	}

	//Validate JWT not expired
	if time.Now().After(time.Unix(decodedBody.Exp, 0)) {
		return "", errors.New("invalid jwt")
	}

	return decodedBody.Id, nil
}

func generateSignature(checkSumString string) string {
	hash := hmac.New(sha256.New, []byte(secret))
	hash.Write([]byte(checkSumString))
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
