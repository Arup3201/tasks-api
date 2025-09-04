package logging

import (
	"log"
	"maps"
	"net/http"
	"net/http/httptest"

	"github.com/Arup3201/gotasks/internal/middlewares"
)

func HttpLogger() middlewares.Middleware {
	return func (fn http.HandlerFunc) http.HandlerFunc {
		return func (w http.ResponseWriter, r *http.Request) {
			recorder := httptest.NewRecorder()
			fn(recorder, r)

			log.Printf("%s %s - %d", r.Method, r.URL.Path, recorder.Code)

			maps.Copy(w.Header(), recorder.Result().Header)
			_, err := recorder.Body.WriteTo(w)
			if err!=nil {
				log.Fatalf("HttpLogger write to response body error: %v", err)
			}
		}
	}
}
