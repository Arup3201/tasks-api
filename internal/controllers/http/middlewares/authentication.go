package middlewares

import (
	"errors"
	"net/http"
	"strings"

	"github.com/Arup3201/gotasks/internal/utils"
	"github.com/gin-gonic/gin"
	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
)

type middlewareCallback func() gin.HandlerFunc

func Authenticate(secureEndpoints []string) middlewareCallback {
	return func() gin.HandlerFunc {
		return func(c *gin.Context) {
			for _, endpoint := range secureEndpoints {
				if strings.Index(c.Request.URL.Path, endpoint) == 0 {
					token, err := verifyToken(c.Request)
					if err != nil {
						c.IndentedJSON(http.StatusBadRequest, gin.H{"reason": err.Error()})
					}
					if token != nil {
						c.Next()
					} else {
						c.IndentedJSON(http.StatusForbidden, gin.H{"message": "Not authenticated"})
					}
					return
				}
			}

			c.Next()
		}
	}
}

func verifyToken(request *http.Request) (jwt.Token, error) {
	strToken, err := getAuthHeader(request)
	if err != nil {
		return nil, err
	}
	jwksKeySet, err := jwk.Fetch(request.Context(), utils.Config.KEYCLOAK_JWT_URL)
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
