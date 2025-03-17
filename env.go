package env

import (
	"os"
	"strconv"
	"strings"
)

const (
	defaultSliceSeparator string = ","
	bitSize64             int    = 64
	decimalBase           int    = 10
)

func Get(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

func GetUint(key string, defaultValue uint) uint {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return uint(number)
}

func GetUint8(key string, defaultValue uint8) uint8 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return uint8(number)
}

func GetBool(key string, defaultValue bool) bool {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	boolean, err := strconv.ParseBool(value)
	if err != nil {
		return defaultValue
	}

	return boolean
}

func GetBytes(key string, defaultValue []byte) []byte {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return []byte(value)
}

func GetStringSlice(key string, defaultValue []string) []string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return strings.Split(value, defaultSliceSeparator)
}

func GetFloat64(key string, defaultValue float64) float64 {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	number, err := strconv.ParseFloat(value, bitSize64)
	if err != nil {
		return defaultValue
	}

	return number
}

func GetInt64(key string, defaultValue int64) int64 {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	number, err := strconv.ParseInt(value, decimalBase, bitSize64)
	if err != nil {
		return defaultValue
	}

	return number
}
