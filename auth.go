package mbus

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"strings"
	"time"
)

const (
	apiKey               = `Qskvu4Z5JDwGEVswqdAVkiA5B`
	frontendHmacKey      = `ZSqCAFdU7bwxHJUHKYfQUxKin06hMxCK`
	formatRFC7231Fixdate = "Mon, 02 Jan 2006 15:04:05 GMT"
)

func extractAPIPath(u string) string {
	apiIndex := strings.Index(u, "/api/")
	return u[apiIndex:]
}

func formatFixdate(t time.Time) string {
	return t.UTC().Format(formatRFC7231Fixdate)
}

func sha256Hmac(data string, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write([]byte(data))
	hmacBytes := mac.Sum(nil)
	return hex.EncodeToString(hmacBytes)
}

func prepareRequestWithV2Auth(req *http.Request) {
	apiPath := extractAPIPath(req.URL.String())
	formatDate := formatFixdate(time.Now())
	reqHash := sha256Hmac(apiPath+formatDate, frontendHmacKey)

	req.Header.Set("key", apiKey)
	req.Header.Set("X-Date", formatDate)
	req.Header.Set("X-Request-ID", reqHash)
}

func prepareRequestWithV3Auth(req *http.Request) {
	args := req.URL.Query()
	args.Set("key", apiKey)
	req.URL.RawQuery = args.Encode()

	apiPath := extractAPIPath(req.URL.String())
	formatDate := formatFixdate(time.Now())
	reqHash := sha256Hmac(apiPath+formatDate, frontendHmacKey)

	req.Header.Set("X-Date", formatDate)
	req.Header.Set("X-Request-ID", reqHash)
}
