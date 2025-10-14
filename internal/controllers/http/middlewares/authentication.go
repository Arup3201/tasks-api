package middlewares

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"

	httperrors "github.com/Arup3201/gotasks/internal/controllers/http/errors"
	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

func Authenticate(secureEndpoints []string) gin.HandlerFunc {
	return func(c *gin.Context) {
		for _, endpoint := range secureEndpoints {
			if strings.Index(c.Request.URL.Path, endpoint) == 0 {
				token, err := verifyToken(c.Request)
				if err != nil {
					log.Printf("authentication error: %v", err)
					abortAuthentication(c)
					return
				}
				if token != nil {
					c.Next()
				} else {
					log.Printf("authentication error: invalid token")
					abortAuthentication(c)
				}
				return
			}
		}

		c.Next()
	}
}

func abortAuthentication(c *gin.Context) {
	c.Header("WWW-Authenticate", "Bearer token68")
	c.Error(httperrors.UnauthorizedError())
	c.Abort()
}

func verifyToken(request *http.Request) (jwt.Token, error) {
	strToken, err := getAuthHeader(request)
	if err != nil {
		return nil, err
	}
	jwksKeySet, err := jwk.Fetch(request.Context(), fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", utils.Config.KeycloakServerUrl, utils.Config.KeycloakRealName))
	if err != nil {
		return nil, err
	}
	token, err := jwt.Parse([]byte(strToken), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
	if err != nil {
		return nil, err
	}
	return token, nil
}

func getAuthHeader(request *http.Request) (string, error) {
	header := strings.Fields(request.Header.Get("Authorization"))
	if len(header) == 0 || header[0] != "Bearer" {
		return "", errors.New("malformed token")
	}
	return header[1], nil
}
