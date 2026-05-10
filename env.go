// Package env provides typed accessors for environment variables with default
// values and silent fallback when parsing fails. Variants suffixed with E
// surface parse errors instead of swallowing them, and variants prefixed with
// Must panic when the variable is missing or invalid.
package env

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	defaultSliceSeparator string = ","
	bitSizePlatform       int    = 0
	bitSize8              int    = 8
	bitSize16             int    = 16
	bitSize32             int    = 32
	bitSize64             int    = 64
	decimalBase           int    = 10
)

// ErrMissing is reported when a Must* function is called for a variable that
// is not set or whose value is empty for slice accessors.
var ErrMissing = errors.New("env: variable is not set")

// errEmptyValue marks an empty environment value as semantically "not
// provided" (used by slice parsers).
var errEmptyValue = errors.New("env: empty value")

// ParseError is returned by *E accessors when the variable is set but its
// value cannot be parsed into the requested type.
type ParseError struct {
	Key   string
	Value string
	Err   error
}

// Error implements the error interface.
func (e *ParseError) Error() string {
	return fmt.Sprintf("env: parse %q value %q: %v", e.Key, e.Value, e.Err)
}

// Unwrap returns the underlying parse error.
func (e *ParseError) Unwrap() error { return e.Err }

// parse looks up key, applies parser, and returns defaultValue when the key is
// missing, the value is empty (for parsers signaling errEmptyValue), or the
// parser returns any other error.
func parse[T any](key string, defaultValue T, parser func(string) (T, error)) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}

	result, err := parser(value)
	if err != nil {
		return defaultValue
	}

	return result
}

// parseE looks up key and applies parser. Returns (defaultValue, nil) when the
// variable is missing or its value is treated as empty. Returns
// (defaultValue, *ParseError) when parsing the value fails.
func parseE[T any](key string, defaultValue T, parser func(string) (T, error)) (T, error) {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue, nil
	}

	result, err := parser(value)
	if errors.Is(err, errEmptyValue) {
		return defaultValue, nil
	}
	if err != nil {
		return defaultValue, &ParseError{Key: key, Value: value, Err: err}
	}

	return result, nil
}

// mustParse looks up key and applies parser. Panics with ErrMissing when the
// variable is missing or its value is empty for slice parsers. Panics with
// *ParseError when parsing fails.
func mustParse[T any](key string, parser func(string) (T, error)) T {
	value, exists := os.LookupEnv(key)
	if !exists {
		panic(fmt.Errorf("env: %q: %w", key, ErrMissing))
	}

	result, err := parser(value)
	if errors.Is(err, errEmptyValue) {
		panic(fmt.Errorf("env: %q: %w", key, ErrMissing))
	}
	if err != nil {
		panic(&ParseError{Key: key, Value: value, Err: err})
	}

	return result
}

// ─── Parsers ─────────────────────────────────────────────────────────────────

func parseString(s string) (string, error) { return s, nil }

func parseBytes(s string) ([]byte, error) { return []byte(s), nil }

func parseStringSliceComma(s string) ([]string, error) {
	if s == "" {
		return nil, errEmptyValue
	}
	return strings.Split(s, defaultSliceSeparator), nil
}

func makeStringSliceParser(separator string) func(string) ([]string, error) {
	return func(s string) ([]string, error) {
		if s == "" {
			return nil, errEmptyValue
		}

		parts := strings.Split(s, separator)
		for i, p := range parts {
			parts[i] = strings.TrimSpace(p)
		}

		return parts, nil
	}
}

func parseInt32(s string) (int32, error) {
	n, err := strconv.ParseInt(s, decimalBase, bitSize32)
	return int32(n), err
}

func parseInt64(s string) (int64, error) {
	return strconv.ParseInt(s, decimalBase, bitSize64)
}

func parseUint(s string) (uint, error) {
	n, err := strconv.ParseUint(s, decimalBase, bitSizePlatform)
	return uint(n), err
}

func parseUint8(s string) (uint8, error) {
	n, err := strconv.ParseUint(s, decimalBase, bitSize8)
	return uint8(n), err
}

func parseUint16(s string) (uint16, error) {
	n, err := strconv.ParseUint(s, decimalBase, bitSize16)
	return uint16(n), err
}

func parseUint32(s string) (uint32, error) {
	n, err := strconv.ParseUint(s, decimalBase, bitSize32)
	return uint32(n), err
}

func parseUint64(s string) (uint64, error) {
	return strconv.ParseUint(s, decimalBase, bitSize64)
}

func parseFloat32(s string) (float32, error) {
	n, err := strconv.ParseFloat(s, bitSize32)
	return float32(n), err
}

