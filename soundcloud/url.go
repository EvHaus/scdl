package soundcloud

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/PuerkitoBio/goquery"
)

// GetSongName uses regex to grab the song name from the URL string
func GetSongName(url string) string {
	/*
		>>> https://soundcloud.com/chillem-637935049/1400-999-freestyle
		<<< 1400-999-freestyle

		>>> https://soundcloud.com/a-boogie-wit-da-hoodie/demons-and-angels
		<<< demons-and-angels
	*/

	var re = regexp.MustCompile(`([^\/]*)$`)

	name := re.FindString(url)

	return name
}

// IsSong looks at a given SoundCloud URL and determine if the URL is a song or not
func IsSong(url string) bool {
	/*
		>>> https://soundcloud.com/uiceheidd
		<<< fasle

		>>> https://soundcloud.com/uiceheidd/tell-me-you-love-me
		<<< true
	*/
	var song bool
	// request to soundcloud url
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	// response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	// turn byte array into reader
	data := bytes.NewReader(body)
	// create new document from reader
	doc, err := goquery.NewDocumentFromReader(data)
	if err != nil {
		panic(err)
	}
	x := doc.Find("meta[property='soundcloud:like_count']").First()
	_, exists := x.Attr("content")
	if exists {
		song = true
	} else {
		song = false
	}
	return song
}
