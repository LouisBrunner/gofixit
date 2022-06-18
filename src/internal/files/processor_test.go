package files

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/LouisBrunner/gofixit/src/contracts"
	"github.com/stretchr/testify/assert"
)

func Test_ProcessFiles_echo(t *testing.T) {
	echoProcessor := func(filepath string) (string, error) {
		return filepath, nil
	}

	pwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("could not read working directory: %v", err)
	}

	tests := []struct {
		name    string
		config  contracts.FilesProcessorConfig[string]
		inputs  []string
		want    map[string]string
		wantErr bool
	}{
		{
			name: "works with relative files",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
			},
			inputs: []string{
				"testdata/file.c",
				"testdata/sub/file2.c",
			},
			want: map[string]string{
				"testdata/file.c":      "testdata/file.c",
				"testdata/sub/file2.c": "testdata/sub/file2.c",
			},
		},
		{
			name: "works with relative files & directories",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
				Recursive: true,
			},
			inputs: []string{
				"testdata/file.c",
				"testdata/sub",
			},
			want: map[string]string{
				"testdata/file.c":      "testdata/file.c",
				"testdata/sub/file2.c": "testdata/sub/file2.c",
			},
		},
		{
			name: "works with absolute files",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
				Recursive: true,
			},
			inputs: []string{
				filepath.Join(pwd, "testdata/file.c"),
				filepath.Join(pwd, "testdata/sub/file2.c"),
			},
			want: map[string]string{
				filepath.Join(pwd, "testdata/file.c"):      filepath.Join(pwd, "testdata/file.c"),
				filepath.Join(pwd, "testdata/sub/file2.c"): filepath.Join(pwd, "testdata/sub/file2.c"),
			},
		},
		{
			name: "works with absolute files & directories",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
				Recursive: true,
			},
			inputs: []string{
				filepath.Join(pwd, "testdata/file.c"),
				filepath.Join(pwd, "testdata/sub"),
			},
			want: map[string]string{
				filepath.Join(pwd, "testdata/file.c"):      filepath.Join(pwd, "testdata/file.c"),
				filepath.Join(pwd, "testdata/sub/file2.c"): filepath.Join(pwd, "testdata/sub/file2.c"),
			},
		},
		{
			name: "don't reprocess same file twice",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
				Recursive: true,
			},
			inputs: []string{
				"testdata/file.c",
				filepath.Join(pwd, "testdata/file.c"),
			},
			want: map[string]string{
				"testdata/file.c": "testdata/file.c",
			},
		},
		{
			name: "report processor failure",
			config: contracts.FilesProcessorConfig[string]{
				Processor: func(filepath string) (string, error) {
					return "", fmt.Errorf("failed")
				},
				Recursive: true,
			},
			inputs: []string{
				"testdata/file.c",
			},
			wantErr: true,
		},
		{
			name: "fail with missing files",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
				Recursive: true,
			},
			inputs: []string{
				"testdata/unknown.c",
			},
			wantErr: true,
		},
		{
			name: "fail with directories without a flag",
			config: contracts.FilesProcessorConfig[string]{
				Processor: echoProcessor,
			},
			inputs: []string{
				"testdata/sub",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			processor, err := New(tt.config)
			if err != nil {
				t.Fatalf("could not create processor: %v", err)
			}

			got, err := processor.ProcessFiles(tt.inputs)
			if (err != nil) != tt.wantErr {
				t.Errorf("error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
