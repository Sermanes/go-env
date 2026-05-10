package env

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMustGet(t *testing.T) {
	const key = "KEY"

	t.Run("when the key exists", func(t *testing.T) {
		t.Setenv(key, "VALUE")
		assert.Equal(t, "VALUE", MustGet(key))
	})

	t.Run("when the key does not exist", func(t *testing.T) {
		assert.PanicsWithError(t, `env: "KEY": env: variable is not set`, func() {
			MustGet(key)
		})
	})
}

func TestMustGetBool(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    bool
		expectPanic bool
	}{
		{"when the key exists", true, "true", true, false},
		{"when the key does not exist", false, "", false, true},
		{"when the value is invalid", true, "VALUE", false, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetBool(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetBool(key))
			}
		})
	}
}

func TestMustGetBytes(t *testing.T) {
	const key = "KEY"

	t.Run("when the key exists", func(t *testing.T) {
		t.Setenv(key, "VALUE")
		assert.Equal(t, []byte("VALUE"), MustGetBytes(key))
	})

	t.Run("when the key does not exist", func(t *testing.T) {
		assert.Panics(t, func() { MustGetBytes(key) })
	})
}

func TestMustGetStringSlice(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    []string
		expectPanic bool
	}{
		{"when the key exists", true, "a,b,c", []string{"a", "b", "c"}, false},
		{"when the key does not exist", false, "", nil, true},
		{"when the key exists but is empty", true, "", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetStringSlice(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetStringSlice(key))
			}
		})
	}
}

func TestMustGetStringSliceWithSep(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		separator   string
		expected    []string
		expectPanic bool
	}{
		{"when the key exists", true, "a;b;c", ";", []string{"a", "b", "c"}, false},
		{"when the items have whitespace", true, " a , b ", ",", []string{"a", "b"}, false},
		{"when the key does not exist", false, "", ",", nil, true},
		{"when the key exists but is empty", true, "", ",", nil, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetStringSliceWithSep(key, tc.separator) })
			} else {
				assert.Equal(t, tc.expected, MustGetStringSliceWithSep(key, tc.separator))
			}
		})
	}
}

func TestMustGetInt(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    int
		expectPanic bool
	}{
		{"when the key exists", true, "10", 10, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is not a number", true, "VALUE", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetInt(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetInt(key))
			}
		})
	}
}

func TestMustGetInt32(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    int32
		expectPanic bool
	}{
		{"when the key exists", true, "10", 10, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value overflows int32", true, "9999999999", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetInt32(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetInt32(key))
			}
		})
	}
}

func TestMustGetInt64(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    int64
		expectPanic bool
	}{
		{"when the key exists", true, "10", 10, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is not a number", true, "VALUE", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetInt64(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetInt64(key))
			}
		})
	}
}

func TestMustGetUint(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    uint
		expectPanic bool
	}{
		{"when the key exists", true, "10", 10, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is negative", true, "-1", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetUint(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetUint(key))
			}
		})
	}
}

func TestMustGetUint8(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    uint8
		expectPanic bool
	}{
		{"when the key exists", true, "10", 10, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value overflows uint8", true, "300", 0, true},
		{"when the value is negative", true, "-1", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetUint8(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetUint8(key))
			}
		})
	}
}

func TestMustGetUint16(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    uint16
		expectPanic bool
	}{
		{"when the key exists", true, "1024", 1024, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value overflows uint16", true, "70000", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetUint16(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetUint16(key))
			}
		})
	}
}

func TestMustGetUint32(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    uint32
		expectPanic bool
	}{
		{"when the key exists", true, "65536", 65536, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value overflows uint32", true, "9999999999", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetUint32(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetUint32(key))
			}
		})
	}
}

func TestMustGetUint64(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    uint64
		expectPanic bool
	}{
		{"when the key exists", true, "9999999999", 9999999999, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is negative", true, "-1", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetUint64(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetUint64(key))
			}
		})
	}
}

func TestMustGetFloat32(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    float32
		expectPanic bool
	}{
		{"when the key exists", true, "3.14", 3.14, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is not a number", true, "VALUE", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetFloat32(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetFloat32(key))
			}
		})
	}
}

func TestMustGetFloat64(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    float64
		expectPanic bool
	}{
		{"when the key exists", true, "10.5", 10.5, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is not a number", true, "VALUE", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetFloat64(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetFloat64(key))
			}
		})
	}
}

func TestMustGetDuration(t *testing.T) {
	const key = "KEY"

	tests := []struct {
		name        string
		setEnv      bool
		value       string
		expected    time.Duration
		expectPanic bool
	}{
		{"when the key exists", true, "5s", 5 * time.Second, false},
		{"when the key does not exist", false, "", 0, true},
		{"when the value is not a duration", true, "VALUE", 0, true},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			if tc.setEnv {
				t.Setenv(key, tc.value)
			}

			if tc.expectPanic {
				assert.Panics(t, func() { MustGetDuration(key) })
			} else {
				assert.Equal(t, tc.expected, MustGetDuration(key))
			}
		})
	}
}

func TestMustGetMissingErrIsErrMissing(t *testing.T) {
	defer func() {
		r := recover()
		assert.NotNil(t, r)
		err, ok := r.(error)
		assert.True(t, ok)
		assert.True(t, errors.Is(err, ErrMissing))
	}()

	MustGetInt("UNSET_KEY_FOR_TEST")
}

func TestMustGetParseErrIsParseError(t *testing.T) {
	t.Setenv("KEY", "abc")

	defer func() {
		r := recover()
		assert.NotNil(t, r)
		perr, ok := r.(*ParseError)
		assert.True(t, ok)
		assert.Equal(t, "KEY", perr.Key)
		assert.Equal(t, "abc", perr.Value)
	}()

	MustGetInt("KEY")
}
