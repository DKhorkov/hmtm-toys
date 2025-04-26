package errors

import (
	"errors"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestMasterNotFoundError_Error(t *testing.T) {
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
			expected: "master not found. Base error: test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &MasterNotFoundError{Message: tc.message, BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestMasterNotFoundError_Unwrap(t *testing.T) {
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
			err := &MasterNotFoundError{BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Unwrap())
		})
	}
}

func TestMasterAlreadyExistsError_Error(t *testing.T) {
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
			expected: "master already exists. Base error: test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &MasterAlreadyExistsError{Message: tc.message, BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestMasterAlreadyExistsError_Unwrap(t *testing.T) {
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
			err := &MasterAlreadyExistsError{BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Unwrap())
		})
	}
}

func TestInvalidMasterInfoError_Error(t *testing.T) {
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
			expected: "invalid master info. Base error: test error",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := &InvalidMasterInfoError{Message: tc.message, BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Error())
		})
	}
}

func TestInvalidMasterInfoError_Unwrap(t *testing.T) {
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
			err := &InvalidMasterInfoError{BaseErr: tc.baseErr}
			require.Equal(t, tc.expected, err.Unwrap())
		})
	}
}
