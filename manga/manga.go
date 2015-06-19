package manga

import (
	"fmt"
	"github.com/l-lin/mr-tracker-api/db"
	"log"
	"bytes"
)

var (
	DEFAULT_MANGAS = [...]string{"flying-witch", "shinmai-fukei-kiruko-san", "team-medical-dragon",
	"the-god-of-high-school", "girls-of-the-wilds", "the-breaker-new-waves", "onepunch-man", "dendrobates", "usogui",
	"assassination-classroom", "btooom", "nisekoi", "karate-shoukoushi-kohinata-minoru", "shibatora", "yuureitou",
	"shigatsu-wa-kimi-no-uso", "renai-boukun", "tonari-no-seki-kun", "saiki-kusuo-no-psi-nan", "bonnouji", "black-joke",
	"demon-lord-at-work", "world-trigger", "kotoura-san", "nanatsu-no-taizai", "no-game-no-life", "mahouka-koukou-no-rettousei",
	"shokugeki-no-soma", "uwagaki", "saijou-no-meii", "natsu-no-zenjitsu", "maoyuu-maou-yuusha", "hanza-sky", "horimiya",
	"haikyu", "noblesse", "buyuden", "addicted-to-curry", "feng-shen-ji", "one-piece", "naruto", "samurai-usagi",
	"cavalier-of-the-abyss", "toriko", "again", "evergreen", "hungry-joker", "beelzebub", "kissxsis", "kaiji",
	"celestial-clothes", "city-of-darkness", "rex-fabula", "gamaran", "hammer-session", "datenshi-gakuen-debipara",
	"taiyou-no-ie", "alyosha", "dear-only-you-dont-know", "magician", "rinne-no-lagrange-akatsuki-no-memoria", "mario", "kokoro-connect",
	"claymore", "sket-dance", "darwins-game", "heroes-of-the-spring-and-autumn", "the-lawless", "yamada-kun-to-7-nin-no-majo",
	"historys-strongest-disciple-kenichi", "yotsubato", "703x", "kingdom", "baby-steps", "hitsugime-no-chaika",
	"boku-to-kanojo-no-renai-mokuroku", "soukyuu-no-lapis-lazuli", "renai-kaidan-sayoko-san", "i-am-a-hero", "ubel-blatt",
	"kamisama-no-iutoori-ii", "killer-stall", "shingeki-no-kyojin", "highschool-dxd", "that-future-is-a-lie",
	"until-death-do-us-part", "gantz", "ill-make-you-into-an-otaku-so-make-me-into-a-riajuu", "denpa-kyoushi",
	"the-blue-eyed-material", "aku-no-hana", "hakaijuu", "soul-catchers", "hinamatsuri", "wagatsuma-san-wa-ore-no-yome",
	"chunchu", "blood-lad", "sankarea", "cradle-of-monsters", "ocha-nigosu", "shuujin-riku",
	"kigurumi", "kanojo-ga-flag-o-oraretara", "imawa-no-kuni-no-alice", "doubt!", "dragons-rioting", "shiryoku-kensa",
	"ability", "witch-hunter", "smoky-bb", "new-prince-of-tennis", "hajime-no-ippo", "sprite", "one-shot-demon-king",
	"oukoku-no-ko", "tekken-chinmi-legends", "the-devil-king-is-bored", "mix", "not-lives", "mutou-black", "gakusen-toshi-asterisk",
	"green-boy", "shurabara", "donyatsu", "wallman", "piano-no-mori", "bamboo-blade", "yankee-kun-to-megane-chan",
	"tail-star", "shinonome-yuuko-wa-tanpen-shousetsu-o-aishite-iru", "mercenary-maruhan", "masamune-kun-no-revenge",
	"honorable-baek-dong-soo", "the-heroic-legend-of-arslan-arakawa-hiromu", "81-diver", "ginga-patrol-jako",
	"naze-toudouin-masaya-16-sai-wa-kanojo-ga-dekinai-no-ka", "classmate-kamimura-yuuka-wa-kou-itta",
	"tetsugaku-letra", "ajin", "ties-of-compassion", "nobunaga-no-chef", "kono-bijutsubu-ni-wa-mondai-ga-aru",
	"dice-the-cube-that-changes-everything", "teppu", "extra-existence", "boku-ni-koi-suru-mechanical", "birdmen",
	"ultimate-special-high-school", "jitsu-wa-watashi-wa",
	"dr-duo", "ane-log", "suicide-island", "mysterious-girlfriend-x", "wolf-mary",
	"wakabayashi-toshiyas-4-koma-collection", "moteki", "hime-doll", "rain", "rebirth-knight", "naqua-den", "uq-holder",
	"te-to-kuchi", "terra-formars", "saito-kun-wa-chounouryokusha-rashii", "sword-art-online",
	"shinigami-sama-to-4-nin-no-kanojo", "the-gamer", "lets-lagoon", "dame-na-watashi-ni-koishite-kudasai",
	"gu-fang-bu-zi-shang", "apocalypse-no-toride", "hachi", "adventure-of-sinbad-prototype", "magi",
	"black-lagoon", "lasboss-x-hero", "danshi-koukousei-no-nichijou",
	"tokyo-esp", "life-ho", "accel-world", "dungeon-ni-deai-wo-motomeru-no-wa-machigatte-iru-darou-ka-gaiden-sword-oratoria", "koe-no-katachi",
	"unbreakable-machine-doll", "chong-tai-ji", "akame-ga-kiru",
	"gordian-knot", "14-sai-no-koi", "anagle-mole", "toukyou-kushu", "iris-zero", "dolls-code",
	"kokushi-musou", "god-of-bath", "koutetsu-no-hanappashira", "zai-x-10", "american-ghost-jack", "the-kenseis-calligraphy",
	"soutaisei-moteron", "dragon-collection", "boku-wa-tomodachi-ga-sukunai",
	"hangyaku-no-kagetsukai", "50-million-km", "glory-hill", "iron-knight",
	"gigantomakhia", "mahou-gyoushounin-roma", "eldlive", "chaos-rings",
	"sakurasaku-shoukougun", "silver-spoon", "lucifer-no-migite",
	"chousuinou-kei-makafushigi-jiken-file", "wake-up-deadman", "cross-x-regalia", "hallelujah-overdrive", "sousei-no-onmyouji",
	"aizawa-san-zoushoku", "clockwork-planet", "to-aru-kagaku-no-accelerator", "mob-psycho-100", "looking-for-a-father",
	"all-you-need-is-kill", "hamatora", "destro-246",
	"high-score-girl", "area-d-inou-ryouiki", "bakuon", "once-again", "prison-school", "real-account", "inu-yashiki",
	"bocchiman", "illegal-rare", "stealth-symphony", "inu-to-hasami-wa-tsukaiyou", "yasashii-sekai-no-tsukurikata", "azumi",
	"ballroom-e-youkoso", "stretch", "shinazu-no-ryouken", "sentou-jousai-masurawo",
	"fujimi-lovers", "mokushiroku-alice", "berserk", "big-order", "area-no-kishi", "second-brain",
	"yahari-ore-no-seishun-rabukome-wa-machigatte-iru-mougenroku", "chihayafuru", "girls-go-around", "rakia", "itoshi-no-karin",
	"combat-continent", "bocchi-na-bokura-no-renai-jijou", "hunter-x-hunter", "chrono-monochrome", "gojikanme-no-sensou-home-sweet-home",
	"remonster", "bambino-secondo", "shadow---super-human-assistance-department-office-worker",
	"world-customize-creator", "tsuyokute-new-saga", "aho-girl", "metallica-metalluca",
	"vinland-saga", "waga-na-wa-umishi", "roboticsnotes", "hero",
	"americano-exodus", "colors-of-the-wind", "a-fairytale-for-the-demon-lord",
	"tale-of-eun-aran", "mahouka-koukou-no-yuutousei", "stravaganza---isai-no-hime", "4-god-ranger",
	"gunjou-senki", "3-gatsu-no-lion", "saike-mata-shite-mo", "yuusha-gojo-kumiai-kouryuugata-keijiban",
	"yawara", "mayonaka-no-x-giten", "is-reiroukan-still-alive",
	"tate-no-yuusha-no-nariagari", "rakudai-kishi-no-cavalry",
	"saijaku-muhai-no-bahamut", "delusional-boy", "hare-kon", "supinamarada",
	"kanata-no-togabito-ga-tatakau-riyuu", "mongrel", "tokyo-dted", "jyoshikausei",
	"gekkan-shoujo-nozaki-kun",
}
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

// Check if the user has the given manga
func HasManga(userId, mangaId string) bool {
	database := db.Connect()
	defer database.Close()

	row := database.QueryRow("SELECT CASE WHEN EXISTS(SELECT 1 FROM mangas WHERE user_id = $1 AND manga_id = $2) THEN 1 ELSE 0 END", userId, mangaId)
	var exists int64
	if err := row.Scan(&exists); err != nil {
		log.Printf("[x] Could not check if there is existing mangas for user '%s' and mangaId '%s'. Reason: %s", userId, mangaId, err.Error())
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
		return
	}
	_, err = tx.Exec(BuildQueryForCopyDefault(userId))
	if err != nil {
		tx.Rollback()
		log.Printf("[x] Could not copy the default mangas. Reason: %s", err.Error())
		return
	}
	if err := tx.Commit(); err != nil {
		log.Printf("[x] Could not commit the transaction. Reason: %s", err.Error())
	}
}

// Build the query to copy the default mangas
func BuildQueryForCopyDefault(userId string) string {
	var query bytes.Buffer
	query.WriteString("INSERT INTO mangas (manga_id, user_id, last_chap) VALUES ")
	for index, mangaId := range DEFAULT_MANGAS {
		query.WriteString(fmt.Sprintf("('%s', '%s', %d)", mangaId, userId, 1))
		if index < len(DEFAULT_MANGAS) - 1 {
			query.WriteString(",")
		}
	}

	return query.String()
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
		return
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
		return
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
		return
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
