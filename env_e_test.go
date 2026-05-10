package env

import (
	"errors"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGetE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue string
		expected     string
	}{
		{"when the key exists", true, "VALUE", "DEFAULT", "VALUE"},
		{"when the key does not exist", false, "", "DEFAULT", "DEFAULT"},
		{"when the key exists but is empty", true, "", "DEFAULT", ""},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetE(key, tc.defaultValue)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetBoolE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue bool
		expected     bool
		expectErr    bool
	}{
		{"when the key exists", true, "true", false, true, false},
		{"when the key does not exist", false, "", true, true, false},
		{"when the value is invalid", true, "VALUE", true, true, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetBoolE(key, tc.defaultValue)

			assert.Equal(t, tc.expected, result)
			if tc.expectErr {
				var perr *ParseError
				assert.ErrorAs(t, err, &perr)
				assert.Equal(t, key, perr.Key)
				assert.Equal(t, tc.value, perr.Value)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestGetBytesE(t *testing.T) {
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
		{"when the key exists but is empty", true, "", []byte("DEFAULT"), []byte("")},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetBytesE(key, tc.defaultValue)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetStringSliceE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue []string
		expected     []string
	}{
		{"when the key exists", true, "a,b,c", nil, []string{"a", "b", "c"}},
		{"when the key does not exist", false, "", []string{"DEFAULT"}, []string{"DEFAULT"}},
		{"when the key exists but is empty", true, "", []string{"DEFAULT"}, []string{"DEFAULT"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetStringSliceE(key, tc.defaultValue)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetStringSliceWithSepE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		separator    string
		defaultValue []string
		expected     []string
	}{
		{"when the key exists with semicolon", true, "a;b;c", ";", nil, []string{"a", "b", "c"}},
		{"when the items have whitespace", true, " a , b ,c ", ",", nil, []string{"a", "b", "c"}},
		{"when the key does not exist", false, "", ",", []string{"DEFAULT"}, []string{"DEFAULT"}},
		{"when the key exists but is empty", true, "", ",", []string{"DEFAULT"}, []string{"DEFAULT"}},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetStringSliceWithSepE(key, tc.defaultValue, tc.separator)

			assert.NoError(t, err)
			assert.Equal(t, tc.expected, result)
		})
	}
}

func TestGetIntE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue int
		expected     int
		expectErr    bool
	}{
		{"when the key exists", true, "10", 0, 10, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value is not a number", true, "VALUE", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetIntE(key, tc.defaultValue)

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

func TestGetInt32E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue int32
		expected     int32
		expectErr    bool
	}{
		{"when the key exists", true, "10", 0, 10, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value is not a number", true, "VALUE", 7, 7, true},
		{"when the value overflows int32", true, "9999999999", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetInt32E(key, tc.defaultValue)

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

func TestGetInt64E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue int64
		expected     int64
		expectErr    bool
	}{
		{"when the key exists", true, "10", 0, 10, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value is not a number", true, "VALUE", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetInt64E(key, tc.defaultValue)

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

func TestGetUintE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint
		expected     uint
		expectErr    bool
	}{
		{"when the key exists", true, "10", 0, 10, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value is not a number", true, "VALUE", 7, 7, true},
		{"when the value is negative", true, "-1", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetUintE(key, tc.defaultValue)

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

func TestGetUint8E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint8
		expected     uint8
		expectErr    bool
	}{
		{"when the key exists", true, "10", 0, 10, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value overflows uint8", true, "300", 7, 7, true},
		{"when the value is negative", true, "-1", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetUint8E(key, tc.defaultValue)

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

func TestGetUint16E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint16
		expected     uint16
		expectErr    bool
	}{
		{"when the key exists", true, "1024", 0, 1024, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value overflows uint16", true, "70000", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetUint16E(key, tc.defaultValue)

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

func TestGetUint32E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint32
		expected     uint32
		expectErr    bool
	}{
		{"when the key exists", true, "65536", 0, 65536, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value overflows uint32", true, "9999999999", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetUint32E(key, tc.defaultValue)

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

func TestGetUint64E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue uint64
		expected     uint64
		expectErr    bool
	}{
		{"when the key exists", true, "9999999999", 0, 9999999999, false},
		{"when the key does not exist", false, "", 7, 7, false},
		{"when the value is not a number", true, "VALUE", 7, 7, true},
		{"when the value is negative", true, "-1", 7, 7, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetUint64E(key, tc.defaultValue)

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

func TestGetFloat32E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue float32
		expected     float32
		expectErr    bool
	}{
		{"when the key exists", true, "3.14", 0, 3.14, false},
		{"when the key does not exist", false, "", 1.5, 1.5, false},
		{"when the value is not a number", true, "VALUE", 1.5, 1.5, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetFloat32E(key, tc.defaultValue)

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

func TestGetFloat64E(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue float64
		expected     float64
		expectErr    bool
	}{
		{"when the key exists", true, "10.5", 0, 10.5, false},
		{"when the key does not exist", false, "", 1.5, 1.5, false},
		{"when the value is not a number", true, "VALUE", 1.5, 1.5, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetFloat64E(key, tc.defaultValue)

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

func TestGetDurationE(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name         string
		setEnv       bool
		value        string
		defaultValue time.Duration
		expected     time.Duration
		expectErr    bool
	}{
		{"when the key exists", true, "5s", 0, 5 * time.Second, false},
		{"when the key does not exist", false, "", time.Minute, time.Minute, false},
		{"when the value is not a duration", true, "VALUE", time.Minute, time.Minute, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			result, err := GetDurationE(key, tc.defaultValue)

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

func TestParseErrorFields(t *testing.T) {
	const key, value = "PORT", "abc"

	t.Setenv(key, value)
	_, err := GetIntE(key, 0)

	var perr *ParseError
	assert.ErrorAs(t, err, &perr)
	assert.Equal(t, key, perr.Key)
	assert.Equal(t, value, perr.Value)
	assert.NotNil(t, perr.Unwrap())
	assert.Contains(t, err.Error(), "PORT")
	assert.Contains(t, err.Error(), "abc")
}

func TestParseErrorUnwrap(t *testing.T) {
	t.Setenv("PORT", "abc")
	_, err := GetIntE("PORT", 0)

	assert.True(t, errors.Is(err, strconv.ErrSyntax))
}
