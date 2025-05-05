package config

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
)

func Load(file ...string) error {
	return godotenv.Load(file...)
}

func GetIntEnv(key string, fallback int) int {
	v, isSet := os.LookupEnv(key)
	if !isSet {
		return fallback
	}
	log.Printf("Env var %s: %v\n", key, v)
	i, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return i
}

func GetStrEnv(key string, fallback string) string {
	v, isSet := os.LookupEnv(key)
	log.Printf("Env var %s: %v\n", key, v)
	if !isSet {
		return fallback
	}
	return v
}

func GetBoolEnv(key string, fallback bool) bool {
	v, isSet := os.LookupEnv(key)
	if !isSet {
		return fallback
	}
	b, err := strconv.ParseBool(v)
	if err != nil {
		return fallback
	}
	return b
}
