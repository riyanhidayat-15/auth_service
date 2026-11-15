package handler

import "net/http"

func NewRouter(auth *AuthHandler) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok"))
	})

	mux.HandleFunc("POST /register", auth.Register)
	mux.HandleFunc("POST /login", auth.Login)

	return mux
}