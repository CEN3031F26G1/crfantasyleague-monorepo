package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"context"
	"crfantasyleague/backend/database"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
)

const (
	sessionCookie = "crfl_session"
	sessionTTL    = 12 * time.Hour
)

type sessionClaims struct {
	Username   string `json:"username,omitempty"`
	GlobalName string `json:"globalName,omitempty"`
	AvatarURL  string `json:"avatarUrl,omitempty"`
	jwt.RegisteredClaims
}

// Continue accepting bearer tokens for API clients, in addition to browser cookies.
func AuthMiddleware(pub *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		var raw string
		if header := c.GetHeader("Authorization"); header != "" {
			parts := strings.Fields(header)
			if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}
			raw = parts[1]
		} else {
			raw, _ = c.Cookie(sessionCookie)
		}
		claims := &sessionClaims{}
		tok, err := jwt.ParseWithClaims(raw, claims, func(t *jwt.Token) (interface{}, error) {
			return pub, nil
		}, jwt.WithValidMethods([]string{"RS256"}), jwt.WithIssuer("crfantasyleague"), jwt.WithExpirationRequired())
		if err != nil || !tok.Valid || claims.Subject == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		c.Set("userID", claims.Subject)
		c.Set("session", claims)
		c.Next()
	}
}

type appConfig struct {
	FrontendURL    string
	SecureCookies  bool
	DiscordUserURL string
	OAuth          *oauth2.Config
	DB             *pgxpool.Pool
}

func setCookie(c *gin.Context, name, value, path string, maxAge int, secure bool) {
	c.SetSameSite(http.SameSiteLaxMode)
	c.SetCookie(name, value, maxAge, path, "", secure, true)
}

func newRouter(config appConfig, key *rsa.PrivateKey) *gin.Engine {
	r := gin.Default()
	// Do not trust arbitrary X-Forwarded-* headers from the public internet.
	_ = r.SetTrustedProxies(nil)
	frontend := strings.TrimRight(config.FrontendURL, "/")
	loginError := func(c *gin.Context, reason string) {
		c.Redirect(http.StatusSeeOther, frontend+"/login?error="+url.QueryEscape(reason))
	}
	login := func(c *gin.Context) {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			loginError(c, "unavailable")
			return
		}
		state := base64.RawURLEncoding.EncodeToString(b)
		setCookie(c, "oauth_state", state, "/auth/callback", 300, config.SecureCookies)
		c.Header("Cache-Control", "no-store")
		c.Redirect(http.StatusFound, config.OAuth.AuthCodeURL(state))
	}
	r.GET("/login", login) // Keep the original backend login URL working.
	r.GET("/auth/login", login)

	r.GET("/auth/callback", func(c *gin.Context) {
		c.Header("Cache-Control", "no-store")
		expected, err := c.Cookie("oauth_state")
		setCookie(c, "oauth_state", "", "/auth/callback", -1, config.SecureCookies)
		state := c.Query("state")
		if err != nil || expected == "" || state == "" || subtle.ConstantTimeCompare([]byte(state), []byte(expected)) != 1 {
			loginError(c, "state")
			return
		}
		if c.Query("error") != "" || c.Query("code") == "" {
			loginError(c, "cancelled")
			return
		}

		ctx := c.Request.Context()
		token, err := config.OAuth.Exchange(ctx, c.Query("code"))
		if err != nil {
			loginError(c, "unavailable")
			return
		}
		client := config.OAuth.Client(ctx, token)
		client.Timeout = 15 * time.Second
		req, err := http.NewRequestWithContext(ctx, http.MethodGet, config.DiscordUserURL, nil)
		if err != nil {
			loginError(c, "unavailable")
			return
		}
		res, err := client.Do(req)
		if err != nil {
			loginError(c, "unavailable")
			return
		}
		defer res.Body.Close()
		user := struct {
			ID         string `json:"id"`
			Username   string `json:"username"`
			GlobalName string `json:"global_name"`
			Avatar     string `json:"avatar"`
		}{}
		if res.StatusCode != http.StatusOK || json.NewDecoder(res.Body).Decode(&user) != nil || user.ID == "" {
			loginError(c, "unavailable")
			return
		}
		avatarURL := ""
		if user.Avatar != "" {
			avatarURL = fmt.Sprintf("https://cdn.discordapp.com/avatars/%s/%s.png?size=128", url.PathEscape(user.ID), url.PathEscape(user.Avatar))
		}
		if config.DB != nil {
			err := database.SaveDiscordUser(
				ctx,
				config.DB,
				user.ID,
				user.Username,
				user.GlobalName,
				avatarURL,
			)

			if err != nil {
				log.Printf("Failed to save Discord user: %v", err)
				loginError(c, "unavailable")
				return
			}
		}
		claims := sessionClaims{
			Username: user.Username, GlobalName: user.GlobalName, AvatarURL: avatarURL,
			RegisteredClaims: jwt.RegisteredClaims{
				Subject: user.ID, Issuer: "crfantasyleague",
				IssuedAt: jwt.NewNumericDate(time.Now()), ExpiresAt: jwt.NewNumericDate(time.Now().Add(sessionTTL)),
			},
		}
		signed, err := jwt.NewWithClaims(jwt.SigningMethodRS256, claims).SignedString(key)
		if err != nil {
			loginError(c, "unavailable")
			return
		}
		// The JWT never enters a URL or JavaScript-accessible storage.
		setCookie(c, sessionCookie, signed, "/", int(sessionTTL.Seconds()), config.SecureCookies)
		c.Redirect(http.StatusSeeOther, frontend+"/")
	})

	r.POST("/auth/logout", func(c *gin.Context) {
		// Cookies authenticate requests automatically; require a trusted browser origin
		// on this state-changing endpoint. Non-browser clients may omit Origin.
		origin := c.GetHeader("Origin")
		if origin != "" && origin != frontend {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		if c.GetHeader("Sec-Fetch-Site") == "cross-site" {
			c.AbortWithStatus(http.StatusForbidden)
			return
		}
		c.Header("Cache-Control", "no-store")
		setCookie(c, sessionCookie, "", "/", -1, config.SecureCookies)
		c.Status(http.StatusNoContent)
	})

	api := r.Group("/api")
	api.Use(func(c *gin.Context) { c.Header("Cache-Control", "no-store"); c.Next() })
	api.Use(AuthMiddleware(&key.PublicKey))
	api.GET("/me", func(c *gin.Context) {
		claims := c.MustGet("session").(*sessionClaims)
		c.JSON(http.StatusOK, gin.H{
			"userID": claims.Subject, "username": claims.Username,
			"globalName": claims.GlobalName, "avatarUrl": claims.AvatarURL,
			"expiresAt": claims.ExpiresAt.Time.UTC().Format(time.RFC3339),
		})
	})
	return r
}

