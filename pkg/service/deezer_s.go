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
	"net/url"
	"regexp"
	"strconv"
	"time"
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

type ResultDzSearch struct {
	Data []struct {
		ID int `json:"id"`
	} `json:"data"`
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
	urlS := "https://api.deezer.com/user/me/tracks?access_token=" + accessToken.AccessTokenFind

	for {
		var result struct {
			Data    []DeezerTrack `json:"data"`
			Total   int           `json:"total"`
			NextURL *string       `json:"next,omitempty"`
		}

		err := getDZUrl(urlS, &result)
		if err != nil {
			return nil
		}

		tracks = append(tracks, result.Data...)
		if result.NextURL == nil {
			break
		}

		urlS = *result.NextURL
	}

	var generalMS []model.GeneralMusicStruct
	for _, track := range tracks {
		generalMS = append(generalMS, model.GeneralMusicStruct{ID: strconv.Itoa(track.ID), ArtistName: track.Artist.Name, SongName: track.Title, AlbumName: track.Album.Title})
	}

	return generalMS
}

func (ds *DeezerService) MoveToDeezer(accessToken string, tracks []model.GeneralMusicStruct) []string {
	var found []int
	var notFound []string

	// Search tracks
	for _, track := range tracks {
		reg := regexp.MustCompile(`(?i)\(.*|feat.*|- feat.*|- with.*`)

		searchString := fmt.Sprintf("%s %s", track.ArtistName, reg.Split(track.SongName, -1)[0])
		artistName := url.PathEscape(track.ArtistName)
		shortSongName := url.PathEscape(reg.Split(track.SongName, -1)[0])
		shortAlbumName := url.PathEscape(reg.Split(track.AlbumName, -1)[0])

		searchUrl := "https://api.deezer.com/search?order=RANKING&q=artist:" + "\"" + artistName + "\"" + "track:" + "\"" + shortSongName + "\"" + "album:" + "\"" + shortAlbumName + "\"" + "&limit=1"
		var result ResultDzSearch
		getDZUrl(searchUrl, &result)

		// Deep music search
		if len(result.Data) == 0 {
			time.Sleep(1 * time.Second)
			searchUrl = "https://api.deezer.com/search?order=RANKING&q=artist:" + "\"" + artistName + "\"" + "track:" + "\"" + url.PathEscape(track.SongName) + "\"" + "&limit=1"
			getDZUrl(searchUrl, &result)

			if len(result.Data) == 0 {
				time.Sleep(1 * time.Millisecond)
				searchUrl = "https://api.deezer.com/search/track?order=RANKING&q=" + url.PathEscape(searchString) + "&limit=1"
				getDZUrl(searchUrl, &result)

				if len(result.Data) == 0 {
					time.Sleep(1 * time.Second)
					searchUrl = "https://api.deezer.com/search/track?order=RANKING&q=track:" + "\"" + shortSongName + "\"" + "album:" + "\"" + url.PathEscape(track.AlbumName) + "\"" + "&limit=1"
					getDZUrl(searchUrl, &result)

					if len(result.Data) == 0 {
						notFound = append(notFound, searchString, track.AlbumName)
					} else {
						found = append(found, result.Data[0].ID)
					}
				} else {
					found = append(found, result.Data[0].ID)
				}
			} else {
				found = append(found, result.Data[0].ID)
			}
		} else {
			found = append(found, result.Data[0].ID)
		}
	}

	// Move tracks
	if len(found) > 0 {
		c := 0
		for _, id := range found {
			c += 1

			time.Sleep(1 * time.Second)
			resp, err := http.Post("https://api.deezer.com/user/me/tracks?access_token="+accessToken+"&track_id="+strconv.Itoa(id), "application/x-www-form-urlencoded", nil)

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println(resp.StatusCode)
			fmt.Println(c)

			resp.Body.Close()
		}
	}

	fmt.Println("NOTFOUND", notFound, len(notFound))
	return notFound
}

// Function helper
func getDZUrl(url string, result interface{}) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
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
