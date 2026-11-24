package handler

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"golang.org/x/crypto/argon2"

	"github.com/csecsoliver/abCMS/internal/auth"
	"github.com/csecsoliver/abCMS/internal/config"
	"github.com/csecsoliver/abCMS/internal/db"
	"github.com/csecsoliver/abCMS/internal/media"
)

func setupTestServer(t *testing.T) (*Handler, *db.Store, *httptest.Server) {
	t.Helper()
	tmpDir := t.TempDir()
	staticDir := filepath.Join(tmpDir, "static")
	uploadDir := filepath.Join(staticDir, "uploads")
	dbPath := filepath.Join(tmpDir, "test.sqlite")

	cfg := config.Config{
		Addr:           ":0",
		DBPath:         dbPath,
		StaticDir:      staticDir,
		UploadDir:      uploadDir,
		SessionSecret:  []byte("0123456789abcdef0123456789abcdef"),
		SessionTTL:     24 * time.Hour,
		MaxUploadBytes: 15 << 20,
	}

	store, err := db.Open(dbPath)
	if err != nil {
		t.Fatalf("db.Open: %v", err)
	}
	t.Cleanup(func() { _ = store.Close() })

	sessions := auth.NewSessionManager(cfg.SessionSecret, cfg.SessionTTL, false)
	mediaSvc := media.NewService(uploadDir)
	h := New(cfg, store, sessions, mediaSvc)

	srv := httptest.NewServer(h.Routes())
	t.Cleanup(srv.Close)
	t.Cleanup(h.Close)
	return h, store, srv
}

func noRedirectClient(cookies []*http.Cookie) *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
}

func makePNGBytes(t *testing.T) []byte {
	t.Helper()
	img := image.NewRGBA(image.Rect(0, 0, 16, 16))
	img.Set(0, 0, color.RGBA{R: 255, A: 255})
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		t.Fatalf("encode png: %v", err)
	}
	return buf.Bytes()
}

