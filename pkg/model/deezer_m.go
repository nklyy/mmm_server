package model

type UserInfo struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type AccessToken struct {
	AccessToken string `json:"access_token"`
	Expires     int
}

type Track struct {
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
