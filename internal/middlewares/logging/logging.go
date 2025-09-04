package logging

import (
	"log"
	"maps"
	"net/http"
	"net/http/httptest"
)

func HttpLogger(fn http.HandlerFunc) http.HandlerFunc {
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
