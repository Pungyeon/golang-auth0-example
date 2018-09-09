package handlers

import (
	"log"
	"net/http"

	"github.com/auth0-community/auth0"
	"github.com/gin-gonic/gin"
	jose "gopkg.in/square/go-jose.v2"
	jwt "gopkg.in/square/go-jose.v2/jwt"
)

// AuthHandler is a endpoint handler for checking validity of JWT tokens
type AuthHandler struct {
	validator *auth0.JWTValidator
}

// NewAuthHandler will return an AuthHandler initialising a JWT validator in the process
func NewAuthHandler(audience string, domain string) *AuthHandler {
	client := auth0.NewJWKClient(auth0.JWKClientOptions{URI: domain + ".well-known/jwks.json"}, nil)
	configuration := auth0.NewConfiguration(client, []string{audience}, domain, jose.RS256)

	return &AuthHandler{
		validator: auth0.NewValidator(configuration, nil),
	}
}

// Required will verify that a token received from an http request
// is valid and signy by authority
func (handler *AuthHandler) Required() gin.HandlerFunc {
	return func(c *gin.Context) {
		token, err := handler.validator.ValidateRequest(c.Request)
		if err != nil {
			log.Println(err)
			terminateWithError(http.StatusUnauthorized, "token is not valid", c)
			return
		}
		claims := jwt.Claims{}
		err = handler.validator.Claims(c.Request, token, &claims)
		if err != nil {
			terminateWithError(http.StatusUnauthorized, "could not retrieve subject from claim", c)
			return
		}
		c.Request.Header.Add("username", claims.Subject)
		c.Next()
	}
}

func terminateWithError(statusCode int, message string, c *gin.Context) {
	c.JSON(statusCode, gin.H{"error": message})
	c.Abort()
}
