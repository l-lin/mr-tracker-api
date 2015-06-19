package manga

import (
	"fmt"
	"github.com/l-lin/mr-tracker-api/db"
	"log"
	"bytes"
)

type Manga struct {
	MangaId  string `json:"mangaId"`
	UserId 	 string `json:"userId"`
	LastChap int	`json:"lastChap"`
}

func (m Manga) String() string {
	return fmt.Sprintf("MangaId = %s, UserId = %s, LastChap = %v", m.MangaId, m.UserId, m.LastChap)
}

func New() *Manga {
	return &Manga{}
}

// Check if the given user has mangas
func Exists(userId string) bool {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT CASE WHEN EXISTS(SELECT 1 FROM mangas WHERE user_id = $1) THEN 1 ELSE 0 END", userId)
	var exists int64
	if err := row.Scan(&exists); err != nil {
		log.Printf("[x] Could not check if there is existing mangas for user '%s'. Reason: %s", userId, err.Error())
	}
	return exists == 1;
}

// Copy the default manga to the newly subscribed user
func CopyDefaultFor(userId string) {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("INSERT INTO mangas (manga_id, user_id, last_chap) VALUES ($1, $2, $3), ($4, $5, $6)", "one-piece", userId, 1, "naruto", userId, 1)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not copy the default mangas. Reason: %s", err.Error())
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Fetch the list of mangas
func GetList(userId string) []*Manga {
	mangas := make([]*Manga, 0)
	database := db.Connect()
	defer database.Close()

	rows, err := database.Query(`
	SELECT manga_id, user_id, last_chap
	FROM mangas
	WHERE user_id = $1`,
		userId)
	if err != nil {
		log.Printf("[x] Error when getting the list of mangas. Reason: %s", err.Error())
		return mangas
	}
	for rows.Next() {
		m := toManga(rows)
		if m.IsValid() {
			mangas = append(mangas, m)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("[x] Error when getting the list of mangas. Reason: %s", err.Error())
	}
	return mangas
}

// Fetch all mangas
func GetAll() []*Manga {
	mangas := make([]*Manga, 0)
	database := db.Connect()
	defer database.Close()

	rows, err := database.Query("SELECT manga_id, user_id, last_chap FROM mangas")
	if err != nil {
		log.Printf("[x] Error when getting the list of mangas. Reason: %s", err.Error())
		return mangas
	}
	for rows.Next() {
		m := toManga(rows)
		if m.IsValid() {
			mangas = append(mangas, m)
		}
	}
	if err := rows.Err(); err != nil {
		log.Printf("[x] Error when getting the list of mangas. Reason: %s", err.Error())
	}
	return mangas
}

// Fetch the manga from a given manga id and user id
func Get(mangaId, userId string) *Manga {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT manga_id, user_id, last_chap FROM mangas WHERE manga_id = $1 AND user_id = $2", mangaId, userId)
	return toManga(row)
}

// Save the manga in the db
func (m *Manga) Save() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("INSERT INTO mangas (manga_id, user_id, last_chap) VALUES ($1, $2, $3)", m.MangaId, m.UserId, m.LastChap)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not save the user. Reason: %s", err.Error())
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Update the manga
func (m *Manga) Update() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Printf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("UPDATE mangas SET last_chap = $1 WHERE manga_id = $2 AND user_id = $3", m.LastChap, m.MangaId, m.UserId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not update the manga. Reason: %s", err.Error())
		return
	}

	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Delete a notification
func (m *Manga) Delete() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("DELETE FROM mangas WHERE manga_id = $1 AND user_id = $2", m.MangaId, m.UserId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not delete the manga. Reason: %s", err.Error())
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Delete multiple mangas
func DeleteMultiple(userId string, mangaIds []string) {
	if len(mangaIds) == 0 {
		return
	}
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}

	_, err = tx.Exec(BuildDeleteMultipleQuery(mangaIds), userId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not delete the manga. Reason: %s", err.Error())
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Build the SQL query to delete multiple mangas
func BuildDeleteMultipleQuery(mangaIds []string) string {
	var query bytes.Buffer
	query.WriteString("DELETE FROM mangas WHERE user_id = $1 AND manga_id IN (")
	for index, mangaId := range mangaIds {
		query.WriteString("'")
		query.WriteString(mangaId)
		query.WriteString("'")
		if index < len(mangaIds) - 1 {
			query.WriteString(",")
		}
	}
	query.WriteString(")")
	return query.String()
}

// Check if the manga has valid attributes
func (m *Manga) IsValid() bool {
	return m.MangaId != "" && m.UserId != ""
}

//func (m *Manga) IsNewChap(url string) bool {
//
//}

// Fetch the content of the rows and build a new manga
func toManga(rows db.RowMapper) *Manga {
	m := New()
	err := rows.Scan(
		&m.MangaId,
		&m.UserId,
		&m.LastChap,
	)
	if err != nil {
		log.Printf("[-] Could not scan the manga. Reason: %s", err.Error())
	}
	return m
}
