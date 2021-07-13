package service

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
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
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func (ss *SpotifyService) GetSpotifyAccessToken(code string) string {
	accessT := getSPAccessToken(code)

	return accessT.AccessToken
}

func (ss *SpotifyService) CheckSpotifyAccessToken(token string) bool {
	user := getSPUserInfo(token)

	fmt.Println(user)

	return true
}

func getSPAccessToken(code string) SpotifyAccessToken {
	urlD := url.Values{}
	urlD.Add("grant_type", "authorization_code")
	urlD.Add("code", code)
	urlD.Add("redirect_uri", "http://localhost:4000/v1/spotify/callback")
	urlD.Add("client_id", "a45422e6fcc04cc6932840b3372581f5")
	urlD.Add("client_secret", "c7bf659fdd5d40aca23125e19bf5d706")

	respAccess, err := http.Post("https://accounts.spotify.com/api/token", "application/x-www-form-urlencoded", strings.NewReader(urlD.Encode()))

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

func getSPUserInfo(accessT string) *SpotifyUserInfo {
	fmt.Println("1231312312")
	urlSP := "https://api.spotify.com/v1/me"

	var result SpotifyUserInfo
	err := getSPUrl(urlSP, &result, accessT)
	if err != nil {
		return nil
	}

	return &result
}

func getSPUrl(url string, result interface{}, token string) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := http.DefaultClient.Do(req)
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

	fmt.Println(string(body))

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}

	return nil
}
