package config

type Configurations struct {
	PORT string `mapstructure:"PORT"`

	FrontEndUrl string `mapstructure:"FRONT_END_URL"`

	SpotifyRedirectUrl string `mapstructure:"SPOTIFY_R_U"`
	SpotifyClientKey   string `mapstructure:"SPOTIFY_C_K"`
	SpotifySecretKey   string `mapstructure:"SPOTIFY_S_K"`
	SpotifyScope       string `mapstructure:"SPOTIFY_SCOPE"`

	DeezerRedirectUrl string `mapstructure:"DEEZER_R_U"`
	DeezerClientKey   string `mapstructure:"DEEZER_C_K"`
	DeezerSecretKey   string `mapstructure:"DEEZER_S_K"`
	DeezerScope       string `mapstructure:"DEEZER_SCOPE"`
}
