package notification

import (
	"github.com/l-lin/mr-tracker-api/db"
	_ "github.com/lib/pq"
	"log"
)

type Notification struct {
	NotificationId int    `json:"notificationId"`
	MangaId        string `json:"mangaId"`
	UserId		   string `json:"-"`
	Title          string `json:"title"`
	Url            string `json:"url"`
	ImageUrl	   string `json:"umageUrl"`
}

// Instanciate a new Notification
func New() *Notification {
	return &Notification{}
}

// Fetch all notifications from the db
func GetList(userId string) []*Notification {
	notifications := make([]*Notification, 0)
	database := db.Connect()
	defer database.Close()

	rows, err := database.Query(`
		SELECT notification_id, manga_id, user_id, title, url, image_url
		FROM notifications
		WHERE user_id = $1
	`, userId)
	if err != nil {
		log.Printf("[x] Error when getting the list of feeds. Reason: %s", err.Error())
		return notifications
	}
	for rows.Next() {
		notifications = append(notifications, toNotification(rows))
	}
	if err := rows.Err(); err != nil {
		log.Printf("[x] Error when getting the list of feeds. Reason: %s", err.Error())
	}
	return notifications
}

// Get the notification from a given id
func Get(notificationId int, userId string) *Notification {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow(`
	SELECT notification_id, manga_id, user_id, title, url, image_url
	FROM notifications
	WHERE notification_id = $1 AND user_id = $2`,
		notificationId, userId)
	return toNotification(row)
}

// Save the notification in the database
func (n *Notification) Save() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	row := tx.QueryRow("INSERT INTO notifications (manga_id, user_id, title, url, image_url) VALUES ($1, $2, $3, $4, $5) RETURNING notification_id",
		n.MangaId, n.UserId, n.Title, n.Url, n.ImageUrl)
	var lastId int
	if err := row.Scan(&lastId); err != nil {
		tx.Rollback()
		log.Printf("[x] Could not fetch the notification_id of the newly created notification. Reason: %s", err.Error())
	}
	n.NotificationId = lastId
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Delete a notification
func (n *Notification) Delete() {
	database := db.Connect()
	defer database.Close()
	tx, err := database.Begin()
	if err != nil {
		log.Fatalf("[x] Could not start the transaction. Reason: %s", err.Error())
	}
	_, err = tx.Exec("DELETE FROM notifications WHERE notification_id = $1", n.NotificationId)
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not delete the notification. Reason: %s", err.Error())
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Fetch the content of the rows and build a new default notification
func toNotification(rows db.RowMapper) *Notification {
	n := New()
	err := rows.Scan(
		&n.NotificationId,
		&n.MangaId,
		&n.UserId,
		&n.Title,
		&n.Url,
		&n.ImageUrl,
	)
	if err != nil {
		log.Printf("[-] Could not scan the notification. Reason: %s", err.Error())
	}
	return n
}
