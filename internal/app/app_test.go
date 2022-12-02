package app

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestFuncName(t *testing.T) {
	tests := []struct {
		name             string
		request          string
		expectedResponse string
	}{
		{
			name:             "Test name here",
			request:          "request",
			expectedResponse: "response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Println("testing logic")
			require.Equal(t, tt.expectedResponse, "response")
		})
	}
}
