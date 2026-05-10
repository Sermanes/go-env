// Package env provides typed accessors for environment variables with default
// values and silent fallback when parsing fails.
package env

import (
	"os"
	"strconv"
	"strings"
)

const (
	defaultSliceSeparator string = ","
	bitSizePlatform       int    = 0
	bitSize8              int    = 8
	bitSize64             int    = 64
	decimalBase           int    = 10
)

// Get returns the value of the environment variable key, or defaultValue when
// the variable is not set.
func Get(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}

// GetUint parses key as a uint. Returns defaultValue when the variable is not
// set, cannot be parsed, or is negative.
func GetUint(key string, defaultValue uint) uint {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	number, err := strconv.ParseUint(value, decimalBase, bitSizePlatform)
	if err != nil {
		return defaultValue
	}

	return uint(number)
}

// GetUint8 parses key as a uint8. Returns defaultValue when the variable is
// not set, cannot be parsed, is negative, or overflows uint8 (>255).
func GetUint8(key string, defaultValue uint8) uint8 {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	number, err := strconv.ParseUint(value, decimalBase, bitSize8)
	if err != nil {
		return defaultValue
	}

	return uint8(number)
}

// GetBool parses key as a bool using strconv.ParseBool. Returns defaultValue
// when the variable is not set or cannot be parsed.
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

// GetBytes returns the value of key as a byte slice, or defaultValue when the
// variable is not set.
func GetBytes(key string, defaultValue []byte) []byte {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	return []byte(value)
}

// GetStringSlice splits the value of key by commas. Returns defaultValue when
// the variable is unset or empty.
func GetStringSlice(key string, defaultValue []string) []string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}

	return strings.Split(value, defaultSliceSeparator)
}

// GetFloat64 parses key as a float64. Returns defaultValue when the variable
// is not set or cannot be parsed.
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

// GetInt64 parses key as a base-10 int64. Returns defaultValue when the
// variable is not set or cannot be parsed.
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

// GetInt parses key as an int. Returns defaultValue when the variable is not
// set or cannot be parsed.
func GetInt(key string, defaultValue int) int {
	value, exist := os.LookupEnv(key)
	if !exist {
		return defaultValue
	}

	number, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return number
}
