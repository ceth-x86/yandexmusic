package yandexmusic

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type MusicPlaylist struct {
	Id                     uint
	Name                   string
	SpotifyId              string
	DeezerId               int64
	YandexId               int64
	SnapshotId             string
	TrackCount             uint
	SkipSyncContent        uint
	Deleted                uint
	LastChanged            *time.Time
	SkipReportWhileSyncing uint
	Manual                 uint
	Tracks                 []*MusicPlaylistTrack
}

type MusicPlaylistTrack struct {
	Id               uint
	PlaylistId       uint
	Name             string
	SpotifyId        string
	DeezerId         int64
	YandexId         string
	Popularity       int
	TrackNumber      int
	Artist           string
	ArtistId         string
	Album            string
	SourceRef        *uint
	Transferred      uint
	Deleted          uint
	SourcePlaylistId *uint
	ReleaseDate      string
}

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

func GetPlaylist(id int64) *MusicPlaylist {

	url := fmt.Sprintf("https://music.yandex.ru/handlers/playlist.jsx?owner=yamusic-new&kinds=%d", id)
	content := get(url)
	extendedPlaylist := ExtendedPlaylist{}
	if err := json.Unmarshal(content, &extendedPlaylist); err != nil {
		log.Fatal(err)
	}

	resultPlaylist := &MusicPlaylist{
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

		resultTrack := &MusicPlaylistTrack{
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
