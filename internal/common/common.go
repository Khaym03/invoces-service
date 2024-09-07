package common

import "os"

func Getenv(key, defaultValue string) string {
	value, defined := os.LookupEnv(key)
	if !defined {
		return defaultValue
	}
	return value
}
