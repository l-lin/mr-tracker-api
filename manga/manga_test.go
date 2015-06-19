package manga

import (
	"testing"
	"log"
	"sort"
)

func Test_BuildDeleteMultipleQuery(t *testing.T) {
	mangaIds := []string{"one-piece", "naruto"}
	query := BuildDeleteMultipleQuery(mangaIds)
	log.Printf("[-] Query to delete multiple mangas is %s", query)
}

func Test_BuildQueryForCopyDefault(t *testing.T) {
	userId := "123"
	query := BuildQueryForCopyDefault(userId)
	log.Printf("[-] Query to copy default mangas is %s", query)
}

func Test_Order(t *testing.T) {
	a := make([]string, len(DEFAULT_MANGAS))
	for _, m := range DEFAULT_MANGAS {
		a = append(a, m)
	}
	sort.Strings(a)
	log.Println(a)
}
