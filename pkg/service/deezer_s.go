package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
	"net/http"
	"strconv"
)

type DeezerService struct {
	repo repository.User
}

func NewDeezerService(repo repository.User) *DeezerService {
	return &DeezerService{
		repo: repo,
	}
}

type DeezerUserInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DeezerAccessToken struct {
	AccessToken string `json:"access_token"`
	Expires     int    `json:"expires"`
}

type DeezerTrack struct {
	ID     int    `json:"id"`
	Title  string `json:"title"`
	Artist struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"artist"`
	Album struct {
		ID    int    `json:"id"`
		Title string `json:"title"`
	} `json:"album"`
	Type string `json:"type"`
}

func (ds *DeezerService) GetDeezerAccessToken(code string) string {
	postBody, _ := json.Marshal(map[string]string{
		"code": code,
	})
	responseBody := bytes.NewBuffer(postBody)

	// Make POST request
	respAccess, err := http.Post("https://connect.deezer.com/oauth/access_token.php?app_id=491682&secret=3288c76621f0c3a4fa83f3d1cdc1a55f&code="+code+"&output=json", "application/json", responseBody)
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}

	defer respAccess.Body.Close()

	// Read access body
	body, err := ioutil.ReadAll(respAccess.Body)
	if err != nil {
		log.Fatalln(err)
	}

	// Unmarshal access token
	var dzAccess DeezerAccessToken
	err = json.Unmarshal(body, &dzAccess)

	return dzAccess.AccessToken
}

func (ds *DeezerService) CheckDeezerAccessToken(guestID string) bool {
	user, _ := ds.repo.GetUserDB(guestID)

	if user.AccessTokenFind != "" {
		return true
	}

	return false
}

func (ds *DeezerService) GetDeezerUserMusic(guestID string) []model.GeneralMusicStruct {
	accessToken, _ := ds.repo.GetUserDB(guestID)

	var tracks []DeezerTrack
	url := "https://api.deezer.com/user/me/tracks?access_token=" + accessToken.AccessTokenFind

	for {
		var result struct {
			Data    []DeezerTrack `json:"data"`
			Total   int           `json:"total"`
			NextURL *string       `json:"next,omitempty"`
		}

		err := getDZUrl(url, &result)
		if err != nil {
			return nil
		}

		tracks = append(tracks, result.Data...)
		if result.NextURL == nil {
			break
		}

		url = *result.NextURL
	}

	var generalMS []model.GeneralMusicStruct
	for _, track := range tracks {
		generalMS = append(generalMS, model.GeneralMusicStruct{ID: strconv.Itoa(track.ID), ArtistName: track.Artist.Name, SongName: track.Title, AlbumName: track.Album.Title})
	}

	return generalMS
}

func (ds *DeezerService) MoveToDeezer(accessToken string, tracks []model.GeneralMusicStruct) {

}

// Function helper
func getDZUrl(url string, result interface{}) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Non success status code %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}
