package middleware

import (
	"log"
	"net/http"

	"profile-service/internal/utils"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {
		log.Println("in auth middleware 2")
		ok, _ := utils.VerifyJWT(w, r, false)
		if !ok {
			return
		}
		next.ServeHTTP(w, r)
	})
}
