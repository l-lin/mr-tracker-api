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
	Id 		string `json:"id"`
	Email 	string `json:"email"`
	Picture string `json:"picture"`
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
				s := sessions.GetSession(r)
				userId := s.Get(SESSION_USER_ID)

				// If userId not found, then fetch the info from Google
				if userId == nil {
					userInfo, oauthT, err := GetUserInfo(r)
					if err == nil {
						userId = userInfo.Id
						saveOrUpdateUser(userInfo, oauthT)
					}
				}

				user := user.Get(fmt.Sprintf("%v", userId))
				if user != nil {
					log.Printf("[-] Refreshing the token %s", user.RefreshToken)
					if user.Refresh() {
						handlerFunc.ServeHTTP(w, r)
					}
				}
			} else {
				handlerFunc.ServeHTTP(w, r)
			}
		}
	}
}

func reject(w http.ResponseWriter) {
	w.WriteHeader(http.StatusUnauthorized)
	if err := json.NewEncoder(w).Encode(JsonErr{Code: http.StatusUnauthorized, Text: "You are not authenticated!"}); err != nil {
		log.Fatalf("[x] Error when encoding the json. Reason: %s", err.Error())
	}
}

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
			// If the refresh is not empty => the user had revoked the permissions => we have to update the token
			u.RefreshToken = oauthT.Refresh()
		}
		log.Printf("[-] Updating the user %v", u)
		u.Update()
	}
}

// Get the user ID from a given token.
// It will make a GET request to https://www.googleapis.com/oauth2/v1/userinfo?access_token=...
func GetUserInfo(r *http.Request) (UserInfo, oauth2.Tokens, error) {
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
