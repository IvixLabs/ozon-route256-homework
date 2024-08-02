package env

import (
	"log"
	"os"
	"strconv"
	"time"
)

const (
	EnvAppGracefulShutdownTimeout = "APP_GRACEFUL_SHUTDOWN_TIMEOUT"
)

func GetEnvVar(name string) string {
	envVar := os.Getenv(name)
	if envVar == "" {
		panic(name + " is required")
	}

	return envVar
}

var globalGracefulShutdownTimeout *time.Duration

func GetGracefulShutdownTimeout() time.Duration {
	if globalGracefulShutdownTimeout == nil {
		intGracefulShutdownTimeout, err := strconv.Atoi(GetEnvVar(EnvAppGracefulShutdownTimeout))
		if err != nil {
			log.Panic(err.Error())
		}
		if intGracefulShutdownTimeout < 2 {
			log.Fatalln("Graceful shutdown timeout min is 2 seconds")
		}

		gracefulShutdownTimeout := time.Duration(intGracefulShutdownTimeout) * time.Second
		globalGracefulShutdownTimeout = &gracefulShutdownTimeout
	}

	return *globalGracefulShutdownTimeout
}
