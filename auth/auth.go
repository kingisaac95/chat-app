package auth

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

type authHandler struct {
	next http.Handler
}

func (h *authHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("auth")

	// not authenticated
	if err == http.ErrNoCookie {
		w.Header().Set("Location", "/login")
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	}

	// other errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	// call the passed in handler via next
	h.next.ServeHTTP(w, r)
}

// Required creates a wrapper for other http.Handler
func Required(handler http.Handler) http.Handler {
	return &authHandler{next: handler}
}

// LoginHandler handles third-party login process
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	const (
		LOGIN    string = "login"
		CALLBACK string = "callback"
	)

	segs := strings.Split(r.URL.Path, "/")

	if len(segs) < 4 {
		http.Error(w, "Error: incomplete aruguments", http.StatusBadRequest)
		return
	}

	action := segs[2]
	provider := segs[3]
	switch action {
	case LOGIN:
		log.Println("TODO: handle login for", provider)
	default:
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Auth action %s not supported", action)
	}
}
