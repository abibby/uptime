package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

func env(key string, def string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		return def
	}
	return v
}
func mustEnv(key string) string {
	v, ok := os.LookupEnv(key)
	if !ok {
		panic(fmt.Sprintf("%s must be set in the environment", key))
	}
	return v
}

func envBool(key string, def bool) bool {
	strDef := "false"
	if def {
		strDef = "true"
	}
	str := strings.ToLower(env(key, strDef))
	return str != "false" && str != "0"
}

func envDuration(key string, def time.Duration) time.Duration {
	strValue := env(key, def.String())
	value, err := time.ParseDuration(strValue)
	if err != nil {
		log.Printf("failed to parse duration '%s' for env %s", strValue, key)
		return def
	}
	return value
}

var LogPath string
var AggregatePath string
var CheckFrequency time.Duration
var Verbose bool

func Init() error {
	err := godotenv.Load("./.env")
	if errors.Is(err, fs.ErrNotExist) {
	} else if err != nil {
		return err
	}

	LogPath = env("UPTIME_LOG_PATH", "./uptime-log/uptime")
	AggregatePath = env("UPTIME_AGGREGATE_PATH", "./uptime-log/aggregate")

	CheckFrequency = envDuration("UPTIME_CHECK_FREQUENCY", time.Second)

	Verbose = envBool("VERBOSE", false)

	return nil
}
