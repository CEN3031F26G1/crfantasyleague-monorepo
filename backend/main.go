package main

import (
	"net/http"
	"net/url"

	"log"
	"os"

	"crypto/rand"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Could not load environment variable file")
	}

	oauth_config := &oauth2.Config{
		ClientID:     "1565415990093684767",
		ClientSecret: os.Getenv("DISCORD_CLIENT_SECRET"),
		RedirectURL:  "http://localhost:8080/auth/callback",
		Scopes:       []string{"identify"},
		Endpoint: oauth2.Endpoint{
			AuthURL:  "https://discord.com/oauth2/authorize",
			TokenURL: "https://discord.com/api/oauth2/token",
		},
	}

	private_key_bytes, err := os.ReadFile("./private.pem")
	log.Fatal("Could not load private key: ", err)

	private_key, err := jwt.ParseRSAPrivateKeyFromPEM(private_key_bytes)
	log.Fatal("Could not parse private key: ", err)

	r := gin.Default()

	r.GET("/auth/callback", func(c *gin.Context) {
		expected_state, err := c.Cookie("oauth_state")
		
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(
			"oauth_state",
			"",
			-1,
			"/auth/callback",
			"",
			false,
			true,
		)

		state := c.Query("state")

		if err != nil || state != expected_state {
			c.String(http.StatusBadRequest, "Invalid OAuth state; please re-login!")
		}

		code := c.Query("code")

		if code == "" {
			c.String(http.StatusUnauthorized, "Failed to login!")
			return
		}

		ctx := c.Request.Context()

		token, err := oauth_config.Exchange(ctx, code)
		if err != nil {
			c.String(http.StatusBadGateway, "Failed to exchange authorization code")
			return
		}

		client := oauth_config.Client(ctx, token)

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			"https://discord.com/api/users/@me",
			nil,
		)

		if err != nil {
			c.String(http.StatusInternalServerError, "Failed to create user request")
			return
		}

		res, err := client.Do(req)
		if err != nil {
			c.String(http.StatusBadGateway, "Failed to fetch Discord user")
			return
		}
		defer res.Body.Close()

		if res.StatusCode != http.StatusOK {
			c.String(http.StatusBadGateway, "Discord rejected user request")
			return
		}

		user := struct {
			ID 			string `json:"id"`
			Username 	string `username:"user"`
		}{}

		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			c.String(http.StatusBadGateway, "Failed to decode Discord response")
			return
		}

		c.JSON(http.StatusOK, user)
	})

	r.GET("/login", func(c *gin.Context) {
		b := make([]byte, 32)
		if _, err := rand.Read(b); err != nil {
			c.String(http.StatusInternalServerError, "Failed to initialize OAuth login")
			return
		}

		state := base64.RawURLEncoding.EncodeToString(b)
		c.SetSameSite(http.SameSiteLaxMode)
		c.SetCookie(
			"oauth_state",
			state,
			300,
			"/auth/callback",
			"",
			false, // switch to true with HTTPS
			true,
		)

		c.Redirect(http.StatusFound, oauth_config.AuthCodeURL(state))
	})

	r.Run()
}
