package web

import (
	"net/http"
	"encoding/json"
	"fmt"
	sessions "github.com/goincremental/negroni-sessions"
	oauth2 "github.com/goincremental/negroni-oauth2"
	"github.com/codegangsta/negroni"
	"log"
)

// Handler to sign in Google account
func SignIn(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are now authenticated! You can close this tab.")
}

// Handler that redirects user to the login page
func SignOut() negroni.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		log.Println("[-] Signing out...")

		s := sessions.GetSession(r)
		s.Delete(SESSION_USER_ID)

		// Set token to null to avoid redirection loop
		oauth2.SetToken(r, nil)
		http.Redirect(rw, r, oauth2.PathLogout + "?next=/authTest", http.StatusFound)
	}
}

// This Handler is used only to check if the user is indeed authenticated
func AuthTest(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "You are now authenticated! You can close this tab.")
}

// Write the response in JSON Content-type
func write(w http.ResponseWriter, status int, n interface{}) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(status)

	if n != nil {
		if err := json.NewEncoder(w).Encode(n); err != nil {
			panic(err)
		}
	}
}