func parseFloat64(s string) (float64, error) {
	return strconv.ParseFloat(s, bitSize64)
}

func parseURL(s string) (*url.URL, error) {
	u, err := url.Parse(s)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "" || u.Host == "" {
		return nil, fmt.Errorf("URL must be absolute (scheme and host required), got %q", s)
	}
	return u, nil
}

func parseTimeRFC3339(s string) (time.Time, error) {
	return time.Parse(time.RFC3339, s)
}

func makeTimeParser(layout string) func(string) (time.Time, error) {
	return func(s string) (time.Time, error) {
		return time.Parse(layout, s)
	}
}

// ─── String ──────────────────────────────────────────────────────────────────

// Get returns the value of the environment variable key, or defaultValue when
// the variable is not set.
func Get(key, defaultValue string) string {
	return parse(key, defaultValue, parseString)
}

// GetE returns the value of key. The error is always nil; the variant exists
// for API symmetry with the other accessors.
func GetE(key, defaultValue string) (string, error) {
	return parseE(key, defaultValue, parseString)
}

// MustGet returns the value of key or panics with ErrMissing when unset.
func MustGet(key string) string {
	return mustParse(key, parseString)
}

// ─── Bool ────────────────────────────────────────────────────────────────────

// GetBool parses key as a bool using strconv.ParseBool. Returns defaultValue
// when the variable is not set or cannot be parsed.
func GetBool(key string, defaultValue bool) bool {
	return parse(key, defaultValue, strconv.ParseBool)
}

// GetBoolE is the error-returning variant of GetBool.
func GetBoolE(key string, defaultValue bool) (bool, error) {
	return parseE(key, defaultValue, strconv.ParseBool)
}

// MustGetBool returns the value of key parsed as bool, panicking when the
// variable is missing or cannot be parsed.
func MustGetBool(key string) bool {
	return mustParse(key, strconv.ParseBool)
}

// ─── Bytes ───────────────────────────────────────────────────────────────────

// GetBytes returns the value of key as a byte slice, or defaultValue when the
// variable is not set.
func GetBytes(key string, defaultValue []byte) []byte {
	return parse(key, defaultValue, parseBytes)
}

// GetBytesE is the error-returning variant of GetBytes.
func GetBytesE(key string, defaultValue []byte) ([]byte, error) {
	return parseE(key, defaultValue, parseBytes)
}

// MustGetBytes returns the value of key as bytes, panicking when the variable
// is missing.
func MustGetBytes(key string) []byte {
	return mustParse(key, parseBytes)
}

// ─── String slices ───────────────────────────────────────────────────────────

// GetStringSlice splits the value of key by commas. Returns defaultValue when
// the variable is unset or empty.
func GetStringSlice(key string, defaultValue []string) []string {
	return parse(key, defaultValue, parseStringSliceComma)
}

// GetStringSliceE is the error-returning variant of GetStringSlice. An empty
// value is treated as missing and returns (defaultValue, nil).
func GetStringSliceE(key string, defaultValue []string) ([]string, error) {
	return parseE(key, defaultValue, parseStringSliceComma)
}

// MustGetStringSlice returns the value of key split by commas, panicking when
// the variable is missing or empty.
func MustGetStringSlice(key string) []string {
	return mustParse(key, parseStringSliceComma)
}

// GetStringSliceWithSep splits the value of key by separator and trims
// surrounding whitespace from each item. Returns defaultValue when the
// variable is unset or empty.
func GetStringSliceWithSep(key string, defaultValue []string, separator string) []string {
	return parse(key, defaultValue, makeStringSliceParser(separator))
}

// GetStringSliceWithSepE is the error-returning variant of
// GetStringSliceWithSep.
func GetStringSliceWithSepE(key string, defaultValue []string, separator string) ([]string, error) {
	return parseE(key, defaultValue, makeStringSliceParser(separator))
}

// MustGetStringSliceWithSep returns the value of key split by separator with
// each item trimmed, panicking when the variable is missing or empty.
func MustGetStringSliceWithSep(key, separator string) []string {
	return mustParse(key, makeStringSliceParser(separator))
}

// ─── Signed ints ─────────────────────────────────────────────────────────────

// GetInt parses key as an int. Returns defaultValue when the variable is not
// set or cannot be parsed.
func GetInt(key string, defaultValue int) int {
	return parse(key, defaultValue, strconv.Atoi)
}

// GetIntE is the error-returning variant of GetInt.
func GetIntE(key string, defaultValue int) (int, error) {
	return parseE(key, defaultValue, strconv.Atoi)
}

