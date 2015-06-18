package feed

import "testing"

func Test_GetChapNumber(t *testing.T) {
	c := &Chapter{
		Title: "Seifuku Aventure - Chemical Reaction Of High School Students 10",
		Url: "http://www.mangareader.net/seifuku-aventure-chemical-reaction-of-high-school-students/10",
	}
	chap := c.GetChapNumber()
	if chap != 10 {
		t.Fail()
	}

	c = &Chapter{
		Title: "Dagashi Kashi 13",
		Url: "http://www.mangareader.net/dagashi-kashi/13",
	}
	chap = c.GetChapNumber()
	if chap != 13 {
		t.Fail()
	}
}
