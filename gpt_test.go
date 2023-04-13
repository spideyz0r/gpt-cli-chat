package main

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func TestGetUserInput(t *testing.T) {
	testCases := []struct {
		name          string
		mockInput     string
		expectedInput string
		delimiter     string
	}{
		{
			name:          "Test with newline delimiter",
			mockInput:     "Hello World\n",
			expectedInput: "Hello World",
			delimiter:     "\n",
		},
		{
			name:          "Test with , delimiter",
			mockInput:     "Hello World,",
			expectedInput: "Hello World",
			delimiter:     ",",
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockStdin := strings.NewReader(tc.mockInput)
			oldStdin := os.Stdin

			f, err := ioutil.TempFile("", "")
			if err != nil {
				t.Fatal(err)
			}
			defer os.Remove(f.Name())
			if _, err := io.Copy(f, mockStdin); err != nil {
				t.Fatal(err)
			}
			if _, err := f.Seek(0, 0); err != nil {
				t.Fatal(err)
			}

			os.Stdin = f
			actualInput := getUserInput(tc.delimiter)
			os.Stdin = oldStdin

			if actualInput != tc.expectedInput {
				t.Errorf("getUserInput() with newline delimiter returned unexpected result: expected %q, got %q", tc.expectedInput, actualInput)
			}
		})
	}
}
