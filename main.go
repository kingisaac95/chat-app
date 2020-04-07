package main

import (
	"chat_app/auth"
	"chat_app/chat"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"text/template"

	"github.com/joho/godotenv"
	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/google"
	"github.com/stretchr/objx"
)

// templ is a single template
type templateHandler struct {
	once     sync.Once
	filename string
	templ    *template.Template
}

func (t *templateHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	t.once.Do(func() {
		t.templ = template.Must(template.ParseFiles(filepath.Join("templates", t.filename)))
	})
	data := map[string]interface{}{
		"Host": r.Host,
	}
	if authCookie, err := r.Cookie("auth"); err == nil {
		data["UserData"] = objx.MustFromBase64(authCookie.Value)
	}
	t.templ.Execute(w, data)
}

func main() {
	var port = flag.String("port", ":8080", "The port where our application runs.")
	flag.Parse() // parse flags and extract appropriate information

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	googleClientID := os.Getenv("GOOGLE_CLIENT_ID")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectURL := os.Getenv("GOOGLE_REDIRECT_URL")
	securityKey := os.Getenv("SECURITY_KEY")

	gomniauth.SetSecurityKey(securityKey)
	gomniauth.WithProviders(google.New(googleClientID, googleClientSecret, googleRedirectURL))

	// new room setup
	r := chat.NewRoom()

	http.Handle("/login", &templateHandler{filename: "login.html"})
	http.Handle("/chat", auth.Required(&templateHandler{filename: "chat.html"}))
	http.HandleFunc("/auth/", auth.LoginHandler)
	http.Handle("/room", r)

	// initialize the room
	go r.Run()

	// start the web server or log error
	log.Println("Starting server on", *port)
	if err := http.ListenAndServe(*port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
