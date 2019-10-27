package services

import "net/http"

// Setup services
func Setup(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/register",register)
}
