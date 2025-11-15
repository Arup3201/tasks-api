package middlewares

import (
	"encoding/json"
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
		if !utils.Config.Testing {
			for _, endpoint := range secureEndpoints {
				if strings.Index(c.Request.URL.Path, endpoint) == 0 {
					token, err := verifyToken(c.Request)
					if err != nil {
						log.Printf("authentication error: %v", err)
						abortAuthentication(c)
						return
					}

					request, err := http.NewRequest("GET", fmt.Sprintf("%s/realms/%s/protocol/openid-connect/userinfo", utils.Config.KeycloakServerUrl, utils.Config.KeycloakRealName), nil)
					if err != nil {
						log.Fatalf("http new request error: %v", err)
					}
					request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
					client := &http.Client{}
					response, err := client.Do(request)
					if err != nil {
						log.Fatalf("http client do error: %v", err)
					}
					defer response.Body.Close()

					if response.StatusCode != http.StatusOK {
						log.Printf("login error: keycloak userInfo response error: got response with status %d", response.StatusCode)
						c.Error(httperrors.InternalServerError(fmt.Errorf("failed to fetch userInfo from Auth server")))
						return
					}

					var userInfo struct {
						UserId   string `json:"sub"`
						Username string `json:"preferred_username"`
					}
					if err = json.NewDecoder(response.Body).Decode(&userInfo); err != nil {
						log.Printf("login error: keycloak userInfo response encoding error: %v", err)
						c.Error(httperrors.InternalServerError(fmt.Errorf("failed to decode userInfo from Auth server")))
						return
					}

					c.Set("username", userInfo.Username)
					c.Next()
				}
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

func verifyToken(request *http.Request) (string, error) {
	strToken, err := getAuthHeader(request)
	if err != nil {
		return "", err
	}

	jwksKeySet, err := jwk.Fetch(request.Context(), fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", utils.Config.KeycloakServerUrl, utils.Config.KeycloakRealName))
	if err != nil {
		return "", err
	}

	token, err := jwt.Parse([]byte(strToken), jwt.WithKeySet(jwksKeySet), jwt.WithValidate(true))
	if err != nil {
		return "", err
	}

	if token == nil {
		return "", fmt.Errorf("parsed token is null")
	}
	return strToken, nil
}

func getAuthHeader(request *http.Request) (string, error) {
	header := strings.Fields(request.Header.Get("Authorization"))
	if len(header) == 0 || header[0] != "Bearer" {
		return "", errors.New("malformed token")
	}
	return header[1], nil
}