func TestEndToEndFlow(t *testing.T) {
	_, store, srv := setupTestServer(t)
	client := noRedirectClient(nil)

	// 1. Homepage renders 200.
	resp, err := client.Get(srv.URL + "/")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET / = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()

	// 2. 404 renders NotFoundPage.
	resp, err = client.Get(srv.URL + "/does-not-exist")
	if err != nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("GET /does-not-exist = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()

	// 3. Unauthenticated /dashboard redirects to /login.
	resp, err = client.Get(srv.URL + "/dashboard")
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("GET /dashboard = (%v, %v)", resp, err)
	}
	if loc := resp.Header.Get("Location"); !strings.HasPrefix(loc, "/login") {
		t.Fatalf("expected redirect to /login, got %q", loc)
	}
	_ = resp.Body.Close()

	// 4. Sign up Alice.
	form := url.Values{
		"option":   {"signup"},
		"username": {"Alice"},
		"password": {"password123"},
	}
	resp, err = client.PostForm(srv.URL+"/login", form)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("POST /login signup = (%v, %v)", resp, err)
	}
	aliceCookies := resp.Cookies()
	_ = resp.Body.Close()
	if len(aliceCookies) == 0 {
		t.Fatal("expected session cookie after signup")
	}

	// Duplicate signup fails with 409.
	resp, err = client.PostForm(srv.URL+"/login", form)
	if err != nil || resp.StatusCode != http.StatusConflict {
		t.Fatalf("duplicate signup status = %d, want 409", resp.StatusCode)
	}
	_ = resp.Body.Close()

	// 5. Create text-only post as Alice.
	var bodyBuf bytes.Buffer
	mw := multipart.NewWriter(&bodyBuf)
	_ = mw.WriteField("title", "First Text Post")
	_ = mw.WriteField("content", "Hello **Markdown**!")
	_ = mw.Close()

	req, _ := http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/add", &bodyBuf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	for _, c := range aliceCookies {
		req.AddCookie(c)
	}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("create text post = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()

	// 6. View text-only post (verifying no nil/missing image crash).
	posts, err := store.ListPosts(context.Background())
	if err != nil || len(posts) != 1 {
		t.Fatalf("ListPosts = (%d, %v)", len(posts), err)
	}
	postID := posts[0].ID

	resp, err = client.Get(srv.URL + "/posts/1")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET /posts/1 = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()

	// 7. Create image post as Alice.
	bodyBuf.Reset()
	mw = multipart.NewWriter(&bodyBuf)
	_ = mw.WriteField("title", "Image Post")
	_ = mw.WriteField("content", "With PNG")
	part, _ := mw.CreateFormFile("image", "sample.png")
	_, _ = part.Write(makePNGBytes(t))
	_ = mw.Close()

	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/add", &bodyBuf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	for _, c := range aliceCookies {
		req.AddCookie(c)
	}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("create image post = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()

	// 8. Edit Alice's first post via /formapi/posts/mod/{post_id}.
	bodyBuf.Reset()
	mw = multipart.NewWriter(&bodyBuf)
	_ = mw.WriteField("title", "Edited Text Post")
	_ = mw.WriteField("content", "Updated content")
	_ = mw.Close()

	req, _ = http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/mod/1", &bodyBuf)
	req.Header.Set("Content-Type", mw.FormDataContentType())
	for _, c := range aliceCookies {
		req.AddCookie(c)
	}
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("edit post = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()

	updated, err := store.FindPostByID(context.Background(), postID)
	if err != nil || updated.Title != "Edited Text Post" {
		t.Fatalf("expected updated title, got %+v (%v)", updated, err)
	}

	// 9. Cross-site POST is blocked by CrossOriginProtection CSRF middleware.
	csrfReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/delete/1", nil)
	csrfReq.Header.Set("Origin", "https://evil.example")
	csrfReq.Header.Set("Sec-Fetch-Site", "cross-site")
	for _, c := range aliceCookies {
		csrfReq.AddCookie(c)
	}
	resp, err = client.Do(csrfReq)
	if err != nil || resp.StatusCode != http.StatusForbidden {
		t.Fatalf("expected 403 for cross-site POST, got %d (%v)", resp.StatusCode, err)
	}
	_ = resp.Body.Close()

	// 10. Sign up Bob and verify he cannot delete Alice's post.
	bobForm := url.Values{
		"option":   {"signup"},
		"username": {"Bob"},
		"password": {"bobpass"},
	}
	resp, err = client.PostForm(srv.URL+"/login", bobForm)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("signup Bob = (%v, %v)", resp, err)
	}
	bobCookies := resp.Cookies()
	_ = resp.Body.Close()

	delAsBob, _ := http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/delete/1", nil)
	for _, c := range bobCookies {
		delAsBob.AddCookie(c)
	}
	resp, err = client.Do(delAsBob)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Bob delete attempt = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()
	if _, err := store.FindPostByID(context.Background(), postID); err != nil {
		t.Fatalf("Bob should not have deleted Alice's post, err=%v", err)
	}

	// 11. Alice deletes her post.
	delAsAlice, _ := http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/delete/1", nil)
	for _, c := range aliceCookies {
		delAsAlice.AddCookie(c)
	}
	resp, err = client.Do(delAsAlice)
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("Alice delete = (%v, %v)", resp, err)
	}
	_ = resp.Body.Close()
	if _, err := store.FindPostByID(context.Background(), postID); err == nil {
		t.Fatal("Alice's post should have been deleted")
	}
}

func TestLegacyArgon2UserLogin(t *testing.T) {
	_, store, srv := setupTestServer(t)
	client := noRedirectClient(nil)

	salt := []byte("0123456789abcdef")
	key := argon2.IDKey([]byte("argonpass"), salt, 1, 8192, 1, 32)
	// Base64-encode without padding for PHC string.
	phc := "$argon2id$v=19$m=8192,t=1,p=1$" +
		encodeRawB64(salt) + "$" + encodeRawB64(key)

	if _, err := store.CreateUser(context.Background(), "legacy", phc); err != nil {
		t.Fatalf("CreateUser legacy: %v", err)
	}

	resp, err := client.PostForm(srv.URL+"/login", url.Values{
		"option":   {"login"},
		"username": {"Legacy"},
		"password": {"argonpass"},
	})
	if err != nil || resp.StatusCode != http.StatusSeeOther {
		t.Fatalf("legacy login = (%d, %v)", resp.StatusCode, err)
	}
	_ = resp.Body.Close()
}

