package feed

import (
	"fmt"
	"net/http"
	"log"
	"io/ioutil"
	"encoding/json"
	"github.com/l-lin/mr-tracker-api/notification"
	"regexp"
	"strings"
	"strconv"
)

const (
	MANGA_FEEDER_URL = "http://mangafeeder.herokuapp.com/latest.json"
	MANGA_READER_URL = "http://www.mangareader.net/"
)

type Feeds []Feed

type Feed struct {
	Completed bool 		`json:"completed"`
	Title 	  string	`json:"title"`
	Slug 	  string	`json:"slug"`
	Url 	  string	`json:"url"`
	ImageUrl  string	`json:"image_url"`
	New 	  bool		`json:"new"`
	Hot 	  bool		`json:"hot"`
	Chapters  []Chapter `json:"chapters"`
}

type Chapter struct {
	Title string `json:"title"`
	Url   string `json:"url"`
}

func (f Feed) String() string {
	return fmt.Sprintf("Completed = %v, Title = %s, Slug = %s, Url = %s, ImageUrl = %s, New = %v, Hot = %v, Chapters = %v",
		f.Completed, f.Title, f.Slug, f.Url, f.ImageUrl, f.New, f.Hot, f.Chapters)
}

func (c Chapter) String() string {
	return fmt.Sprintf("Title = %s, Url = %s", c.Title, c.Url)
}

// Build the chap number
func (c *Chapter) GetChapNumber() int {
	re := regexp.MustCompile("/([0-9]+)")
	chapStr := re.FindString(c.Url)

	chapStr = strings.Replace(chapStr, "/", "", 1)
	chap, err := strconv.Atoi(chapStr)
	if err != nil {
		log.Printf("[x] Could not parse the chapter number")
		return 0
	}
	return chap
}

// Get the feeds from ${MANGA_FEEDER_URL}
func GetFeeds() Feeds {
	resp, err := http.Get(MANGA_FEEDER_URL)
	if err != nil {
		log.Printf("[x] Could not fetch content of %s. Reason: %s", MANGA_FEEDER_URL, err.Error())
		return nil
	}

	defer resp.Body.Close()
	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("[x] Error reading content of %s. Reason: %s", MANGA_FEEDER_URL, err.Error())
		return nil
	}

	var feeds Feeds
	err = json.Unmarshal(response, &feeds)
	if err != nil {
		log.Printf("[x] Could not unmarshal the rss json feed. Reason: %s", err.Error())
		return nil
	}

	return feeds
}

func GetNewMangaNotifications() []*notification.Notification {
	notifications := make([]*notification.Notification, 0)
	feeds := GetFeeds()
	for _, f := range feeds {
		if f.New {
			n := notification.New()
			n.MangaId = f.Slug
			n.Title = f.Title
			n.Url = f.Url
			n.ImageUrl = f.ImageUrl
			notifications = append(notifications, n)
		}
	}

	return notifications
}
