package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const sessionCookieName = "abcms_session"

// SessionManager handles signing and verifying cookie-based user sessions.
type SessionManager struct {
	secret []byte
	ttl    time.Duration
	secure bool
	now    func() time.Time
}

// NewSessionManager creates a SessionManager using the provided HMAC secret, lifetime, and cookie security flag.
func NewSessionManager(secret []byte, ttl time.Duration, secure bool) *SessionManager {
	secretCopy := make([]byte, len(secret))
	copy(secretCopy, secret)
	return &SessionManager{
		secret: secretCopy,
		ttl:    ttl,
		secure: secure,
		now:    time.Now,
	}
}

// Set writes a signed session cookie for username onto the HTTP response.
func (m *SessionManager) Set(w http.ResponseWriter, username string) {
	expiry := m.now().Add(m.ttl).Unix()
	payloadRaw := username + "|" + strconv.FormatInt(expiry, 10)
	payloadEnc := base64.RawURLEncoding.EncodeToString([]byte(payloadRaw))
	sigEnc := base64.RawURLEncoding.EncodeToString(m.sign(payloadEnc))

	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    payloadEnc + "." + sigEnc,
		Path:     "/",
		MaxAge:   int(m.ttl.Seconds()),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})
}

// Clear expires the session cookie on the HTTP response.
func (m *SessionManager) Clear(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookieName,
		Value:    "",
		Path:     "/",
		MaxAge:   -1,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   m.secure,
	})
}

// GetUsername validates the session cookie on r and returns the authenticated username.
func (m *SessionManager) GetUsername(r *http.Request) (string, bool) {
	cookie, err := r.Cookie(sessionCookieName)
	if err != nil || cookie.Value == "" {
		return "", false
	}

	dotIdx := strings.IndexByte(cookie.Value, '.')
	if dotIdx <= 0 || dotIdx >= len(cookie.Value)-1 {
		return "", false
	}
	payloadEnc := cookie.Value[:dotIdx]
	sigEnc := cookie.Value[dotIdx+1:]

	gotSig, err := base64.RawURLEncoding.DecodeString(sigEnc)
	if err != nil {
		return "", false
	}
	expectedSig := m.sign(payloadEnc)
	if subtle.ConstantTimeCompare(gotSig, expectedSig) != 1 {
		return "", false
	}

	payloadBytes, err := base64.RawURLEncoding.DecodeString(payloadEnc)
	if err != nil {
		return "", false
	}
	payload := string(payloadBytes)
	sepIdx := strings.LastIndexByte(payload, '|')
	if sepIdx <= 0 {
		return "", false
	}

	username := payload[:sepIdx]
	expiryUnix, err := strconv.ParseInt(payload[sepIdx+1:], 10, 64)
	if err != nil {
		return "", false
	}
	if expiryUnix <= m.now().Unix() {
		return "", false
	}

	return username, true
}

func (m *SessionManager) sign(payloadEnc string) []byte {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(payloadEnc))
	return mac.Sum(nil)
}
