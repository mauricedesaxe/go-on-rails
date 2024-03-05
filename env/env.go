package env

import (
	"log"
	"os"
	"reflect"

	"github.com/joho/godotenv"
)

// Env is a struct that holds all the environment variables
// that the application needs. It's used to initialize the
// environment variables and provide a single source of truth
// for the application's configuration.
type Env struct {
	BaseUrl         string `env:"BASE_URL;default=http://localhost:3000"`
	MjApiKeyPublic  string `env:"MJ_APIKEY_PUBLIC"`
	MjApiKeyPrivate string `env:"MJ_APIKEY_PRIVATE"`
	// * Add more environment variables here
}

// Config is a struct that holds the configuration for the
// environment initialization process.
type Config struct {
	UseDotEnv bool
}

// Try to load environment variables from a .env file (if UseDotEnv is true), panic if it fails.
// Attempt to find all environment variables based on the `env` tag and set the value of the field.
// If a value is empty and doesn't have a default tag, panic.
func (e *Env) Init(c Config) {
	if c.UseDotEnv {
		err := godotenv.Load()
		if err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
	}

	val := reflect.ValueOf(e).Elem()
	for i := 0; i < val.NumField(); i++ {
		field := val.Type().Field(i)
		envVar, ok := field.Tag.Lookup("env")
		if !ok {
			continue
		}
		envValue := os.Getenv(envVar)
		if envValue == "" {
			envValue, ok = field.Tag.Lookup("default")
			if !ok {
				log.Panicf("Environment variable %s not found", envVar)
			}
		}
		val.Field(i).SetString(envValue)
	}
}
