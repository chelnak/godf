package internal

import (
	"github.com/spf13/viper"
)

func GetString(key string) string {

	value := viper.GetString(key)

	if value == "" {
		panic("Configuration value for %s is missing")
	}

	return value
}

func GetInt64(key string) int64 {

	value := viper.GetInt64(key)

	if value == 0 {
		panic("Configuration value for %s is missing")
	}

	return value
}
