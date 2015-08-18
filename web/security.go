package web

import (
	"github.com/l-lin/mr-tracker-api/user"
	"github.com/codegangsta/negroni"
	sessions "github.com/goincremental/negroni-sessions"
	oauth2 "github.com/goincremental/negroni-oauth2"
	"os"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"time"
	"fmt"
)

const (
	SESSION_USER_ID = "user_id"
	googleUserInfoEndPoint = "https://www.googleapis.com/oauth2/v1/userinfo"
)

// The user info for Google account
type UserInfo struct {
	Id 		string
	Email 	string
	Picture string
}

// Returns a new Negroni middleware using Google OAuth2
func NewOAuth() negroni.Handler {
	return oauth2.Google(&oauth2.Config{
	ClientID: os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL: os.Getenv("GOOGLE_REDIRECT_URI"),
	Scopes: []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
})
}

// Wrap the HandlerFunc by checking if the user is indeed authenticated
func WrapWithCheckAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		oauthT := oauth2.GetToken(r)
		if oauthT == nil {
			reject(w)
		} else {
			if !oauthT.Valid() {
				log.Printf("[-] The oauthToken is not valid")
				userId := getUserId(r, saveOrUpdateUser)

				u := user.Get(fmt.Sprintf("%v", userId))
				if u != nil {
					log.Printf("[-] Refreshing the token %s", u.RefreshToken)
					if u.Refresh() {
						handlerFunc.ServeHTTP(w, r)
					}
				}
			} else {
				log.Printf("[-] The oauthToken is valid")
				userId := getUserId(r, nil)
				u := user.Get(fmt.Sprintf("%v", userId))
				u.LastConnection = time.Now()
				u.Update()

				handlerFunc.ServeHTTP(w, r)
			}
		}
	}
}

// Get the user id.
// First fetch it from the session
// If not present, then fetch it from Google service
func getUserId(r *http.Request, callback func (UserInfo, oauth2.Tokens)) string {
	s := sessions.GetSession(r)
	userId := s.Get(SESSION_USER_ID)
	// If userId not found, then fetch the info from Google
	if userId == nil {
		userInfo, oauthT, err := getUserInfo(r)
		if err == nil {
			userId = userInfo.Id
			s.Set(SESSION_USER_ID, userId)
			if callback != nil {
				// Save or updating with fresh data of the user
				callback(userInfo, oauthT)
			}
		}
	} else {
		log.Printf("[-] Updating last connection date for userId %v", userId)
		user.UpdateLastConnection(fmt.Sprintf("%v", userId))
	}
	return fmt.Sprintf("%v", userId)
}

// Reject the request by sending a HTTP 401
func reject(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(JsonErr{Code: http.StatusUnauthorized, Text: "You are not authenticated!"}); err != nil {
		log.Fatalf("[x] Error when encoding the json. Reason: %s", err.Error())
	}
}

// Save or update the given user info
func saveOrUpdateUser(userInfo UserInfo, oauthT oauth2.Tokens) {
	if !user.Exists(userInfo.Id) {
		u := user.New()
		u.UserId = userInfo.Id
		u.Email = userInfo.Email
		u.Picture = userInfo.Picture
		u.LastConnection = time.Now()
		u.RefreshToken = oauthT.Refresh()
		log.Printf("[-] Saving user %v", u)
		u.Save()
	} else {
		u := user.Get(userInfo.Id)
		u.Email = userInfo.Email
		u.Picture = userInfo.Picture
		u.LastConnection = time.Now()
		if oauthT.Refresh() != "" {
			log.Printf("[-] The refresh token is not empty => the user had revoked the permissions")
			u.RefreshToken = oauthT.Refresh()
		}
		log.Printf("[-] Updating the user %v", u)
		u.Update()
	}
}

// Get the user ID from a given token.
// It will make a GET request to https://www.googleapis.com/oauth2/v1/userinfo?access_token=...
func getUserInfo(r *http.Request) (UserInfo, oauth2.Tokens, error) {
	var userInfo UserInfo
	oauthT := oauth2.GetToken(r)
	if oauthT == nil || !oauthT.Valid() {
		log.Printf("[x] The user is not authenticated yet!")
	}
	accessToken := oauthT.Access()

	log.Printf("[-] Getting the user id from access token %s", accessToken)
	endPoint := googleUserInfoEndPoint + "?access_token=" + accessToken
	resp, err := http.Get(endPoint)
	if err != nil {
		log.Printf("[x] Could not find the user info with token %s. Reason: %s", accessToken, err.Error())
		return userInfo, oauthT, err
	}
	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[x] Error reading content of %s. Reson: %s", endPoint, err.Error())
		return userInfo, oauthT, err
	}
	err = json.Unmarshal(response, &userInfo)
	if err != nil {
		log.Printf("[x] Could not unmarshal the user info. Reason: %s", err.Error())
		return userInfo, oauthT, err
	}

	return userInfo, oauthT, nil
}
