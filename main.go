package main

import (
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
)

type Login struct {
	Username string `json:"username"`
	Password string `json:"-"`
}

var (
	store *sessions.CookieStore
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file", err)
	}

	store = sessions.NewCookieStore([]byte(os.Getenv("SESSION_KEY")))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.ServeFile(w, r, "index.html")
			return
		}

		login := Login{
			Username: r.FormValue("username"),
			Password: r.FormValue("password"),
		}

		if login.Password == os.Getenv("PASSWORD") {
			session, _ := store.Get(r, "session")
			session.Values["authenticated"] = true
			session.Values["username"] = login.Username
			session.Save(r, w)
			http.Redirect(w, r, "/chat?username="+login.Username, http.StatusSeeOther)
		} else {
			http.Error(w, "Forbidden", http.StatusForbidden)
		}
	})

	http.HandleFunc("/chat", func(w http.ResponseWriter, r *http.Request) {
		session, _ := store.Get(r, "session")

		if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		http.ServeFile(w, r, "chat.html")
	})

	http.HandleFunc("/conn", handleConnection)

	go handleMessages()

	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic("Error starting server: " + err.Error())
	}
}
