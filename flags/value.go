package flags

import (
	"fmt"
	"strings"
)

func OnlyArg(args []string) (string, error) {
	if len(args) == 0 {
		return "", fmt.Errorf("requires an argument")
	}
	if len(args) > 1 {
		return "", fmt.Errorf("expect one argument, got extra: %s", strings.Join(args[1:], ","))
	}
	return args[0], nil
}