// MustGetInt returns the value of key parsed as int, panicking when the
// variable is missing or cannot be parsed.
func MustGetInt(key string) int {
	return mustParse(key, strconv.Atoi)
}

// GetInt32 parses key as an int32. Returns defaultValue when the variable is
// not set, cannot be parsed, or overflows int32.
func GetInt32(key string, defaultValue int32) int32 {
	return parse(key, defaultValue, parseInt32)
}

// GetInt32E is the error-returning variant of GetInt32.
func GetInt32E(key string, defaultValue int32) (int32, error) {
	return parseE(key, defaultValue, parseInt32)
}

// MustGetInt32 returns the value of key parsed as int32, panicking when the
// variable is missing, cannot be parsed, or overflows int32.
func MustGetInt32(key string) int32 {
	return mustParse(key, parseInt32)
}

// GetInt64 parses key as a base-10 int64. Returns defaultValue when the
// variable is not set or cannot be parsed.
func GetInt64(key string, defaultValue int64) int64 {
	return parse(key, defaultValue, parseInt64)
}

// GetInt64E is the error-returning variant of GetInt64.
func GetInt64E(key string, defaultValue int64) (int64, error) {
	return parseE(key, defaultValue, parseInt64)
}

// MustGetInt64 returns the value of key parsed as int64, panicking when the
// variable is missing or cannot be parsed.
func MustGetInt64(key string) int64 {
	return mustParse(key, parseInt64)
}

// ─── Unsigned ints ───────────────────────────────────────────────────────────

// GetUint parses key as a uint. Returns defaultValue when the variable is not
// set, cannot be parsed, or is negative.
func GetUint(key string, defaultValue uint) uint {
	return parse(key, defaultValue, parseUint)
}

// GetUintE is the error-returning variant of GetUint.
func GetUintE(key string, defaultValue uint) (uint, error) {
	return parseE(key, defaultValue, parseUint)
}

// MustGetUint returns the value of key parsed as uint, panicking when the
// variable is missing, cannot be parsed, or is negative.
func MustGetUint(key string) uint {
	return mustParse(key, parseUint)
}

// GetUint8 parses key as a uint8. Returns defaultValue when the variable is
// not set, cannot be parsed, is negative, or overflows uint8 (>255).
func GetUint8(key string, defaultValue uint8) uint8 {
	return parse(key, defaultValue, parseUint8)
}

// GetUint8E is the error-returning variant of GetUint8.
func GetUint8E(key string, defaultValue uint8) (uint8, error) {
	return parseE(key, defaultValue, parseUint8)
}

// MustGetUint8 returns the value of key parsed as uint8, panicking when the
// variable is missing, cannot be parsed, is negative, or overflows uint8.
func MustGetUint8(key string) uint8 {
	return mustParse(key, parseUint8)
}

// GetUint16 parses key as a uint16. Returns defaultValue when the variable is
// not set, cannot be parsed, is negative, or overflows uint16 (>65535).
func GetUint16(key string, defaultValue uint16) uint16 {
	return parse(key, defaultValue, parseUint16)
}

// GetUint16E is the error-returning variant of GetUint16.
func GetUint16E(key string, defaultValue uint16) (uint16, error) {
	return parseE(key, defaultValue, parseUint16)
}

// MustGetUint16 returns the value of key parsed as uint16, panicking when the
// variable is missing, cannot be parsed, is negative, or overflows uint16.
func MustGetUint16(key string) uint16 {
	return mustParse(key, parseUint16)
}

// GetUint32 parses key as a uint32. Returns defaultValue when the variable is
// not set, cannot be parsed, is negative, or overflows uint32.
func GetUint32(key string, defaultValue uint32) uint32 {
	return parse(key, defaultValue, parseUint32)
}

// GetUint32E is the error-returning variant of GetUint32.
func GetUint32E(key string, defaultValue uint32) (uint32, error) {
	return parseE(key, defaultValue, parseUint32)
}

// MustGetUint32 returns the value of key parsed as uint32, panicking when the
// variable is missing, cannot be parsed, is negative, or overflows uint32.
func MustGetUint32(key string) uint32 {
	return mustParse(key, parseUint32)
}

// GetUint64 parses key as a uint64. Returns defaultValue when the variable is
// not set, cannot be parsed, or is negative.
func GetUint64(key string, defaultValue uint64) uint64 {
	return parse(key, defaultValue, parseUint64)
}

// GetUint64E is the error-returning variant of GetUint64.
func GetUint64E(key string, defaultValue uint64) (uint64, error) {
	return parseE(key, defaultValue, parseUint64)
}

