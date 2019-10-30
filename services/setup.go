package services

import "net/http"

// Setup services
func Setup(mux *http.ServeMux) {
	mux.HandleFunc("/api/user/users",getUsers)
	mux.HandleFunc("/api/user/user",userHandler)
	mux.HandleFunc("/api/user/changePassword",changePassword)
	mux.HandleFunc("/api/user/login",login)
	mux.HandleFunc("/api/user/logout",logout)
	mux.HandleFunc("/api/user/session",sessionHandler)
}
