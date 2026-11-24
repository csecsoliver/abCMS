package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"golang.org/x/crypto/argon2"
)

func TestSlugify(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"Alice", "alice"},
		{"  Bob Smith! ", "bob-smith"},
		{"user_name.42", "user-name-42"},
		{"---test---", "test"},
		{"", ""},
	}
	for _, tc := range cases {
		if got := Slugify(tc.in); got != tc.want {
			t.Errorf("Slugify(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

func TestBcryptHashAndVerify(t *testing.T) {
	hash, err := HashPassword("s3cr3t-pass")
	if err != nil {
		t.Fatalf("HashPassword error: %v", err)
	}
	if !VerifyPassword("s3cr3t-pass", hash) {
		t.Fatal("VerifyPassword returned false for matching password")
	}
	if VerifyPassword("wrong-pass", hash) {
		t.Fatal("VerifyPassword returned true for wrong password")
	}
	if VerifyPassword("s3cr3t-pass", "") {
		t.Fatal("VerifyPassword returned true for empty hash")
	}
}

func TestVerifyArgon2PHC(t *testing.T) {
	salt := []byte("0123456789abcdef")
	idKey := argon2.IDKey([]byte("legacy-pass"), salt, 2, 16384, 1, 32)
	phcID := fmt.Sprintf(
		"$argon2id$v=19$m=16384,t=2,p=1$%s$%s",
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(idKey),
	)
	if !VerifyPassword("legacy-pass", phcID) {
		t.Fatal("VerifyPassword failed on valid argon2id PHC hash")
	}
	if VerifyPassword("bad-pass", phcID) {
		t.Fatal("VerifyPassword accepted wrong password for argon2id PHC hash")
	}

	iKey := argon2.Key([]byte("legacy-i-pass"), salt, 2, 16384, 1, 32)
	phcI := fmt.Sprintf(
		"$argon2i$v=19$m=16384,t=2,p=1$%s$%s",
		base64.StdEncoding.EncodeToString(salt),
		base64.StdEncoding.EncodeToString(iKey),
	)
	if !VerifyPassword("legacy-i-pass", phcI) {
		t.Fatal("VerifyPassword failed on valid argon2i PHC hash")
	}

	if VerifyPassword("legacy-pass", "$argon2id$v=19$bad$salt$hash") {
		t.Fatal("VerifyPassword accepted malformed argon2 hash")
	}
}

func TestSessionManager(t *testing.T) {
	now := time.Unix(1700000000, 0)
	mgr := NewSessionManager([]byte("test-secret-key-32-bytes-minimum"), time.Hour, false)
	mgr.now = func() time.Time { return now }

	rec := httptest.NewRecorder()
	mgr.Set(rec, "alice")
	cookies := rec.Result().Cookies()
	if len(cookies) != 1 {
		t.Fatalf("expected 1 cookie, got %d", len(cookies))
	}

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req.AddCookie(cookies[0])

	user, ok := mgr.GetUsername(req)
	if !ok || user != "alice" {
		t.Fatalf("GetUsername = (%q, %v), want (\"alice\", true)", user, ok)
	}

	// Tampered cookie should fail verification.
	tampered := *cookies[0]
	tampered.Value += "x"
	reqTampered := httptest.NewRequest(http.MethodGet, "/", nil)
	reqTampered.AddCookie(&tampered)
	if _, ok := mgr.GetUsername(reqTampered); ok {
		t.Fatal("tampered cookie was accepted")
	}

	// Expired cookie should fail verification.
	now = now.Add(2 * time.Hour)
	if _, ok := mgr.GetUsername(req); ok {
		t.Fatal("expired cookie was accepted")
	}

	// Clear should expire the cookie.
	clearRec := httptest.NewRecorder()
	mgr.Clear(clearRec)
	cleared := clearRec.Result().Cookies()
	if len(cleared) != 1 || cleared[0].MaxAge != -1 {
		t.Fatalf("expected MaxAge=-1 on clear cookie, got %+v", cleared)
	}
}
