package utils

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestToString(t *testing.T) {
	t.Run("should convert string to string", func(t *testing.T) {
		result := ToString("hello")
		assert.Equal(t, "hello", result)
	})

	t.Run("should convert int to string", func(t *testing.T) {
		result := ToString(123)
		assert.Equal(t, "123", result)
	})
}

func TestToInt(t *testing.T) {
	t.Run("should convert string to int", func(t *testing.T) {
		result := ToInt("123")
		assert.Equal(t, 123, result)
	})

	t.Run("should convert int to int", func(t *testing.T) {
		result := ToInt(456)
		assert.Equal(t, 456, result)
	})
}

func TestToIntE(t *testing.T) {
	t.Run("should convert valid string to int", func(t *testing.T) {
		result, err := ToIntE("789")
		assert.NoError(t, err)
		assert.Equal(t, 789, result)
	})

	t.Run("should return error for invalid string", func(t *testing.T) {
		_, err := ToIntE("invalid")
		assert.Error(t, err)
	})
}

func TestNow(t *testing.T) {
	t.Run("should return current time", func(t *testing.T) {
		result := Now()
		assert.NotNil(t, result)
		// Should be within 1 second of now
		assert.WithinDuration(t, time.Now(), result, time.Second)
	})
}

func TestNowRFC3339(t *testing.T) {
	t.Run("should return RFC3339 formatted string", func(t *testing.T) {
		result := NowRFC3339()
		assert.NotEmpty(t, result)
		// Should be parseable
		_, err := time.Parse(time.RFC3339, result)
		assert.NoError(t, err)
	})
}

func TestParseRFC3339(t *testing.T) {
	t.Run("should parse valid RFC3339 string", func(t *testing.T) {
		input := "2024-01-01T12:00:00Z"
		result, err := ParseRFC3339(input)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should return error for invalid string", func(t *testing.T) {
		_, err := ParseRFC3339("invalid")
		assert.Error(t, err)
	})
}

func TestFormatRFC3339(t *testing.T) {
	t.Run("should format time to RFC3339", func(t *testing.T) {
		input := time.Date(2024, 1, 1, 12, 0, 0, 0, time.UTC)
		result := FormatRFC3339(input)
		assert.Equal(t, "2024-01-01T12:00:00Z", result)
	})
}

func TestGenerateID(t *testing.T) {
	t.Run("should generate valid ULID", func(t *testing.T) {
		id := GenerateID()
		assert.NotEmpty(t, id)
		assert.Len(t, id, 26) // ULID is 26 characters
	})

	t.Run("should generate unique IDs", func(t *testing.T) {
		id1 := GenerateID()
		id2 := GenerateID()
		assert.NotEqual(t, id1, id2)
	})
}

func TestParseULID(t *testing.T) {
	t.Run("should parse valid ULID", func(t *testing.T) {
		id := GenerateID()
		result, err := ParseULID(id)
		assert.NoError(t, err)
		assert.NotNil(t, result)
	})

	t.Run("should return error for invalid ULID", func(t *testing.T) {
		_, err := ParseULID("invalid")
		assert.Error(t, err)
	})
}

func TestGetTimestampFromID(t *testing.T) {
	t.Run("should extract timestamp from ULID", func(t *testing.T) {
		id := GenerateID()
		timestamp, err := GetTimestampFromID(id)
		assert.NoError(t, err)
		assert.NotNil(t, timestamp)
		// Should be within 1 second of now
		assert.WithinDuration(t, time.Now(), timestamp, time.Second)
	})

	t.Run("should return error for invalid ULID", func(t *testing.T) {
		_, err := GetTimestampFromID("invalid")
		assert.Error(t, err)
	})
}
