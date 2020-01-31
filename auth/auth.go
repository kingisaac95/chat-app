package auth

import "net/http"

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
