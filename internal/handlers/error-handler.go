package handlers

import (
	"log"
	"net/http"

	"github.com/Arup3201/gotasks/internal/handlers/apperr"
	"github.com/gin-gonic/gin"
)

func HandleErrors() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last().Err

			clientError, ok := err.(apperr.ClientError)
			if !ok {
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occured"})
				return
			}

			body, err := clientError.ResponseBody()
			if err != nil {
				log.Printf("ClientError.ResponseBody error: %v", err)
				c.IndentedJSON(http.StatusInternalServerError, gin.H{"message": "An unexpected error occured"})
				return
			}

			status, headers := clientError.ResponseHeader()
			for k, v := range headers {
				c.Writer.Header().Set(k, v)
			}
			c.Writer.WriteHeader(status)
			c.Writer.Write(body)
		}
	}
}
