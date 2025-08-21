package env

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGet(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := value
		t.Setenv(key, value)

		result := Get(key, "")

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", "DEFAULT"
		expected := defaultValue

		result := Get(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is empty", func(t *testing.T) {
		key, value := "KEY", ""
		expected := value
		t.Setenv(key, value)

		result := Get(key, "")

		assert.Equal(t, expected, result)
	})
}

func TestGetUint(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "10"
		expected := uint(10)
		t.Setenv(key, value)

		result := GetUint(key, 0)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", uint(0)
		expected := defaultValue

		result := GetUint(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is not a number", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := uint(0)
		t.Setenv(key, value)

		result := GetUint(key, 0)

		assert.Equal(t, expected, result)
	})
}

func TestGetUint8(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "10"
		expected := uint8(10)
		t.Setenv(key, value)

		result := GetUint8(key, 0)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", uint8(0)
		expected := defaultValue

		result := GetUint8(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is not a number", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := uint8(0)
		t.Setenv(key, value)

		result := GetUint8(key, 0)

		assert.Equal(t, expected, result)
	})
}

func TestGetBool(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "true"
		expected := true
		t.Setenv(key, value)

		result := GetBool(key, false)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", false
		expected := defaultValue

		result := GetBool(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is not a boolean", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := false
		t.Setenv(key, value)

		result := GetBool(key, false)

		assert.Equal(t, expected, result)
	})
}

func TestGetBytes(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := []byte(value)
		t.Setenv(key, value)

		result := GetBytes(key, nil)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", []byte("DEFAULT")
		expected := defaultValue

		result := GetBytes(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is empty", func(t *testing.T) {
		key, value := "KEY", ""
		expected := []byte(value)
		t.Setenv(key, value)

		result := GetBytes(key, nil)

		assert.Equal(t, expected, result)
	})
}

func TestGetStringSlice(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "value1,value2,value3"
		expected := []string{"value1", "value2", "value3"}
		t.Setenv(key, value)

		result := GetStringSlice(key, nil)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is empty", func(t *testing.T) {
		key, value := "KEY", ""
		t.Setenv(key, value)

		result := GetStringSlice(key, nil)

		assert.Nil(t, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", []string{"DEFAULT"}
		expected := defaultValue

		result := GetStringSlice(key, defaultValue)

		assert.Equal(t, expected, result)
	})
}

func TestGetFloat64(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "10.5"
		expected := 10.5
		t.Setenv(key, value)

		result := GetFloat64(key, 0)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", 0.0
		expected := defaultValue

		result := GetFloat64(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is not a number", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := 0.0
		t.Setenv(key, value)

		result := GetFloat64(key, 0)

		assert.Equal(t, expected, result)
	})
}

func TestGetInt64(t *testing.T) {
	t.Run("When the key exists", func(t *testing.T) {
		key, value := "KEY", "10"
		expected := int64(10)
		t.Setenv(key, value)

		result := GetInt64(key, 0)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key does not exist", func(t *testing.T) {
		key, defaultValue := "KEY", int64(0)
		expected := defaultValue

		result := GetInt64(key, defaultValue)

		assert.Equal(t, expected, result)
	})

	t.Run("When the key exists but is not a number", func(t *testing.T) {
		key, value := "KEY", "VALUE"
		expected := int64(0)
		t.Setenv(key, value)

		result := GetInt64(key, 0)

		assert.Equal(t, expected, result)
	})
}
