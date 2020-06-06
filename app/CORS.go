package app

import (
	"net/http"
)

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Set headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods",
			"GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Accept, Content-Type, Content-Length, Accept-Encoding, "+
				"X-CSRF-Token, Authorization, Origin, X-Auth-Token")
		w.Header().Set("Access-Control-Expose-Headers",
			"Authorization")
		// Next
		next.ServeHTTP(w, r)
		return
	})
}
