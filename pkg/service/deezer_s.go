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
)

type DeezerService struct {
	repo repository.User
}

func NewDeezerService(repo repository.User) *DeezerService {
	return &DeezerService{
		repo: repo,
	}
}

func (ds *DeezerService) GetDeezerAccessToken(code string) string {
	accessT := getAccessToken(code)

	return accessT.AccessToken
}

func (ds *DeezerService) CheckAccessToken(token string) bool {
	user := getUserInfo(token)

	if user.Name != "" {
		return true
	}

	return false
}

func (ds *DeezerService) GetDeezerUserMusic(token string) []model.Track {
	userInfo := getUserInfo(token)
	userTracks := getUserTracks(token)

	fmt.Println(userInfo)

	return userTracks
}

// Get Access Token
func getAccessToken(code string) model.AccessToken {
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
	var result model.AccessToken
	err = json.Unmarshal(body, &result)

	return result
}

// Get Deezer User
func getUserInfo(accessT string) *model.UserInfo {
	url := "https://api.deezer.com/user/me?access_token=" + accessT

	// Unmarshal user info
	var result model.UserInfo
	err := getUrl(url, &result)
	if err != nil {
		return nil
	}

	return &result
}

// Get User Music
func getUserTracks(accessT string) []model.Track {
	var tracks []model.Track
	url := "https://api.deezer.com/user/me/tracks?access_token=" + accessT

	for {
		var result struct {
			Data    []model.Track `json:"data"`
			Total   int           `json:"total"`
			NextURL *string       `json:"next,omitempty"`
		}

		err := getUrl(url, &result)
		if err != nil {
			return nil
		}

		tracks = append(tracks, result.Data...)
		if result.NextURL == nil {
			break
		}

		url = *result.NextURL
	}

	return tracks
}

func getUrl(url string, result interface{}) error {
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
