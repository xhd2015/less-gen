package git

import "github.com/xhd2015/xgo/support/cmd"

func GetToplevel() (string, error) {
	return cmd.Output("git", "rev-parse", "--show-toplevel")
}
