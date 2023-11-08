package configs

import "os"

type Configs struct {
	BaseUrl string
}

func NewConfig() *Configs {
	var conf Configs

	conf.BaseUrl = getString("BaseUrl", "https://google.com")

	return &conf
}

func getString(key, defaultStr string) string {
	env := os.Getenv(key)
	if env == "" {
		return defaultStr
	}
	return env
}
