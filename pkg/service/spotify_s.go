package service

import (
	"encoding/json"
	"fmt"
	"github.com/gofiber/websocket/v2"
	"io/ioutil"
	"log"
	"mmm_server/config"
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type SpotifyService struct {
	repo repository.User
	cfg  *config.Configurations
}

func NewSpotifyService(repo repository.User, cfg *config.Configurations) *SpotifyService {
	return &SpotifyService{
		repo: repo,
		cfg:  cfg,
	}
}

type SpotifyAccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type SpotifyUserInfo struct {
	Country     string `json:"country"`
	DisplayName string `json:"display_name"`
	Email       string `json:"email"`
}

type SpotifyTrack struct {
	Track struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
	} `json:"track"`
}

type ResultSpSearch struct {
	Tracks struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	} `json:"tracks"`
}

func (ss *SpotifyService) GetSpotifyAccessToken(code string) string {
	// Create url search
	urlD := url.Values{}
	urlD.Add("grant_type", "authorization_code")
	urlD.Add("code", code)
	urlD.Add("redirect_uri", ss.cfg.SpotifyRedirectUrl)

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(urlD.Encode()))
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth(ss.cfg.SpotifyClientKey, ss.cfg.SpotifySecretKey)

	respAccess, err := client.Do(req)

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
	var spAccess SpotifyAccessToken
	err = json.Unmarshal(body, &spAccess)

	return spAccess.AccessToken
}

func (ss *SpotifyService) CheckSpotifyAccessToken(guestID string) bool {
	user, _ := ss.repo.GetGuestUserDB(guestID)

	if user.AccessTokenFind != "" {
		return true
	}

	return false
}

func (ss *SpotifyService) GetSpotifyUserMusic(guestID string) []model.GeneralMusicStruct {
	accessToken, _ := ss.repo.GetGuestUserDB(guestID)

	var tracks []SpotifyTrack
	u := "https://api.spotify.com/v1/me/tracks"

	for {
		var result struct {
			Items []SpotifyTrack `json:"items"`
			//Total   int           `json:"total"`
			NextURL *string `json:"next,omitempty"`
		}

		err := getSPUrl(u, &result, accessToken.AccessTokenFind)
		if err != nil {
			return nil
		}

		tracks = append(tracks, result.Items...)
		if result.NextURL == nil {
			break
		}

		u = *result.NextURL
	}

	var generalMS []model.GeneralMusicStruct
	for _, track := range tracks {
		generalMS = append(generalMS, model.GeneralMusicStruct{ID: track.Track.ID, ArtistName: track.Track.Artists[0].Name, SongName: track.Track.Name, AlbumName: track.Track.Album.Name})
	}

	return generalMS
}

func (ss *SpotifyService) MoveToSpotify(accessToken string, tracks []model.GeneralMusicStruct, con *websocket.Conn, mt int) {
	var found []string
	var notFound []string
	var moveArr [][]string

	// Search tracks
	for _, track := range tracks {
		searchString := fmt.Sprintf("%s %s", track.ArtistName, track.SongName)

		decodeSearchS := url.PathEscape(searchString)
		decodeAlbumS := url.PathEscape(track.AlbumName)

		searchUrl := "https://api.spotify.com/v1/search?q=" + decodeSearchS + "%20album:" + decodeAlbumS + "&type=track&limit=1"

		var result ResultSpSearch
		getSPUrl(searchUrl, &result, accessToken)

		if len(result.Tracks.Items) == 0 {
			searchUrl = "https://api.spotify.com/v1/search?q=" + decodeSearchS + "&type=track&limit=1"
			getSPUrl(searchUrl, &result, accessToken)

			if len(result.Tracks.Items) == 0 {
				notFound = append(notFound, searchString)
			} else {
				found = append(found, result.Tracks.Items[0].ID)
			}
		} else {
			found = append(found, result.Tracks.Items[0].ID)
		}
	}

	if len(found) > 0 {
		client := &http.Client{}

		countMusic, _ := json.Marshal(map[string]int{"lenTracks": len(found)})
		err := con.WriteMessage(mt, countMusic)
		if err != nil {
			return
		}

		// Make chunk array
		for i := 0; i < len(found); i += 50 {
			end := i + 50

			if end > len(found) {
				end = len(found)
			}

			moveArr = append(moveArr, found[i:end])
		}

		// Move tracks
		for _, ids := range moveArr {
			c := 0
			c += len(ids)

			countMusic, _ := json.Marshal(map[string]int{"countM": c})

			time.Sleep(2 * time.Second)

			err := con.WriteMessage(mt, countMusic)
			if err != nil {
				return
			}

			req, err := http.NewRequest("PUT", "https://api.spotify.com/v1/me/tracks?ids="+string(strings.Join(ids, ",")), nil)
			if err != nil {
				log.Fatalf("ERROR %v", err)
			}

			req.Header.Add("Content-Type", "application/json")
			req.Header.Add("Accept", "application/json")
			req.Header.Add("Authorization", "Bearer "+accessToken)

			resp, err := client.Do(req)
			if err != nil {
				log.Fatal(err)
			}

			resp.Body.Close()
		}
	}

	notFoundMusic, _ := json.Marshal(map[string][]string{"notFoundTracks": notFound})
	err := con.WriteMessage(mt, notFoundMusic)
	if err != nil {
		return
	}

	con.Close()
}

// Functions - Helpers
func getSPUrl(url string, result interface{}, token string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
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
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}
