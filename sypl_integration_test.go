// Copyright 2021 The sypl Authors. All rights reserved.
// Use of this source code is governed by a MIT
// license that can be found in the LICENSE file.

package sypl

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/saucelabs/sypl/level"
)

// List files test files.
func listFiles(t *testing.T, filePath string) []string {
	t.Helper()

	files, err := filepath.Glob(strings.TrimSuffix(filePath, filepath.Ext(filePath)) + "*")
	if err != nil {
		t.Error("Failed to list test files", err)
	}

	return files
}

// Delete test files.
func deleteTestFiles(t *testing.T, filePath string) {
	t.Helper()

	for _, f := range listFiles(t, filePath) {
		if err := os.Remove(f); err != nil {
			t.Error("Failed to delete test files", err)
		}
	}
}

func TestNewIntegration(t *testing.T) {
	if !strings.EqualFold(os.Getenv("SYPL_TEST_MODE"), "integration") {
		t.SkipNow()
	}

	type args struct {
		component string
		content   string
		dir       string
		filename  string
		level     level.Level
		maxLevel  level.Level

		run func(a args) string
	}

	realFileArgs := args{
		component: "componentNameTest",
		content:   "contentTest",
		dir:       "/tmp",
		filename:  "spyl-integration-test.log",
		level:     level.Info,
		maxLevel:  level.Debug,
		run: func(a args) string {
			filePath := filepath.Join(a.dir, a.filename)

			deleteTestFiles(t, filePath)

			New(a.component).
				AddOutput(FileWithRotation(filePath, level.Debug, &FileRotationOptions{
					Compress:   true,
					MaxAge:     28, // Days.
					MaxBackups: 5,
					MaxBytes:   50,
				}, Prefixer("Test Prefix - "))).
				Printf(a.level, "%s", a.content).
				Printf(a.level, "%s 1", a.content).
				Printf(a.level, "%s 2", a.content)

			b, err := ioutil.ReadFile(filePath)
			if err != nil {
				t.Error("Failed to read file", err)
			}

			// Rotation based on bytes.
			testFiles := listFiles(t, filePath)
			testFilesSize := len(testFiles)
			minExpectedTestFiles := 2

			if testFilesSize < minExpectedTestFiles {
				t.Errorf(
					"Failed rotation based on size. Files: %s Received size: %d, want %d",
					testFiles,
					testFilesSize,
					minExpectedTestFiles,
				)
			}

			deleteTestFiles(t, filePath)

			return string(b)
		},
	}

	tests := []struct {
		name string
		args args
		want func(a args) string
	}{
		{
			name: "Should print - File based with size rotation",
			args: realFileArgs,
			want: func(a args) string {
				return "Test Prefix - contentTest 2"
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			message := tt.args.run(tt.args)
			want := tt.want(tt.args)

			if !strings.EqualFold(message, want) {
				t.Errorf("New() = %v, want %v", message, want)
			}
		})
	}
}
