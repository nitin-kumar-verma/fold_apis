package config

import (
	"os"
	"sync"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Port        string
	ES_URL      string
	ES_UserName string
	ES_Password string
	AppName     string //name of microservice
}

var GlobalConfig *AppConfig
var once sync.Once

func LoadConfig() {
	once.Do(func() {
		if err := godotenv.Load(); err != nil {
			panic("Error loading envs")
		}

		port, ok1 := os.LookupEnv("PORT")
		AppName, ok2 := os.LookupEnv("APP_NAME")

		EsUrl, ok3 := os.LookupEnv("ES_URL")
		EsUsername, ok4 := os.LookupEnv("ES_USERNAME")
		EsPassword, ok5 := os.LookupEnv("ES_PASSWORD")
		if !(ok1 && ok2 && ok3 && ok4 && ok5) {
			panic("One or more envs missing")
		}

		GlobalConfig = &AppConfig{

			Port:        port,
			ES_URL:      EsUrl,
			AppName:     AppName,
			ES_UserName: EsUsername,
			ES_Password: EsPassword,
		}
	})
}