func encodeRawB64(b []byte) string {
	const table = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var out strings.Builder
	for i := 0; i < len(b); i += 3 {
		var val uint32
		rem := len(b) - i
		val |= uint32(b[i]) << 16
		if rem > 1 {
			val |= uint32(b[i+1]) << 8
		}
		if rem > 2 {
			val |= uint32(b[i+2])
		}
		out.WriteByte(table[(val>>18)&0x3F])
		out.WriteByte(table[(val>>12)&0x3F])
		if rem > 1 {
			out.WriteByte(table[(val>>6)&0x3F])
		}
		if rem > 2 {
			out.WriteByte(table[val&0x3F])
		}
	}
	return out.String()
}

func readBody(t *testing.T, resp *http.Response) string {
	t.Helper()
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read body: %v", err)
	}
	return string(b)
}

func TestHTMXSPANavigationAndHistoryRestore(t *testing.T) {
	_, _, srv := setupTestServer(t)
	client := noRedirectClient(nil)

	// 1. Standard GET / renders full HTMX 4 SPA shell with explicit :inherited boost, history element, and SSE bridge.
	resp, err := client.Get(srv.URL + "/")
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("GET / = (%v, %v)", resp, err)
	}
	fullHTML := readBody(t, resp)
	for _, want := range []string{
		"<!doctype html>",
		`hx-boost:inherited="true"`,
		`hx-target:inherited="#app-main"`,
		`hx-swap:inherited="innerHTML show:top showTarget:window"`,
		`hx-sse:connect="/events"`,
		`hx-history-elt`,
		`id="reduced-motion-toggle"`,
		`role="switch"`,
		`/static/js/hx-sse.min.js`,
		`/static/js/hx-preload.min.js`,
	} {
		if !strings.Contains(fullHTML, want) {
			t.Errorf("full GET / missing %q", want)
		}
	}
	vary := strings.Join(resp.Header.Values("Vary"), ",")
	if !strings.Contains(vary, "HX-Request") || !strings.Contains(vary, "HX-Target") {
		t.Errorf("expected Vary header to include HX-Request and HX-Target, got %q", vary)
	}

	// 2. HTMX 4 GET / returns partial fragment with <template hx type="partial"> nav/stats and title.
	req, _ := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	req.Header.Set("HX-Request", "true")
	req.Header.Set("HX-Target", "main#app-main")
	resp, err = client.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("HTMX GET / = (%v, %v)", resp, err)
	}
	frag := readBody(t, resp)
	if strings.Contains(frag, "<!doctype html>") {
		t.Error("HTMX partial response should not contain <!doctype html>")
	}
	for _, want := range []string{
		"<title>abCMS</title>",
		`<template hx type="partial" hx-target="#site-nav" hx-swap="outerHTML">`,
		`<template hx type="partial" hx-target="#header-stats" hx-swap="innerHTML">`,
		`id="post-grid-region"`,
	} {
		if !strings.Contains(frag, want) {
			t.Errorf("HTMX partial GET / missing %q", want)
		}
	}

	// 3. History restore request returns full HTML layout so htmx can extract [hx-history-elt].
	histReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/", nil)
	histReq.Header.Set("HX-Request", "true")
	histReq.Header.Set("HX-History-Restore-Request", "true")
	resp, err = client.Do(histReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("history restore GET / = (%v, %v)", resp, err)
	}
	histHTML := readBody(t, resp)
	if !strings.Contains(histHTML, "<!doctype html>") || !strings.Contains(histHTML, "hx-history-elt") {
		t.Error("history restore request must return full layout containing hx-history-elt")
	}

	// 4. HTMX 404 returns fragment for response-targets swap into #app-main.
	nfReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/nope", nil)
	nfReq.Header.Set("HX-Request", "true")
	resp, err = client.Do(nfReq)
	if err != nil || resp.StatusCode != http.StatusNotFound {
		t.Fatalf("HTMX 404 = (%v, %v)", resp, err)
	}
	nfBody := readBody(t, resp)
	if strings.Contains(nfBody, "<!doctype html>") || !strings.Contains(nfBody, "Page Not Found") {
		t.Errorf("unexpected HTMX 404 body: %s", nfBody)
	}
}

