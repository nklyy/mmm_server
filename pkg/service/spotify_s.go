package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"mmm_server/pkg/model"
	"mmm_server/pkg/repository"
	"net/http"
	"net/url"
	"strings"
)

type SpotifyService struct {
	repo repository.User
}

func NewSpotifyService(repo repository.User) *SpotifyService {
	return &SpotifyService{
		repo: repo,
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
	//AddedAt string `json:"added_at"`
	Track struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Type  string `json:"type"`
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
	} `json:"track"`
}

func (ss *SpotifyService) GetSpotifyAccessToken(code string) string {
	accessT := getSPAccessToken(code)

	return accessT.AccessToken
}

func (ss *SpotifyService) CheckSpotifyAccessToken(token string) bool {
	user := getSPUserInfo(token)

	if user.Email != "" {
		return true
	}

	return false
}

func (ss *SpotifyService) GetSpotifyUserMusic(code string) []model.GeneralMusicStruct {
	userMusic := getSPUserTracks(code)
	return userMusic
}

func (ss *SpotifyService) MoveToSpotify(tracks []model.GeneralMusicStruct, code string) {
	for _, track := range tracks {
		s := fmt.Sprintf("%s - %s", track.AlbumName, track.SongName)
		err := searchSPTrack(track.AlbumName, track.SongName, s, code)
		if err != nil {
			return
		}
		break
	}
}

func getSPAccessToken(code string) SpotifyAccessToken {
	urlD := url.Values{}
	urlD.Add("grant_type", "authorization_code")
	urlD.Add("code", code)
	urlD.Add("redirect_uri", "http://localhost:4000/v1/spotify/callback")

	client := &http.Client{}
	req, err := http.NewRequest("POST", "https://accounts.spotify.com/api/token", strings.NewReader(urlD.Encode()))
	if err != nil {
		log.Fatalf("ERROR %v", err)
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Accept", "application/json")
	req.SetBasicAuth("6b990a58d275455da234d248fda89722", "bfa229942d1a444f9ab9e91266a42d73")

	//respAccess, err := http.Post("https://accounts.spotify.com/api/token", "application/x-www-form-urlencoded", strings.NewReader(urlD.Encode()))

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
	var result SpotifyAccessToken
	err = json.Unmarshal(body, &result)

	return result
}

func getSPUserTracks(accessT string) []model.GeneralMusicStruct {
	var tracks []SpotifyTrack
	u := "https://api.spotify.com/v1/me/tracks"

	for {
		var result struct {
			Items []SpotifyTrack `json:"items"`
			//Total   int           `json:"total"`
			NextURL *string `json:"next,omitempty"`
		}

		err := getSPUrl(u, &result, accessT)
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
		generalMS = append(generalMS, model.GeneralMusicStruct{ID: track.Track.ID, AlbumName: track.Track.Album.Name, SongName: track.Track.Name})
	}

	return generalMS
}

func getSPUserInfo(accessT string) *SpotifyUserInfo {
	urlSP := "https://api.spotify.com/v1/me"

	var result SpotifyUserInfo
	err := getSPUrl(urlSP, &result, accessT)
	if err != nil {
		return nil
	}

	return &result
}

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

func searchSPTrack(title string, artist string, fullStr string, code string) error {
	//t := url.PathEscape(title)
	//a := url.PathEscape(artist)
	fS := url.PathEscape(fullStr)

	sUrl := "https://api.spotify.com/v1/search?q=" + fS + "type=track"

	aToken := getSPAccessToken(code)

	client := &http.Client{}
	req, err := http.NewRequest("GET", sUrl, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+aToken.AccessToken)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	//fmt.Println(string(body))

	//err = json.Unmarshal(body, result)
	//if err != nil {
	//	return err
	//}

	return nil
}
