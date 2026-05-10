package env

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue string
		expected     string
	}{
		{"when the key exists", true, "VALUE", "", "VALUE"},
		{"when the key does not exist", false, "", "DEFAULT", "DEFAULT"},
		{"when the key exists but is empty", true, "", "", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := Get(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetUint(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint
		expected     uint
	}{
		{"when the key exists", true, "10", 0, 10},
		{"when the key does not exist", false, "", 0, 0},
		{"when the value is not a number", true, "VALUE", 0, 0},
		{"when the value is negative", true, "-1", 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetUint(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetUint8(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint8
		expected     uint8
	}{
		{"when the key exists", true, "10", 0, 10},
		{"when the key does not exist", false, "", 0, 0},
		{"when the value is not a number", true, "VALUE", 0, 0},
		{"when the value overflows uint8", true, "300", 7, 7},
		{"when the value is negative", true, "-1", 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetUint8(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetBool(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue bool
		expected     bool
	}{
		{"when the key exists", true, "true", false, true},
		{"when the key does not exist", false, "", false, false},
		{"when the value is not a boolean", true, "VALUE", false, false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetBool(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetBytes(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue []byte
		expected     []byte
	}{
		{"when the key exists", true, "VALUE", nil, []byte("VALUE")},
		{"when the key does not exist", false, "", []byte("DEFAULT"), []byte("DEFAULT")},
		{"when the key exists but is empty", true, "", nil, []byte("")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetBytes(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetStringSlice(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue []string
		expected     []string
	}{
		{"when the key exists", true, "value1,value2,value3", nil, []string{"value1", "value2", "value3"}},
		{"when the key exists but is empty", true, "", nil, nil},
		{"when the key does not exist", false, "", []string{"DEFAULT"}, []string{"DEFAULT"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetStringSlice(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetStringSliceWithSep(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		separator    string
		defaultValue []string
		expected     []string
	}{
		{"when the key exists with comma separator", true, "a,b,c", ",", nil, []string{"a", "b", "c"}},
		{"when the key exists with semicolon separator", true, "a;b;c", ";", nil, []string{"a", "b", "c"}},
		{"when the items have surrounding whitespace", true, " a , b ,c ", ",", nil, []string{"a", "b", "c"}},
		{"when the key has a single item", true, "only", ",", nil, []string{"only"}},
		{"when the key does not exist", false, "", ",", []string{"DEFAULT"}, []string{"DEFAULT"}},
		{"when the key exists but is empty", true, "", ",", []string{"DEFAULT"}, []string{"DEFAULT"}},
		{"when the separator is multi-char", true, "a||b||c", "||", nil, []string{"a", "b", "c"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetStringSliceWithSep(key, tc.defaultValue, tc.separator)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetFloat64(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue float64
		expected     float64
	}{
		{"when the key exists", true, "10.5", 0, 10.5},
		{"when the key does not exist", false, "", 0, 0},
		{"when the value is not a number", true, "VALUE", 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetFloat64(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetInt64(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue int64
		expected     int64
	}{
		{"when the key exists", true, "10", 0, 10},
		{"when the key exists with negative", true, "-10", 0, -10},
		{"when the key does not exist", false, "", 0, 0},
		{"when the value is not a number", true, "VALUE", 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetInt64(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetInt(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue int
		expected     int
	}{
		{"when the key exists", true, "10", 0, 10},
		{"when the key exists with negative", true, "-10", 0, -10},
		{"when the key does not exist", false, "", 0, 0},
		{"when the value is not a number", true, "VALUE", 0, 0},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetInt(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetInt32(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue int32
		expected     int32
	}{
		{"when the key exists", true, "10", 0, 10},
		{"when the key exists with negative", true, "-10", 0, -10},
		{"when the key does not exist", false, "", 7, 7},
		{"when the value is not a number", true, "VALUE", 7, 7},
		{"when the value overflows int32", true, "9999999999", 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetInt32(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetUint16(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint16
		expected     uint16
	}{
		{"when the key exists", true, "1024", 0, 1024},
		{"when the key does not exist", false, "", 7, 7},
		{"when the value is not a number", true, "VALUE", 7, 7},
		{"when the value is negative", true, "-1", 7, 7},
		{"when the value overflows uint16", true, "70000", 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetUint16(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetUint32(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint32
		expected     uint32
	}{
		{"when the key exists", true, "65536", 0, 65536},
		{"when the key does not exist", false, "", 7, 7},
		{"when the value is not a number", true, "VALUE", 7, 7},
		{"when the value is negative", true, "-1", 7, 7},
		{"when the value overflows uint32", true, "9999999999", 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetUint32(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetUint64(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint64
		expected     uint64
	}{
		{"when the key exists", true, "9999999999", 0, 9999999999},
		{"when the key does not exist", false, "", 7, 7},
		{"when the value is not a number", true, "VALUE", 7, 7},
		{"when the value is negative", true, "-1", 7, 7},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetUint64(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetFloat32(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue float32
		expected     float32
	}{
		{"when the key exists", true, "3.14", 0, 3.14},
		{"when the key does not exist", false, "", 1.5, 1.5},
		{"when the value is not a number", true, "VALUE", 1.5, 1.5},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetFloat32(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetDuration(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue time.Duration
		expected     time.Duration
	}{
		{"when the key exists with seconds", true, "5s", 0, 5 * time.Second},
		{"when the key exists with mixed units", true, "1h30m", 0, 90 * time.Minute},
		{"when the key does not exist", false, "", time.Minute, time.Minute},
		{"when the value is not a duration", true, "VALUE", time.Minute, time.Minute},
		{"when the value lacks units", true, "10", time.Minute, time.Minute},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result := GetDuration(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
		})
	}
}
