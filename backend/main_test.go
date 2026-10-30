package main

import (
	"crypto/rand"
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
)

func testApp(t *testing.T) (*gin.Engine, *rsa.PrivateKey) {
	t.Helper()
	gin.SetMode(gin.TestMode)
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		t.Fatal(err)
	}
	// Exercise the real OAuth exchange and user fetch without contacting Discord.
	discord := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/token":
			if err := r.ParseForm(); err != nil || r.FormValue("code") != "valid-code" {
				http.Error(w, "invalid code", http.StatusBadRequest)
				return
			}
			_, _ = w.Write([]byte(`{"access_token":"discord-token","token_type":"Bearer","expires_in":3600}`))
		case "/me":
			if r.Header.Get("Authorization") != "Bearer discord-token" {
				t.Error("Discord access token was not sent")
			}
			_, _ = w.Write([]byte(`{"id":"12345","username":"testfan","global_name":"Test Fan","avatar":"abc123"}`))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(discord.Close)
	return newRouter(appConfig{
		FrontendURL: "http://localhost:5173", DiscordUserURL: discord.URL + "/me",
		OAuth: &oauth2.Config{
			ClientID: "test-client", ClientSecret: "test-secret", RedirectURL: "http://localhost:8080/auth/callback", Scopes: []string{"identify"},
			Endpoint: oauth2.Endpoint{AuthURL: discord.URL + "/authorize", TokenURL: discord.URL + "/token"},
		},
	}, key), key
}

func TestOAuthRoundTrip(t *testing.T) {
	r, _ := testApp(t)
	login := httptest.NewRecorder()
	r.ServeHTTP(login, httptest.NewRequest(http.MethodGet, "/auth/login", nil))
	if login.Code != http.StatusFound {
		t.Fatalf("login status: %d", login.Code)
	}
	redirect, _ := url.Parse(login.Header().Get("Location"))
	state := redirect.Query().Get("state")
	if len(state) < 32 || redirect.Query().Get("scope") != "identify" {
		t.Fatal("login must include a random state and only the identify scope")
	}
	cookies := login.Result().Cookies()
	if len(cookies) != 1 || cookies[0].Value != state || !cookies[0].HttpOnly || cookies[0].SameSite != http.SameSiteLaxMode || cookies[0].Path != "/auth/callback" {
		t.Fatal("invalid OAuth state cookie")
	}

	callback := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodGet, "/auth/callback?state="+url.QueryEscape(state)+"&code=valid-code", nil)
	req.AddCookie(cookies[0])
	r.ServeHTTP(callback, req)
	if callback.Code != http.StatusSeeOther || callback.Header().Get("Location") != "http://localhost:5173/" {
		t.Fatalf("callback failed: %d %s", callback.Code, callback.Header().Get("Location"))
	}
	var session *http.Cookie
	for _, cookie := range callback.Result().Cookies() {
		if cookie.Name == sessionCookie {
			session = cookie
		}
	}
	if session == nil || !session.HttpOnly || session.SameSite != http.SameSiteLaxMode || session.Path != "/" || session.MaxAge != int(sessionTTL.Seconds()) {
		t.Fatal("callback must set an HTTP-only session cookie")
	}
	if strings.Contains(callback.Header().Get("Location"), session.Value) || strings.Contains(callback.Body.String(), session.Value) {
		t.Fatal("JWT must not appear in the redirect or response body")
	}
	me := httptest.NewRecorder()
	req = httptest.NewRequest(http.MethodGet, "/api/me", nil)
	req.AddCookie(session)
	r.ServeHTTP(me, req)
	var profile map[string]string
	if err := json.Unmarshal(me.Body.Bytes(), &profile); err != nil {
		t.Fatal(err)
	}
	if me.Code != http.StatusOK || profile["userID"] != "12345" || profile["globalName"] != "Test Fan" || profile["avatarUrl"] == "" || profile["expiresAt"] == "" {
		t.Fatalf("unexpected profile: %s", me.Body.String())
	}
	if me.Header().Get("Cache-Control") != "no-store" {
		t.Fatal("profile responses must not be cached")
	}
}

