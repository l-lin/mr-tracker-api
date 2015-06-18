package manga

import (
	"testing"
	"log"
)

func Test_BuildDeleteMultipleQuery(t *testing.T) {
	mangaIds := []string{"one-piece", "naruto"}
	query := BuildDeleteMultipleQuery(mangaIds)
	log.Printf("[-] Query is %s", query)
}
