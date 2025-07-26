package env

import (
	"fmt"
	"os"
	"strconv"
)

func GetEnvOrFail(key string) string {
	value := GetEnvAllowEmptyOrFail(key)
	if value == "" {
		panic(fmt.Errorf("%s environment variable is empty", key))
	}
	return value
}

func GetEnvAllowEmptyOrFail(key string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	panic(fmt.Errorf("%s environment variable is not set", key))
}

func GetEnvIntOrFail(key string) int {
	port, err := strconv.Atoi(GetEnvOrFail(key))
	if err != nil {
		panic(fmt.Errorf("%s environment variable must be a numeric: %v", key, err))
	}
	return port
}
