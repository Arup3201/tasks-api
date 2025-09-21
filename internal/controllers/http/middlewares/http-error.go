package middlewares

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ErrorResponder interface {
	ResponseBody() ([]byte, error)
	ResponseHeader() (int, map[string]string)
	Error() string
}

func HttpErrorResponse() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			errorResponder, ok := err.(ErrorResponder)
			if !ok {
				log.Printf("Internal server error: %v", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occured"})
				return
			}

			body, err := errorResponder.ResponseBody()
			if err != nil {
				log.Printf("errorResponder.ResponseBody error: %v", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occured"})
				return
			}

			status, headers := errorResponder.ResponseHeader()
			for k, v := range headers {
				c.Writer.Header().Set(k, v)
			}
			c.Writer.WriteHeader(status)
			c.Writer.Write(body)
		}
	}
}
