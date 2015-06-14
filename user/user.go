package user

import (
	"github.com/l-lin/mr-tracker-api/db"
	_ "github.com/lib/pq"
	"log"
	"fmt"
	"os"
	"encoding/json"
	"bytes"
	"net/http"
	"io/ioutil"
	"time"
)

const oauth2RefreshEndPoint  = "https://www.googleapis.com/oauth2/v3/token"

// The feed
type User struct {
	UserId       	string `json:"-"`
	RefreshToken 	string `json:"-"`
	Email		 	string `json:"-"`
	LastConnection 	time.Time `json:"-"`
}

func (t User) String() string {
	return fmt.Sprintf("UserId = %s, RefreshToken = %s, Emaio = %s", t.UserId, t.RefreshToken)
}

type ResfreshTokenConfig struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	GrantType    string `json:"grant_type"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenResponse struct {
	AccessToken string 	`json:"access_token"`
	ExpiresIn   int 	`json:"expires_in"`
	TokenType   string 	`json:"token_type"`
}

// Instanciate a new feed
func New() *User {
	return &User{}
}

// Check if the given userId already exist
func Exists(userId string) bool {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT CASE WHEN EXISTS(SELECT 1 FROM tokens WHERE user_id = $1) THEN 1 ELSE 0 END", userId)
	var exists int64
	if err := row.Scan(&exists); err != nil {
		log.Printf("[x] Could not check if there is existing user for user '%s'. Reason: %s", userId, err.Error())
	}
	return exists == 1;
}

// Get the User from a given userId
func Get(userId string) *User {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT user_id, refresh_token, email, last_connection FROM tokens WHERE user_id = $1", userId)
	return toUser(row)
}

// Get the Novel given an novelId
func GetByAccessToken(accessToken string) *User {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT user_id, refresh_token, email, last_connection FROM tokens WHERE access_token = $1", accessToken)
	return toUser(row)
}

// Save the token in the database
func (user *User) Save() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("INSERT INTO users (user_id, refresh_token, email, last_connection) VALUES ($1, $2, $3, $4)", user.UserId, user.RefreshToken, user.Email, user.LastConnection)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not save the user. Reason: %s", err.Error())
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Update the user
func (user *User) Update() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("UPDATE users SET refresh_token = $1, email = $2, last_connection = $3 WHERE user_id = $4", user.RefreshToken, user.UserId, user.Email, user.LastConnection)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not update the user. Reason: %s", err.Error())
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Refresh the given user
func (user *User) Refresh() bool {
	c := &ResfreshTokenConfig{os.Getenv("GOOGLE_CLIENT_ID"), os.Getenv("GOOGLE_CLIENT_SECRET"), "refresh_token", user.RefreshToken}
	buf, _ := json.Marshal(c)
	body := bytes.NewBuffer(buf)
	r, err := http.Post(oauth2RefreshEndPoint, "application/json", body)
	if err != nil {
		log.Printf("[x] Could not refresh the user. Reason: %s", err.Error())
		return false
	}
	defer r.Body.Close()
	response, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[x] Error reading content of %s. Reson: %s", oauth2RefreshEndPoint, err.Error())
		return false
	}
	var oauthResponse RefreshTokenResponse
	if err := json.Unmarshal(response, &oauthResponse); err != nil {
		log.Printf("[x] Could not read the JSON of the response after refreshing the user. Reason: %s", err.Error())
		return false
	}
	return true
}

// Fetch the content of the rows and build a new user
func toUser(rows db.RowMapper) *User {
	user := New()
	err := rows.Scan(
		&user.UserId,
		&user.RefreshToken,
		&user.Email,
		&user.LastConnection,
	)
	if err != nil {
		log.Printf("[-] Could not scan the user. Reason: %s", err.Error())
	}
	return user
}