func envOr(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func validOrigin(raw string) bool {
	u, err := url.Parse(raw)
	return err == nil && (u.Scheme == "http" || u.Scheme == "https") && u.Host != "" && u.User == nil && (u.Path == "" || u.Path == "/") && u.RawQuery == "" && u.Fragment == ""
}

func main() {
	if err := godotenv.Load(".env"); err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Fatal("Could not load environment file: ", err)
	}
	if os.Getenv("DISCORD_CLIENT_SECRET") == "" {
		log.Fatal("DISCORD_CLIENT_SECRET is required")
	}
	frontend := strings.TrimRight(envOr("FRONTEND_URL", "http://localhost:5173"), "/")
	if !validOrigin(frontend) {
		log.Fatal("FRONTEND_URL must be an http(s) origin without a path")
	}
	redirect := envOr("DISCORD_REDIRECT_URL", "http://localhost:8080/auth/callback")
	callback, err := url.Parse(redirect)
	if err != nil || callback.Host == "" || (callback.Scheme != "http" && callback.Scheme != "https") || callback.Path != "/auth/callback" || callback.User != nil || callback.RawQuery != "" || callback.Fragment != "" {
		log.Fatal("DISCORD_REDIRECT_URL must be an http(s) URL ending in /auth/callback")
	}
	pem, err := os.ReadFile(envOr("JWT_PRIVATE_KEY_PATH", "./private.pem"))
	if err != nil {
		log.Fatal("Could not load private key: ", err)
	}
	key, err := jwt.ParseRSAPrivateKeyFromPEM(pem)
	if err != nil {
		log.Fatal("Could not parse private key: ", err)
	}

	//connect to PostgreSQL
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal("Failed to connect to PostgreSQL: ", err)
	}
	defer db.Close()

	log.Println("Connected to PostgreSQL successfully!")

	// Test reading from the users table
	var userCount int

	err = db.QueryRow(ctx, "SELECT COUNT(*) FROM users").Scan(&userCount)
	if err != nil {
		log.Fatal("Failed to query users table: ", err)
	}

	log.Printf("Total users in database: %d", userCount)

	// Bound OAuth exchange requests as well as Discord profile requests.
	http.DefaultClient.Timeout = 15 * time.Second
	r := newRouter(appConfig{
		DB: db, FrontendURL: frontend, SecureCookies: callback.Scheme == "https",
		DiscordUserURL: "https://discord.com/api/users/@me",
		OAuth: &oauth2.Config{
			ClientID:     envOr("DISCORD_CLIENT_ID", "1565415990093684767"),
			ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"), RedirectURL: redirect,
			Scopes:   []string{"identify"},
			Endpoint: oauth2.Endpoint{AuthURL: "https://discord.com/oauth2/authorize", TokenURL: "https://discord.com/api/oauth2/token"},
		},
	}, key)
	server := &http.Server{Addr: ":" + envOr("PORT", "8080"), Handler: r, ReadHeaderTimeout: 10 * time.Second}
	log.Fatal(server.ListenAndServe())
}