// MustGetUint64 returns the value of key parsed as uint64, panicking when the
// variable is missing, cannot be parsed, or is negative.
func MustGetUint64(key string) uint64 {
	return mustParse(key, parseUint64)
}

// ─── Floats ──────────────────────────────────────────────────────────────────

// GetFloat32 parses key as a float32. Returns defaultValue when the variable
// is not set or cannot be parsed.
func GetFloat32(key string, defaultValue float32) float32 {
	return parse(key, defaultValue, parseFloat32)
}

// GetFloat32E is the error-returning variant of GetFloat32.
func GetFloat32E(key string, defaultValue float32) (float32, error) {
	return parseE(key, defaultValue, parseFloat32)
}

// MustGetFloat32 returns the value of key parsed as float32, panicking when
// the variable is missing or cannot be parsed.
func MustGetFloat32(key string) float32 {
	return mustParse(key, parseFloat32)
}

// GetFloat64 parses key as a float64. Returns defaultValue when the variable
// is not set or cannot be parsed.
func GetFloat64(key string, defaultValue float64) float64 {
	return parse(key, defaultValue, parseFloat64)
}

// GetFloat64E is the error-returning variant of GetFloat64.
func GetFloat64E(key string, defaultValue float64) (float64, error) {
	return parseE(key, defaultValue, parseFloat64)
}

// MustGetFloat64 returns the value of key parsed as float64, panicking when
// the variable is missing or cannot be parsed.
func MustGetFloat64(key string) float64 {
	return mustParse(key, parseFloat64)
}

// ─── Duration ────────────────────────────────────────────────────────────────

// GetDuration parses key as a time.Duration using time.ParseDuration. Returns
// defaultValue when the variable is not set or cannot be parsed.
func GetDuration(key string, defaultValue time.Duration) time.Duration {
	return parse(key, defaultValue, time.ParseDuration)
}

// GetDurationE is the error-returning variant of GetDuration.
func GetDurationE(key string, defaultValue time.Duration) (time.Duration, error) {
	return parseE(key, defaultValue, time.ParseDuration)
}

// MustGetDuration returns the value of key parsed as time.Duration, panicking
// when the variable is missing or cannot be parsed.
func MustGetDuration(key string) time.Duration {
	return mustParse(key, time.ParseDuration)
}

// ─── URL ─────────────────────────────────────────────────────────────────────

// GetURL parses key as an absolute URL using url.ParseRequestURI. Returns
// defaultValue when the variable is not set, empty, or not an absolute URL.
func GetURL(key string, defaultValue *url.URL) *url.URL {
	return parse(key, defaultValue, parseURL)
}

// GetURLE is the error-returning variant of GetURL.
func GetURLE(key string, defaultValue *url.URL) (*url.URL, error) {
	return parseE(key, defaultValue, parseURL)
}

// MustGetURL returns the value of key parsed as an absolute URL, panicking
// when the variable is missing or invalid.
func MustGetURL(key string) *url.URL {
	return mustParse(key, parseURL)
}

// ─── Time ────────────────────────────────────────────────────────────────────

// GetTime parses key as an RFC3339 timestamp. Returns defaultValue when the
// variable is not set or cannot be parsed.
func GetTime(key string, defaultValue time.Time) time.Time {
	return parse(key, defaultValue, parseTimeRFC3339)
}

// GetTimeE is the error-returning variant of GetTime.
func GetTimeE(key string, defaultValue time.Time) (time.Time, error) {
	return parseE(key, defaultValue, parseTimeRFC3339)
}

// MustGetTime returns the value of key parsed as RFC3339, panicking when the
// variable is missing or cannot be parsed.
func MustGetTime(key string) time.Time {
	return mustParse(key, parseTimeRFC3339)
}

// GetTimeWithLayout parses key using the supplied time layout. Returns
// defaultValue when the variable is not set or cannot be parsed.
func GetTimeWithLayout(key string, defaultValue time.Time, layout string) time.Time {
	return parse(key, defaultValue, makeTimeParser(layout))
}

// GetTimeWithLayoutE is the error-returning variant of GetTimeWithLayout.
func GetTimeWithLayoutE(key string, defaultValue time.Time, layout string) (time.Time, error) {
	return parseE(key, defaultValue, makeTimeParser(layout))
}

// MustGetTimeWithLayout returns the value of key parsed with the supplied
// layout, panicking when the variable is missing or cannot be parsed.
func MustGetTimeWithLayout(key, layout string) time.Time {
	return mustParse(key, makeTimeParser(layout))
}