func TestInvalidOAuthStateAndCancellation(t *testing.T) {
	r, _ := testApp(t)
	for _, scenario := range []struct{ query, cookie, reason string }{
		{"?state=abc&code=valid-code", "", "state"},
		{"?state=wrong&code=valid-code", "abc", "state"},
		{"?code=valid-code", "", "state"},
		{"?state=abc&error=access_denied", "abc", "cancelled"},
		{"?state=abc", "abc", "cancelled"},
		{"?state=abc&code=invalid-code", "abc", "unavailable"},
	} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/auth/callback"+scenario.query, nil)
		if scenario.cookie != "" {
			req.AddCookie(&http.Cookie{Name: "oauth_state", Value: scenario.cookie})
		}
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusSeeOther || !strings.HasSuffix(rec.Header().Get("Location"), "error="+scenario.reason) {
			t.Errorf("unexpected error redirect for %s: %s", scenario.query, rec.Header().Get("Location"))
		}
		for _, cookie := range rec.Result().Cookies() {
			if cookie.Name == sessionCookie && cookie.Value != "" {
				t.Error("invalid OAuth attempt created a session")
			}
		}
	}
}

func TestAuthValidationAndBearerCompatibility(t *testing.T) {
	r, key := testApp(t)
	for _, tc := range []struct {
		name, subject, issuer, method string
		expiry                        time.Time
		want                          int
	}{
		{"valid", "123", "crfantasyleague", "RS256", time.Now().Add(time.Hour), 200},
		{"expired", "123", "crfantasyleague", "RS256", time.Now().Add(-time.Hour), 401},
		{"wrong issuer", "123", "other", "RS256", time.Now().Add(time.Hour), 401},
		{"empty subject", "", "crfantasyleague", "RS256", time.Now().Add(time.Hour), 401},
		{"wrong algorithm", "123", "crfantasyleague", "HS256", time.Now().Add(time.Hour), 401},
		{"missing expiration", "123", "crfantasyleague", "RS256", time.Time{}, 401},
	} {
		t.Run(tc.name, func(t *testing.T) {
			claims := jwt.RegisteredClaims{Subject: tc.subject, Issuer: tc.issuer}
			if !tc.expiry.IsZero() {
				claims.ExpiresAt = jwt.NewNumericDate(tc.expiry)
			}
			var signingKey interface{} = key
			if tc.method == "HS256" {
				signingKey = []byte("not-the-rsa-key")
			}
			raw, err := jwt.NewWithClaims(jwt.GetSigningMethod(tc.method), claims).SignedString(signingKey)
			if err != nil {
				t.Fatal(err)
			}
			for _, bearer := range []bool{true, false} {
				req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
				if bearer {
					req.Header.Set("Authorization", "Bearer "+raw)
				} else {
					req.AddCookie(&http.Cookie{Name: sessionCookie, Value: raw})
				}
				rec := httptest.NewRecorder()
				r.ServeHTTP(rec, req)
				if rec.Code != tc.want {
					t.Errorf("bearer=%v: got %d, want %d", bearer, rec.Code, tc.want)
				}
			}
		})
	}
	for _, auth := range []string{"", "invalid", "Bearer broken"} {
		req := httptest.NewRequest(http.MethodGet, "/api/me", nil)
		req.Header.Set("Authorization", auth)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("unauthenticated request returned %d", rec.Code)
		}
	}
}

func TestLogoutOriginAndCookie(t *testing.T) {
	r, _ := testApp(t)
	for _, origin := range []string{"http://localhost:5173", "https://evil.example", ""} {
		req := httptest.NewRequest(http.MethodPost, "/auth/logout", nil)
		req.Header.Set("Origin", origin)
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, req)
		if origin == "https://evil.example" {
			if rec.Code != http.StatusForbidden {
				t.Error("cross-origin logout was accepted")
			}
			continue
		}
		cookies := rec.Result().Cookies()
		if rec.Code != http.StatusNoContent || len(cookies) != 1 || cookies[0].Name != sessionCookie || cookies[0].MaxAge != -1 {
			t.Error("logout did not clear the session cookie")
		}
	}
}

func TestSecureCookies(t *testing.T) {
	rec := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(rec)
	setCookie(c, sessionCookie, "example", "/", 100, true)
	if cookies := rec.Result().Cookies(); len(cookies) != 1 || !cookies[0].Secure || !cookies[0].HttpOnly {
		t.Fatal("HTTPS session cookie must be Secure and HttpOnly")
	}
}
