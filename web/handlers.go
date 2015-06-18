package web

import (
	"net/http"
	"encoding/json"
	"fmt"
	sessions "github.com/goincremental/negroni-sessions"
	oauth2 "github.com/goincremental/negroni-oauth2"
	"github.com/codegangsta/negroni"
	"log"
	"github.com/l-lin/mr-tracker-api/manga"
	"github.com/l-lin/mr-tracker-api/notification"
	"github.com/l-lin/mr-tracker-api/feed"
	"github.com/gorilla/mux"
	"io"
	"io/ioutil"
	"strings"
)

// Handler to fetch the list of mangas
func Mangas(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r, nil)

	if !manga.Exists(userId) {
		log.Printf("[-] No mangas found for user %s. Copy the default ones...", userId)
		manga.CopyDefaultFor(userId)
	}

	write(w, http.StatusOK, manga.GetList(userId))
}

// Handler to fetch a manga
func Manga(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	mangaId := vars["mangaId"]

	userId := getUserId(r, nil)

	m := manga.Get(mangaId, userId)
	if m != nil && m.MangaId != "" {
		log.Printf("[-] Found the manga mangaId %s", mangaId)
		write(w, http.StatusOK, m)
		return
	}

	// If we didn't find it, 404
	log.Printf("[-] Could not find the manga with mangaId %s", mangaId)
	write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Manga not Found for mangaId %s", mangaId)})
}

// Handler to save a manga
func SaveManga(w http.ResponseWriter, r *http.Request)  {
	var m manga.Manga
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Printf("[x] Could not read the body. Reason: %s", err.Error())
		write(w, http.StatusInternalServerError, JsonErr{Code: http.StatusInternalServerError, Text: "Could not read the body."})
		return
	}
	if err := r.Body.Close(); err != nil {
		log.Printf("[x] Could not close ready the body. Reason: %s", err.Error())
		write(w, http.StatusInternalServerError, JsonErr{Code: http.StatusInternalServerError, Text: "Could not close the body."})
		return
	}
	if err := json.Unmarshal(body, &m); err != nil {
		// 422: unprocessable entity
		write(w, 422, JsonErr{Code: 422, Text: "Could not parse the given parameter"})
		return
	}
	m.UserId = getUserId(r, nil)

	if !m.IsValid() {
		write(w, http.StatusPreconditionFailed, JsonErr{
			Code: http.StatusPreconditionFailed, Text: "The mangaId should not be empty!",
		})
		return
	}

	log.Printf("[-] Creating new manga %s", m.MangaId)
	m.Save()
	write(w, http.StatusCreated, m)
}

// Handler to update a manga
func UpdateManga(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	mangaId := vars["mangaId"]

	var m manga.Manga
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalf("[x] Could not read the body. Reason: %s", err.Error())
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalf("[x] Could not close ready the body. Reason: %s", err.Error())
	}
	if err := json.Unmarshal(body, &m); err != nil {
		// 422: unprocessable entity
		write(w, 422, JsonErr{Code: 422, Text: "Could not parse the given parameter"})
		return
	}
	m.UserId = getUserId(r, nil)
	m.MangaId = mangaId

	if !m.IsValid() {
		write(w, http.StatusPreconditionFailed, JsonErr{
			Code: http.StatusPreconditionFailed, Text: "The given manga has incorrect attributes",
		})
		return
	}
	log.Printf("[-] Updating manga mangaId %s", mangaId)
	m.Update()
	write(w, http.StatusOK, m)
}

// Handler to delete mangas
func DeleteMangas(w http.ResponseWriter, r *http.Request)  {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Fatalf("[x] Could not read the body. Reason: %s", err.Error())
	}
	if err := r.Body.Close(); err != nil {
		log.Fatalf("[x] Could not close ready the body. Reason: %s", err.Error())
	}
	mangaIdsStr := string(body[:])
	mangaIds := strings.Split(mangaIdsStr, ",")

	log.Printf("[-] Deleting manga with mangaIds %s", mangaIdsStr)
	manga.DeleteMultiple(getUserId(r, nil), mangaIds)
	write(w, http.StatusNoContent, nil)
}

// Handler to fetch the list of notifications
func Notifications(w http.ResponseWriter, r *http.Request) {
	userId := getUserId(r, nil)
	write(w, http.StatusOK, notification.GetList(userId))
}

// Handler to fetch a notification
func Notification(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	notificationId := vars["notificationId"]

	n := notification.Get(notificationId)
	if n != nil && n.NotificationId != "" {
		log.Printf("[-] Found the notification id %s", notificationId)
		write(w, http.StatusOK, n)
		return
	}

	// If we didn't find it, 404
	log.Printf("[-] Could not find the notification id %s", notificationId)
	write(w, http.StatusNotFound, JsonErr{Code: http.StatusNotFound, Text: fmt.Sprintf("Notification not Found for notificationId %s", notificationId)})
}

// Handler to delete a notification
func DeleteNotification(w http.ResponseWriter, r *http.Request)  {
	vars := mux.Vars(r)
	notificationId := vars["notificationId"]
	n := notification.New()
	n.NotificationId = notificationId

	log.Printf("[-] Deleting notification id %s", notificationId)
	n.Delete()
	write(w, http.StatusNoContent, nil)
}

// Handler to fetch the list of mangas
func NewMangas(w http.ResponseWriter, r *http.Request) {
	write(w, http.StatusOK, feed.GetNewMangaNotifications())
}

// Handler to sign in Google account
func SignIn(w http.ResponseWriter, r *http.Request) {
	userInfo, oauthT, err := getUserInfo(r)
	if err == nil {
		// Save or updating with fresh data of the user
		saveOrUpdateUser(userInfo, oauthT)

		// Save the userId in the session
		s := sessions.GetSession(r)
		s.Set(SESSION_USER_ID, userInfo.Id)
	}

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
		http.Redirect(rw, r, oauth2.PathLogout, http.StatusFound)
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
