package main

//https://curtisvermeeren.github.io/2018/05/13/Golang-Gorilla-Sessions.html

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/danilopolani/gocialite"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

const (
	// DefaultPort to start up server
	DefaultPort = "9000"
)

// User holds a users account information
type SSOUser struct {
	Username      string
	Authenticated bool
}

// store will hold all session data
var store *sessions.CookieStore

var gocial = gocialite.NewDispatcher()

type spaHandler struct {
	staticPath string
	indexPath  string
}

func init() {
	authKeyOne := securecookie.GenerateRandomKey(64)
	encryptionKeyOne := securecookie.GenerateRandomKey(32)

	store = sessions.NewCookieStore(
		authKeyOne,
		encryptionKeyOne,
	)

	// 30mins
	store.Options = &sessions.Options{
		MaxAge:   60 * 30,
		HttpOnly: true,
	}

	gob.Register(SSOUser{})

	// tpl = template.Must(template.ParseGlob("templates/*.gohtml"))
}

func main() {
	var port string
	if port = os.Getenv("PORT"); len(port) == 0 {
		fmt.Printf("Warning, PORT not set. Defaulting to %v\n", DefaultPort)
		port = DefaultPort
	}

	router := mux.NewRouter()

	router.HandleFunc("/auth", redirectHandler)
	router.HandleFunc("/auth/callback", callbackHandler)

	// router.HandleFunc("/", index)
	// router.HandleFunc("/login", login)
	// router.HandleFunc("/logout", logout)
	// router.HandleFunc("/forbidden", forbidden)
	// router.HandleFunc("/secret", secret)

	router.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		// an example API handler
		json.NewEncoder(w).Encode(map[string]bool{"ok": true})
	})

	spa := spaHandler{staticPath: "static", indexPath: "index.html"}
	router.PathPrefix("/").Handler(spa)

	srv := &http.Server{
		Handler:      router,
		Addr:         "127.0.0.1:" + port,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	log.Fatal(srv.ListenAndServe())
}

func (h spaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// get the absolute path to prevent directory traversal
	path, err := filepath.Abs(r.URL.Path)
	if err != nil {
		// if we failed to get the absolute path respond with a 400 bad request and stop
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// prepend the path with the path to the static directory
	path = filepath.Join(h.staticPath, path)

	// check whether a file exists at the given path
	_, err = os.Stat(path)
	if os.IsNotExist(err) {
		// file does not exist, serve index.html
		http.ServeFile(w, r, filepath.Join(h.staticPath, h.indexPath))
		return
	} else if err != nil {
		// if we got an error (that wasn't that the file doesn't exist) stating the
		// file, return a 500 internal server error and stop
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// otherwise, use http.FileServer to serve the static dir
	http.FileServer(http.Dir(h.staticPath)).ServeHTTP(w, r)
}

// login authenticates the user
// func login(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "cookie-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// Where authentication could be done
// 	if r.FormValue("code") != "code" {
// 		if r.FormValue("code") == "" {
// 			session.AddFlash("Must enter a code")
// 		}
// 		session.AddFlash("The code was incorrect")
// 		err = session.Save(r, w)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		http.Redirect(w, r, "/forbidden", http.StatusFound)
// 		return
// 	}

// 	username := r.FormValue("username")

// 	user := &SSOUser{
// 		Username:      username,
// 		Authenticated: true,
// 	}

// 	session.Values["user"] = user

// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/secret", http.StatusFound)
// }

// logout revokes authentication for a user
// func logout(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "cookie-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	session.Values["user"] = SSOUser{}
// 	session.Options.MaxAge = -1

// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	http.Redirect(w, r, "/", http.StatusFound)
// }

// index serves the index html file
// func index(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "cookie-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	user := getUser(session)
// 	fmt.Println(user)
// 	// tpl.ExecuteTemplate(w, "index.gohtml", user)
// }

// secret displays the secret message for authorized users
// func secret(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "cookie-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	user := getUser(session)

// 	if auth := user.Authenticated; !auth {
// 		session.AddFlash("You don't have access!")
// 		err = session.Save(r, w)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		http.Redirect(w, r, "/forbidden", http.StatusFound)
// 		return
// 	}

// 	// tpl.ExecuteTemplate(w, "secret.gohtml", user.Username)
// }

// func forbidden(w http.ResponseWriter, r *http.Request) {
// 	session, err := store.Get(r, "cookie-name")
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}

// 	// flashMessages := session.Flashes()
// 	err = session.Save(r, w)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	// tpl.ExecuteTemplate(w, "forbidden.gohtml", flashMessages)
// }

// getUser returns a user from session s, on error returns an empty user
func getUser(s *sessions.Session) SSOUser {
	val := s.Values["user"]
	var user = SSOUser{}
	user, ok := val.(SSOUser)
	if !ok {
		return SSOUser{Authenticated: false}
	}
	return user
}

func redirectHandler(w http.ResponseWriter, r *http.Request) {
	provider := "github"

	// In this case we use a map to store our secrets, but you can use dotenv or your framework configuration
	// for example, in revel you could use revel.Config.StringDefault(provider + "_clientID", "") etc.
	providerSecrets := map[string]map[string]string{
		"github": {
			"clientID":     "xx",
			"clientSecret": "yy",
			"redirectURL":  "http://localhost:9000/auth/callback",
		},
	}

	providerScopes := map[string][]string{
		"github": []string{"public_repo"},
	}

	providerData := providerSecrets[provider]
	actualScopes := providerScopes[provider]
	authURL, err := gocial.New().
		Driver(provider).
		Scopes(actualScopes).
		Redirect(
			providerData["clientID"],
			providerData["clientSecret"],
			providerData["redirectURL"],
		)

	// Check for errors (usually driver not valid)
	if err != nil {
		//c.Writer.Write([]byte("Error: " + err.Error()))
		fmt.Printf("Error %v", err)
		return
	}

	http.Redirect(w, r, authURL, http.StatusFound)
}

func callbackHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve query params for state and code
	vals := r.URL.Query()

	states, ok := vals["state"]
	var state string
	if ok {
		if len(states) >= 1 {
			state = states[0] // The first
		}
	}

	codes, ok := vals["code"]
	var code string
	if ok {
		if len(codes) >= 1 {
			code = codes[0] // The first
		}
	}

	// Handle callback and check for errors
	user, token, err := gocial.Handle(state, code)
	if err != nil {
		// c.Writer.Write([]byte("Error: " + err.Error()))
		fmt.Printf("Error %v", err)
		return
	}

	// Print in terminal user information
	fmt.Printf("%#v", token)
	fmt.Printf("%#v", user)

	session, err := store.Get(r, "cookie-name")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	//-------
	userFromSession := getUser(session)
	if auth := userFromSession.Authenticated; !auth {
		fmt.Println("Previously user had no active session!")
	}
	//-------

	ssouser := &SSOUser{
		Username:      user.Username,
		Authenticated: true,
	}

	session.Values["user"] = ssouser

	err = session.Save(r, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// If no errors, show provider name
	// c.Writer.Write([]byte("Hi, " + user.FullName))
	http.Redirect(w, r, "/", http.StatusFound)

}
