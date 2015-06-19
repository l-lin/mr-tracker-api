package feed

import (
	"github.com/robfig/cron"
	"log"
	"github.com/l-lin/mr-tracker-api/manga"
	"github.com/l-lin/mr-tracker-api/notification"
)

// Cron to fetch all the rss content
func NewCronRss() *cron.Cron {
	c := cron.New()
	c.AddFunc("0 */1 * * * *", FillNotifications)
	return c
}

func FillNotifications() {
	log.Printf("[-] CRON - Starting to fill the table notifications...")
	feeds := GetFeeds()

	if len(feeds) > 0 {
		mangas := manga.GetAll()
		notifications := make([]*notification.Notification, 0)
		for _, m := range mangas {
			notifications = append(notifications, getNotificationsFromFeeds(feeds, m)...)
		}
		log.Printf("[-] CRON - Saved %d notifications", len(notifications))
	} else {
		log.Printf("[-] CRON - There are no notifications to save")
	}

	log.Printf("[-] CRON - Finished filling the table notifications...")
}

func getNotificationsFromFeeds(feeds Feeds, m *manga.Manga) []*notification.Notification {
	notifications := make([]*notification.Notification, 0)
	for _, f := range feeds {
		notifications = append(notifications, getNotificationsFromFeed(f, m)...)
	}
	return notifications
}

func getNotificationsFromFeed(f Feed, m *manga.Manga) []*notification.Notification {
	notifications := make([]*notification.Notification, 0)
	if m.MangaId == f.Slug {
		for _, chap := range f.Chapters {
			if chap.GetChapNumber() > m.LastChap {
				log.Printf("[-] CRON - Found match for manga %s. New chapter is %d", m.MangaId, chap.GetChapNumber())
				n := notification.New()
				n.UserId = m.UserId
				n.MangaId = m.MangaId
				n.Title = chap.Title
				n.Url = chap.Url
				n.ImageUrl = f.ImageUrl
				n.Save()
				notifications = append(notifications, n)

				// Update the manga last chapter
				m.LastChap = chap.GetChapNumber()
				m.Update()
			}
		}
	}
	return notifications
}
