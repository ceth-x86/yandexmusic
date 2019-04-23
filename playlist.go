package yandex

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/demas/music/models"
	"github.com/demas/music/yandexclient/yandexmodels"
)

type Playlist struct {
	Id          uint    `json:"kind,omitempty"`
	Name        string  `json:"title,omitempty"`
	Description string  `json:"description,omitempty"`
	Tracks      []Track `json:"tracks,omitempty"`
}

type Track struct {
	Id      string   `json:"id,omitempty"`
	Title   string   `json:"title,omitempty"`
	Artists []Artist `json:"artists,omitempty"`
	Albums  []Album  `json:"albums,omitempty"`
}

type Artist struct {
	Id   uint   `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

type Album struct {
	Id    uint   `json:"id,omitempty"`
	Title string `json:"title,omitempty"`
}

type ExtendedPlaylist struct {
	Playlist Playlist `json:"playlist,omitempty"`
}

func get(url string) []byte {

	response, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	defer response.Body.Close()

	return contents
}

func GetPlaylist(id int64) *models.Playlist {

	url := fmt.Sprintf("https://music.yandex.ru/handlers/playlist.jsx?owner=yamusic-new&kinds=%d", id)
	content := get(url)
	extendedPlaylist := yandexmodels.ExtendedPlaylist{}
	if err := json.Unmarshal(content, &extendedPlaylist); err != nil {
		log.Fatal(err)
	}

	resultPlaylist := &models.Playlist{
		Name:            extendedPlaylist.Playlist.Name,
		YandexId:        id,
		SkipSyncContent: 0,
		Deleted:         0,
	}

	for _, serviceTrack := range extendedPlaylist.Playlist.Tracks {

		album := ""
		artist := ""

		if len(serviceTrack.Albums) > 0 {
			album = serviceTrack.Albums[0].Title
		}

		if len(serviceTrack.Artists) > 0 {
			artist = serviceTrack.Artists[0].Name
		}

		resultTrack := &models.PlaylistTrack{
			Name:        serviceTrack.Title,
			YandexId:    serviceTrack.Id,
			Popularity:  0,
			TrackNumber: 0,
			Artist:      artist,
			Album:       album,
			Transferred: 0,
			Deleted:     0,
		}

		resultPlaylist.Tracks = append(resultPlaylist.Tracks, resultTrack)
	}

	return resultPlaylist
}
