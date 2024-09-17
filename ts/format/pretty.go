package format

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"

	"github.com/xhd2015/xgo/support/cmd"
)

// Pretty pretty ts code via
// command: prettier --parser babel-ts --config pretty-ts.json
// prerequisity: npm install -g prettier
func Pretty(code string) (string, error) {
	output, err := cmd.New().Stdin(strings.NewReader(code)).Output("prettier", "--parser=babel-ts", "--tab-width=4", "--no-semi")
	if err != nil {
		var execErr *exec.Error
		if errors.As(err, &execErr) && execErr.Err == exec.ErrNotFound {
			return "", fmt.Errorf("prettier not found, try run: npm install -g prettier")
		}
		return "", err
	}
	return output, nil
}
