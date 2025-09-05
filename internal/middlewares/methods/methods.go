package methods

import (
	"log"
	"net/http"
	"slices"

	"github.com/Arup3201/gotasks/internal/middlewares"
)

func Methods(methods []string) middlewares.Middleware {
	return func (fn http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			if slices.Contains(methods, r.Method) {
				fn(w, r)
			} else {
				log.Printf("%s %s - %d", r.Method, r.URL.Path, http.StatusForbidden)
				log.Printf("[SERVER] Method not allowed")
			}
		}
	}
}
