package util

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

func Get(key string) string {
	viper.SetConfigFile("./.env")
	err := viper.ReadInConfig()

	if err != nil {
		return os.Getenv(key) //for production
	}

	value, ok := viper.Get(key).(string)

	if !ok {
		log.Fatalf("Invalid type assertion")
	}

	return value
}
