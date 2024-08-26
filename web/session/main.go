// Tutors:
//     https://gowebexamples.com/
//     https://www.soberkoder.com/go-rest-api-gorilla-mux/

package main

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/gorilla/sessions"
	"os"
)

const keyLength = 1 << 5

var (
	sessionKey string
	store      *sessions.CookieStore
)

func init() {
	var ok bool
	sessionKey, ok = os.LookupEnv("SESSION_KEY")
	if !ok {
		panic("Environment variable `SESSION_KEY` didn't set.")
	}
	store = sessions.NewCookieStore(resizeSlice([]byte(sessionKey), keyLength))
}

// resizeSlice resize a slice to the specific length.
//
// If slice is less than desired length, pad with zeros.
// If slice is lager than desired length, trim it.
func resizeSlice(key []byte, length int) []byte {
	mylen := len(key)
	if mylen < length {
		padding := make([]byte, length-mylen)
		return append(key, padding...)
	}
	return key[:length]
}

// generateRandomString return a base64 encoded string with an approximate length.
func generateRandomString(length int) (string, error) {
	b := make([]byte, length)
	if _, err := rand.Read(b); err != nil {
		return "", err
	} else {
		return base64.URLEncoding.EncodeToString(b), nil
	}
}

func main() {

}
