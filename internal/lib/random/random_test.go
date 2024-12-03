package random

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestRandomString(t *testing.T) {
	testCases := []struct {
		name   string
		length int
	}{
		{
			name:   "len-1",
			length: 1,
		},
		{
			name:   "len-5",
			length: 5,
		},
		{
			name:   "len-7",
			length: 7,
		},
		{
			name:   "len-15",
			length: 15,
		},
		{
			name:   "len-23",
			length: 23,
		},
		{
			name:   "len-42",
			length: 42,
		},
	}

	for _, tc := range testCases{
		str1 := RandomString(tc.length)
		time.Sleep(time.Microsecond * 1)
		str2 := RandomString(tc.length)

		assert.Len(t, str1, tc.length)
		assert.Len(t, str2, tc.length)

		assert.NotEqual(t, str1, str2)
		
	}
}