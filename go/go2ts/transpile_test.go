package go2ts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/less-gen/ts/format"
	"github.com/xhd2015/xgo/support/assert"
	"github.com/xhd2015/xgo/support/cmd"
)

func TestTranspileFile(t *testing.T) {
	tests := []struct {
		file    string
		skip    string
		wantErr error
	}{
		{
			file: "hello/hello.go",
		},
		{
			file: "type/type.go",
			// skip: "not ready", //
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.file, func(t *testing.T) {
			if tt.skip != "" {
				t.Skipf("SKIP: %s", tt.skip)
				return
			}
			got, err := TranspileFile(filepath.Join("testdata", tt.file))
			if errDiff := assert.Diff(tt.wantErr, err); errDiff != "" {
				t.Errorf("TranspileFile() err: %s", errDiff)
				return
			}

			gotPretty, err := format.Pretty(got)
			if err != nil {
				t.Errorf("pretty %v: %s", err, got)
				return
			}

			expectFile := filepath.Join("testdata", strings.ReplaceAll(tt.file, ".go", ".ts"))
			expectData, err := os.ReadFile(expectFile)
			if err != nil {
				t.Error(err)
				return
			}

			if resultDiff := assert.Diff(string(expectData), gotPretty); resultDiff != "" {
				t.Errorf("TranspileFile(): %v", resultDiff)
			}
		})
	}
}

func getToplevel() (string, error) {
	return cmd.Output("git", "rev-parse", "--show-toplevel")
}
