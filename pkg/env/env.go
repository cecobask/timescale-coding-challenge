package env

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

func LoadEnv() error {
	if err := godotenv.Load(".env"); err != nil {
		return err
	}
	env := os.Getenv("ENV")
	if env == "" {
		return nil
	}
	envFile := fmt.Sprintf(".env.%s", strings.ToLower(env))
	return godotenv.Overload(envFile)
}
