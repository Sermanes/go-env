package env

import (
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func mustParseURL(t *testing.T, raw string) *url.URL {
	t.Helper()
	u, err := url.Parse(raw)
	assert.NoError(t, err)
	return u
}

// ─── URL ─────────────────────────────────────────────────────────────────────

func TestGetURL(t *testing.T) {
	const key = "URL"
	defaultURL := mustParseURL(t, "https://default.example.com")

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue *url.URL
		expected     *url.URL
	}{
		{"when the key exists with absolute URL", true, "https://api.example.com/v1", defaultURL, mustParseURL(t, "https://api.example.com/v1")},
		{"when the key does not exist", false, "", defaultURL, defaultURL},
		{"when the value is a relative URL", true, "/path/only", defaultURL, defaultURL},
		{"when the value is malformed", true, "://broken", defaultURL, defaultURL},
		{"when the value is empty", true, "", defaultURL, defaultURL},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetURL(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetURLE(t *testing.T) {
	const key = "URL"
	defaultURL := mustParseURL(t, "https://default.example.com")

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue *url.URL
		expected     *url.URL
		expectErr    bool
	}{
		{"when the key exists", true, "https://api.example.com", defaultURL, mustParseURL(t, "https://api.example.com"), false},
		{"when the key does not exist", false, "", defaultURL, defaultURL, false},
		{"when the value is relative", true, "/path", defaultURL, defaultURL, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetURLE(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
			if tc.expectErr {
				var perr *ParseError
				assert.ErrorAs(t, err, &perr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMustGetURL(t *testing.T) {
	const key = "URL"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    *url.URL
		expectPanic bool
	}{
		{"when the key exists", true, "https://api.example.com", mustParseURL(t, "https://api.example.com"), false},
		{"when the key does not exist", false, "", nil, true},
		{"when the value is relative", true, "/path", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetURL(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetURL(key))
			}
		})
	}
}

// ─── Time (RFC3339) ──────────────────────────────────────────────────────────

func TestGetTime(t *testing.T) {
	const key = "DEADLINE"
	defaultTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2026, 5, 10, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue time.Time
		expected     time.Time
	}{
		{"when the key exists with RFC3339", true, "2026-05-10T12:30:00Z", defaultTime, expected},
		{"when the key does not exist", false, "", defaultTime, defaultTime},
		{"when the value is not RFC3339", true, "2026-05-10", defaultTime, defaultTime},
		{"when the value is garbage", true, "VALUE", defaultTime, defaultTime},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetTime(key, tc.defaultValue)

			assert.True(t, tc.expected.Equal(result), "expected %v got %v", tc.expected, result)
		})
	}
}

func TestGetTimeE(t *testing.T) {
	const key = "DEADLINE"
	defaultTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2026, 5, 10, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue time.Time
		expected     time.Time
		expectErr    bool
	}{
		{"when the key exists", true, "2026-05-10T12:30:00Z", defaultTime, expected, false},
		{"when the key does not exist", false, "", defaultTime, defaultTime, false},
		{"when the value is invalid", true, "VALUE", defaultTime, defaultTime, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetTimeE(key, tc.defaultValue)

			assert.True(t, tc.expected.Equal(result))
			if tc.expectErr {
				var perr *ParseError
				assert.ErrorAs(t, err, &perr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMustGetTime(t *testing.T) {
	const key = "DEADLINE"
	expected := time.Date(2026, 5, 10, 12, 30, 0, 0, time.UTC)

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    time.Time
		expectPanic bool
	}{
		{"when the key exists", true, "2026-05-10T12:30:00Z", expected, false},
		{"when the key does not exist", false, "", time.Time{}, true},
		{"when the value is invalid", true, "VALUE", time.Time{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetTime(key) })
			} else {
				assert.True(t, tc.expected.Equal(MustGetTime(key)))
			}
		})
	}
}

// ─── Time with custom layout ─────────────────────────────────────────────────

func TestGetTimeWithLayout(t *testing.T) {
	const key = "DEADLINE"
	defaultTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	layoutDate := "2006-01-02"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		layout       string
		defaultValue time.Time
		expected     time.Time
	}{
		{"when the key exists with matching layout", true, "2026-05-10", layoutDate, defaultTime, expected},
		{"when the key does not exist", false, "", layoutDate, defaultTime, defaultTime},
		{"when the value does not match the layout", true, "2026-05-10T12:30:00Z", layoutDate, defaultTime, defaultTime},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetTimeWithLayout(key, tc.defaultValue, tc.layout)

			assert.True(t, tc.expected.Equal(result), "expected %v got %v", tc.expected, result)
		})
	}
}

func TestGetTimeWithLayoutE(t *testing.T) {
	const key = "DEADLINE"
	defaultTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	expected := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	layoutDate := "2006-01-02"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		layout       string
		defaultValue time.Time
		expected     time.Time
		expectErr    bool
	}{
		{"when the key exists", true, "2026-05-10", layoutDate, defaultTime, expected, false},
		{"when the key does not exist", false, "", layoutDate, defaultTime, defaultTime, false},
		{"when the value does not match the layout", true, "VALUE", layoutDate, defaultTime, defaultTime, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetTimeWithLayoutE(key, tc.defaultValue, tc.layout)

			assert.True(t, tc.expected.Equal(result))
			if tc.expectErr {
				var perr *ParseError
				assert.ErrorAs(t, err, &perr)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestMustGetTimeWithLayout(t *testing.T) {
	const key = "DEADLINE"
	expected := time.Date(2026, 5, 10, 0, 0, 0, 0, time.UTC)
	layoutDate := "2006-01-02"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		layout      string
		expected    time.Time
		expectPanic bool
	}{
		{"when the key exists", true, "2026-05-10", layoutDate, expected, false},
		{"when the key does not exist", false, "", layoutDate, time.Time{}, true},
		{"when the value does not match the layout", true, "VALUE", layoutDate, time.Time{}, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetTimeWithLayout(key, tc.layout) })
			} else {
				assert.True(t, tc.expected.Equal(MustGetTimeWithLayout(key, tc.layout)))
			}
		})
	}
}
