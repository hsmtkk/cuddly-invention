package env

import (
	"fmt"
	"os"
)

func RequiredEnv(key string) (string, error) {
	val := os.Getenv(key)
	if val == "" {
		return "", fmt.Errorf("you must define %s environment variable", key)
	}
	return val, nil
}
