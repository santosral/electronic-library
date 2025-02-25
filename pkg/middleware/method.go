package middleware

import (
	"electronic-library/pkg/response"
	"encoding/json"
	"net/http"
)

func Method(m string) Middleware {
	return func(next http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			if r.Method != m {
				w.Header().Set("Content-Type", "application/json")

				w.WriteHeader(http.StatusMethodNotAllowed)
				resp := response.NewErrorResponse("Invalid HTTP method")
				json.NewEncoder(w).Encode(resp)
				return
			}
			next(w, r)
		}
	}
}
