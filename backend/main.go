package main

import (
	"net/http"

	"fmt"
	"log"
	"os"
	"time"

	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"

	"github.com/golang-jwt/jwt/v5"
	"github.com/golang-jwt/jwt/v5/request"
	"github.com/joho/godotenv"
)

func AuthMiddleware(pub *rsa.PublicKey) gin.HandlerFunc {
	return func(c *gin.Context) {
		tok, err := request.ParseFromRequest(
			c.Request,
			request.OAuth2Extractor,
			func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return pub, nil
			},
			request.WithClaims(&jwt.RegisteredClaims{}),
			request.WithParser(jwt.NewParser(
				jwt.WithValidMethods([]string{"RS256"}),
				jwt.WithIssuer("crfantasyleague"),
				jwt.WithExpirationRequired(),
			)),
		)

		if err != nil || !tok.Valid {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims := tok.Claims.(*jwt.RegisteredClaims)
		if claims.Subject == "" {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set("userID", claims.Subject)
		c.Next()
	}
}

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
	if err != nil {
		log.Fatal("Could not load private key: ", err)
	}

	private_key, err := jwt.ParseRSAPrivateKeyFromPEM(private_key_bytes)
	if err != nil {
		log.Fatal("Could not parse private key: ", err)
	}

	public_key := private_key.PublicKey

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
			return
		}

		code := c.Query("code")

		if code == "" {
			c.String(http.StatusUnauthorized, "Failed to login!")
			return
		}

		ctx := c.Request.Context()

		oauth_token, err := oauth_config.Exchange(ctx, code)
		if err != nil {
			c.String(http.StatusBadGateway, "Failed to exchange authorization code")
			return
		}

		client := oauth_config.Client(ctx, oauth_token)

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
			Username 	string `json:"username"`
		}{}

		if err := json.NewDecoder(res.Body).Decode(&user); err != nil {
			c.String(http.StatusBadGateway, "Failed to decode Discord response")
			return
		}

		claims := jwt.RegisteredClaims {
			Subject: user.ID,
			Issuer: "crfantasyleague",
			IssuedAt: jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(5 * time.Minute)),
		}

		jwt_token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		signed, err := jwt_token.SignedString(private_key)
		if err != nil {
			log.Printf("Failed to sign JWT: %v", err)
			c.String(http.StatusInternalServerError, "Failed to sign token")
			return
		}

		c.JSON(200, gin.H{"token": signed})
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
	
	api_group := r.Group("/api")
	api_group.Use(AuthMiddleware(&public_key))

	api_group.GET("/me", func(c *(gin.Context)) {
		c.JSON(http.StatusOK, gin.H {
			"userID": c.GetString("userID"),
		})
	})

	r.Run()
}
