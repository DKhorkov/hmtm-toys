package errors

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTagNotFoundError_Error(t *testing.T) {
	testCases := []struct {
		name     string
		message  string
		baseErr  error
		expected string
	}{
		{
			name:     "with message",
			message:  "test message",
			expected: "test message",
		},
		{
			name:     "with BaseErr",
			baseErr:  errors.New("test error"),
			expected: "tag not found. Base error: test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &TagNotFoundError{Message: tc.message, BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestTagNotFoundError_Unwrap(t *testing.T) {
	testCases := []struct {
		name     string
		baseErr  error
		expected error
	}{
		{
			name:     "with BaseErr",
			baseErr:  errors.New("test error"),
			expected: errors.New("test error"),
		},
		{
			name: "without BaseErr",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &TagNotFoundError{BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Unwrap())
		})
	}
}
