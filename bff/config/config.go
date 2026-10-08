package config

type Config struct {
	Kind     string `json:"kind"`
	ApiKey   string `json:"api_key"`
	Provider string `json:"provider"`
	UserName string `json:"username"`
	Password string `json:"password"`
}