func TestHTMXAuthDashboardSearchPreviewAndDelete(t *testing.T) {
	_, store, srv := setupTestServer(t)
	client := noRedirectClient(nil)

	// 1. Unauthenticated HTMX 4 dashboard tab request (tag#id HX-Target) retargets #app-main to /login.
	unauthReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/dashboard/posts", nil)
	unauthReq.Header.Set("HX-Request", "true")
	unauthReq.Header.Set("HX-Target", "div#dashboard-content")
	resp, err := client.Do(unauthReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("unauth HTMX /dashboard/posts = (%v, %v)", resp, err)
	}
	if got := resp.Header.Get("HX-Retarget"); got != "#app-main" {
		t.Errorf("expected HX-Retarget=#app-main, got %q", got)
	}
	if got := resp.Header.Get("HX-Push-Url"); !strings.HasPrefix(got, "/login?return_to=") {
		t.Errorf("expected HX-Push-Url=/login?return_to=..., got %q", got)
	}
	_ = readBody(t, resp)

	// 2. HTMX signup returns 200, pushes /dashboard, and OOB-updates nav & toast.
	signupReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/login", strings.NewReader(url.Values{
		"option":   {"signup"},
		"username": {"htmxuser"},
		"password": {"secret123"},
	}.Encode()))
	signupReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	signupReq.Header.Set("HX-Request", "true")
	resp, err = client.Do(signupReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("HTMX signup = (%v, %v)", resp, err)
	}
	cookies := resp.Cookies()
	if got := resp.Header.Get("HX-Push-Url"); got != "/dashboard" {
		t.Errorf("expected HX-Push-Url=/dashboard, got %q", got)
	}
	body := readBody(t, resp)
	if !strings.Contains(body, "Logout (htmxuser)") || !strings.Contains(body, "Welcome to your Dashboard") {
		t.Errorf("unexpected HTMX signup response: %s", body)
	}

	// 3. HTMX duplicate signup returns 409 inline error banner (for response-targets).
	dupReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/login", strings.NewReader(url.Values{
		"option":   {"signup"},
		"username": {"htmxuser"},
		"password": {"secret123"},
	}.Encode()))
	dupReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	dupReq.Header.Set("HX-Request", "true")
	resp, err = client.Do(dupReq)
	if err != nil || resp.StatusCode != http.StatusConflict {
		t.Fatalf("HTMX duplicate signup status = (%v, %v)", resp, err)
	}
	errFrag := readBody(t, resp)
	if strings.Contains(errFrag, "<!doctype html>") || !strings.Contains(errFrag, "Username already exists") {
		t.Errorf("expected inline error banner, got: %s", errFrag)
	}

	// 4. Empty HTMX post creation returns 422 inline error banner.
	var emptyBuf bytes.Buffer
	mw := multipart.NewWriter(&emptyBuf)
	_ = mw.WriteField("title", "")
	_ = mw.WriteField("content", "")
	_ = mw.Close()
	emptyReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/add", &emptyBuf)
	emptyReq.Header.Set("Content-Type", mw.FormDataContentType())
	emptyReq.Header.Set("HX-Request", "true")
	for _, c := range cookies {
		emptyReq.AddCookie(c)
	}
	resp, err = client.Do(emptyReq)
	if err != nil || resp.StatusCode != http.StatusUnprocessableEntity {
		t.Fatalf("empty post create = (%v, %v)", resp, err)
	}
	_ = readBody(t, resp)

	// 6. Create two posts via HTMX.
	for _, p := range []struct{ title, content string }{
		{"Alpha HTMX Post", "Live search target alpha"},
		{"Beta Notes", "Something else entirely"},
	} {
		var buf bytes.Buffer
		w := multipart.NewWriter(&buf)
		_ = w.WriteField("title", p.title)
		_ = w.WriteField("content", p.content)
		_ = w.Close()
		req, _ := http.NewRequest(http.MethodPost, srv.URL+"/formapi/posts/add", &buf)
		req.Header.Set("Content-Type", w.FormDataContentType())
		req.Header.Set("HX-Request", "true")
		req.Header.Set("HX-Target", "div#dashboard-content")
		for _, c := range cookies {
			req.AddCookie(c)
		}
		resp, err = client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			t.Fatalf("create %q = (%v, %v)", p.title, resp, err)
		}
		if got := resp.Header.Get("HX-Push-Url"); got != "/dashboard/posts" {
			t.Errorf("expected HX-Push-Url=/dashboard/posts, got %q", got)
		}
		createBody := readBody(t, resp)
		if !strings.Contains(createBody, p.title) {
			t.Errorf("expected created post %q in create response", p.title)
		}
	}

	// 7. Dashboard tab swap targeting div#dashboard-content returns inner content + partial #dashboard-nav.
	tabReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/dashboard/posts", nil)
	tabReq.Header.Set("HX-Request", "true")
	tabReq.Header.Set("HX-Target", "div#dashboard-content")
	for _, c := range cookies {
		tabReq.AddCookie(c)
	}
	resp, err = client.Do(tabReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("tab switch = (%v, %v)", resp, err)
	}
	tabHTML := readBody(t, resp)
	if strings.Contains(tabHTML, `id="dashboard-shell"`) {
		t.Error("dashboard-content target should not re-wrap dashboard-shell")
	}
	if !strings.Contains(tabHTML, `<template hx type="partial" hx-target="#dashboard-nav"`) || !strings.Contains(tabHTML, "Alpha HTMX Post") {
		t.Errorf("unexpected dashboard tab HTML: %s", tabHTML)
	}

	// 8. Live search targeting div#post-grid-region filters posts cleanly.
	searchReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/?q=Alpha", nil)
	searchReq.Header.Set("HX-Request", "true")
	searchReq.Header.Set("HX-Target", "div#post-grid-region")
	resp, err = client.Do(searchReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("live search = (%v, %v)", resp, err)
	}
	searchHTML := readBody(t, resp)
	if !strings.Contains(searchHTML, "Alpha HTMX Post") || strings.Contains(searchHTML, "Beta Notes") {
		t.Errorf("live search did not filter properly: %s", searchHTML)
	}

	// 9. Author filter returns matching posts and active author filter notice.
	authorReq, _ := http.NewRequest(http.MethodGet, srv.URL+"/?author=htmxuser", nil)
	authorReq.Header.Set("HX-Request", "true")
	authorReq.Header.Set("HX-Target", "main#app-main")
	resp, err = client.Do(authorReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("author filter = (%v, %v)", resp, err)
	}
	authorHTML := readBody(t, resp)
	if !strings.Contains(authorHTML, `id="author-filter"`) || !strings.Contains(authorHTML, "Alpha HTMX Post") {
		t.Errorf("author filter missing filter indicator or post: %s", authorHTML)
	}

	// 10. In-place card delete removes single card, and deleting last card swaps empty state.
	posts, err := store.ListPosts(context.Background())
	if err != nil || len(posts) != 2 {
		t.Fatalf("expected 2 posts, got %d (%v)", len(posts), err)
	}

	del1, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/formapi/posts/delete/%d", srv.URL, posts[0].ID), strings.NewReader("from=list"))
	del1.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	del1.Header.Set("HX-Request", "true")
	del1.Header.Set("HX-Target", fmt.Sprintf("article#post-card-%d", posts[0].ID))
	for _, c := range cookies {
		del1.AddCookie(c)
	}
	resp, err = client.Do(del1)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("delete 1 = (%v, %v)", resp, err)
	}
	if resp.Header.Get("HX-Retarget") != "" {
		t.Errorf("first delete should not retarget when posts remain, got %q", resp.Header.Get("HX-Retarget"))
	}
	del1Body := readBody(t, resp)
	if !strings.Contains(del1Body, `<template hx type="partial" hx-target="#header-stats"`) {
		t.Errorf("missing header stats partial on delete: %s", del1Body)
	}

	del2, _ := http.NewRequest(http.MethodPost, fmt.Sprintf("%s/formapi/posts/delete/%d", srv.URL, posts[1].ID), strings.NewReader("from=list"))
	del2.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	del2.Header.Set("HX-Request", "true")
	del2.Header.Set("HX-Target", fmt.Sprintf("article#post-card-%d", posts[1].ID))
	for _, c := range cookies {
		del2.AddCookie(c)
	}
	resp, err = client.Do(del2)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("delete 2 = (%v, %v)", resp, err)
	}
	if got := resp.Header.Get("HX-Retarget"); got != "#user-posts-list" {
		t.Errorf("last delete should retarget #user-posts-list, got %q", got)
	}
	del2Body := readBody(t, resp)
	if !strings.Contains(del2Body, "You have no posts yet.") {
		t.Errorf("expected empty state after deleting last post, got: %s", del2Body)
	}

	// 11. HTMX logout swaps homepage, pushes /, and resets nav to Login.
	logoutReq, _ := http.NewRequest(http.MethodPost, srv.URL+"/logout", nil)
	logoutReq.Header.Set("HX-Request", "true")
	for _, c := range cookies {
		logoutReq.AddCookie(c)
	}
	resp, err = client.Do(logoutReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("HTMX logout = (%v, %v)", resp, err)
	}
	if got := resp.Header.Get("HX-Push-Url"); got != "/" {
		t.Errorf("expected HX-Push-Url=/, got %q", got)
	}
	logoutBody := readBody(t, resp)
	if !strings.Contains(logoutBody, `href="/login"`) {
		t.Errorf("unexpected HTMX logout body: %s", logoutBody)
	}
}

