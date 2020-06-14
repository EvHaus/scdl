package soundcloud

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"

	"github.com/bogem/id3v2"
	"github.com/imthaghost/scdl/mp3"
)

type audioLink struct {
	URL string `json:"url"`
}

// TODO: implement tests
// ExtractSong queries the SoundCloud api and receives a m3u8 file, then binds the segments received into a .mp3 file
func ExtractSong(url string) {

	// request to user inputed SoundCloud URL
	resp, err := http.Get(url)
	if err != nil {
		log.Fatalln(err)
	}
	// response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// parse the response data to grab the song name
	songname := GetTitle(body)

	// parse the response data to grab the artwork URL
	artworkURL := GetArtwork(body)

	// parse the response data and make a reqeust to receive clien_id embedded in the javascript
	clientID := GetClientID(body)

	// TODO: probably cleaner to just move this request into the GetArtwork function
	// request to artwork url to download image data
	artworkresp, err := http.Get(artworkURL)
	if err != nil {
		log.Fatalln(err)
	}
	// image data
	image, err := ioutil.ReadAll(artworkresp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// TODO improve pattern for finding encrypted string ID
	var re = regexp.MustCompile(`https:\/\/api-v2.*\/stream\/hls`) // pattern for finding encrypted string ID
	// TODO not needed if encrypted string ID regex pattern is improved
	var ree = regexp.MustCompile(`.+?(stream)`) // pattern for finding stream URL

	streamURL := re.FindString(string(body)) // stream URL

	baseURL := ree.FindString(streamURL) // baseURL ex: https://api-v2.soundcloud.com/media/soundcloud:tracks:816595765/0ad937d5-a278-4b36-b128-220ac89aec04/stream

	// TODO: replace with format string instead of concatenation
	requestURL := baseURL + "/hls?client_id=" + clientID // API query string ex: https://api-v2.soundcloud.com/media/soundcloud:tracks:805856467/ddfb7463-50f1-476c-9010-729235958822/stream/hls?client_id=iY8sfHHuO2UsXy1QOlxthZoMJEY9v0eI

	// query API
	r, e := http.Get(requestURL)
	if err != nil {
		log.Fatalln(e)
	}

	// API response returns a m3u8 URL
	m3u8Reponse, er := ioutil.ReadAll(r.Body)
	if er != nil {
		log.Fatalln(er)
	}

	var a audioLink

	// unmarshal json data from response
	audioerr := json.Unmarshal(m3u8Reponse, &a)
	if er != nil {
		panic(audioerr)
	}

	// merege segments
	mp3.Merge(a.URL, songname)

	// replace empty cover image with SoundCloud artwork
	tag, err := id3v2.Open(songname+".mp3", id3v2.Options{Parse: true})
	if tag == nil || err != nil {
		log.Fatal("Error while opening mp3 file: ", err)
	}
	pic := id3v2.PictureFrame{
		Encoding:    id3v2.EncodingUTF8,
		MimeType:    "image/jpeg",
		PictureType: id3v2.PTFrontCover,
		Description: "Front cover",
		Picture:     image,
	}
	tag.AddAttachedPicture(pic)
	tag.Save()
}
