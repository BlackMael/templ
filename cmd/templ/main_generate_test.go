package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name           string
		args           []string
		expectedStdout string
		expectedStderr string
		expectedCode   int
	}{
		{
			name: `"templ generate {path to templ file}" to debug new conditional attributes feature`,
			args: []string{
				"templ",
				"generate",
				"-f",
				"C:\\projects\\github\\templ\\generator\\test-attribute-ifelseif\\template.templ",
			},
			expectedStdout: "Don't care at the moment",
			expectedCode:   0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			stdin := strings.NewReader("")
			stdout := bytes.NewBuffer(nil)
			stderr := bytes.NewBuffer(nil)
			actualCode := run(stdin, stdout, stderr, test.args)

			if actualCode != test.expectedCode {
				t.Errorf("expected code %v, got %v", test.expectedCode, actualCode)
			}
			if diff := cmp.Diff(test.expectedStdout, stdout.String()); diff != "" {
				t.Error(diff)
				t.Error("expected stdout:")
				t.Error(test.expectedStdout)
				t.Error("actual stdout:")
				t.Error(stdout.String())
			}
			if diff := cmp.Diff(test.expectedStderr, stderr.String()); diff != "" {
				t.Error(diff)
				t.Error("expected stderr:")
				t.Error(test.expectedStderr)
				t.Error("actual stderr:")
				t.Error(stderr.String())
			}
		})
	}
}
