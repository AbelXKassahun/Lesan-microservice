package api

import(
	"net/http"
)

func Routes() *http.ServeMux {
	router := http.NewServeMux()
	router.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Gamification Service is up!"))
	})
	
	return router
}