func TestSSEEventsBroadcastAndClose(t *testing.T) {
	h, _, srv := setupTestServer(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	sseReq, _ := http.NewRequestWithContext(ctx, http.MethodGet, srv.URL+"/events", nil)
	sseResp, err := http.DefaultClient.Do(sseReq)
	if err != nil || sseResp.StatusCode != http.StatusOK {
		t.Fatalf("GET /events = (%v, %v)", sseResp, err)
	}
	defer sseResp.Body.Close()
	if ct := sseResp.Header.Get("Content-Type"); !strings.HasPrefix(ct, "text/event-stream") {
		t.Fatalf("expected text/event-stream, got %q", ct)
	}

	reader := bufio.NewReader(sseResp.Body)
	// Read initial ": connected" comment line + blank line.
	line, err := reader.ReadString('\n')
	if err != nil || !strings.HasPrefix(line, ": connected") {
		t.Fatalf("expected initial SSE connected comment, got %q (%v)", line, err)
	}
	_, _ = reader.ReadString('\n')

	// Trigger a broadcast and verify both the <template hx type="partial"> stats message and posts-changed event stream through.
	h.broadcastFeedAndStats(sseReq)

	var seenStatsPartial, seenPostsChanged bool
	for !(seenStatsPartial && seenPostsChanged) {
		l, err := reader.ReadString('\n')
		if err != nil {
			t.Fatalf("reading SSE stream: %v", err)
		}
		if strings.Contains(l, `<template hx type="partial" hx-target="#header-stats"`) {
			seenStatsPartial = true
		}
		if strings.TrimSpace(l) == "event: posts-changed" {
			seenPostsChanged = true
		}
	}

	// Closing the handler terminates active SSE streams cleanly.
	h.Close()
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(io.Discard, reader)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("SSE stream did not terminate after Handler.Close()")
	}
}
