package methods

import (
	"net/http"
)

type MethodHandler struct {
	Handler http.HandlerFunc
	Method  string
}

func Map(mHandles []MethodHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		matched := false
		for _, mHandle := range mHandles {
			if r.Method == mHandle.Method {
				mHandle.Handler(w, r)
				matched = true
				break
			}
		}

		if !matched {
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	}
}
