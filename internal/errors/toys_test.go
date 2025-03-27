package errors

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestToyNotFoundError_Error(t *testing.T) {
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
			expected: "toy not found. Base error: test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &ToyNotFoundError{Message: tc.message, BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestToyNotFoundError_Unwrap(t *testing.T) {
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
			err := &ToyNotFoundError{BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Unwrap())
		})
	}
}

func TestToyAlreadyExistsError_Error(t *testing.T) {
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
			expected: "toy already exists. Base error: test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &ToyAlreadyExistsError{Message: tc.message, BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestToyAlreadyExistsError_Unwrap(t *testing.T) {
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
			err := &ToyAlreadyExistsError{BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Unwrap())
		})
	}
}
