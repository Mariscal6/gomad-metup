package main

import (
	"base-image/views"
	"log"
	"net/http"
	_ "net/http/pprof"
)

// this is a base http server
// with http pprof enabled
func main() {
	log.Println("Starting server on :8080")
	server := http.NewServeMux()

	server.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	server.Handle("/debug/pprof/", http.DefaultServeMux)
	server.Handle("/assets/", http.StripPrefix("/assets", http.FileServer(http.Dir("./assets/"))))
	server.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		err := views.HelloGopher().Render(r.Context(), w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	http.ListenAndServe(":8080", server)
}
