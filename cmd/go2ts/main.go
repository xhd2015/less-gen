package main

import (
	"fmt"
	"os"
	"strings"
)

const help = `
go2ts help to transpile go code to ts code

Usage: go2ts <cmd> [OPTIONS]
Options:
  --help   show help message
`

func main() {
	err := handle(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
func handle(args []string) error {
	var some string

	var remainArgs []string
	n := len(args)
	for i := 0; i < n; i++ {
		if args[i] == "--some" {
			if i+1 >= n {
				return fmt.Errorf("%v requires arg", args[i])
			}
			some = args[i+1]
			i++
			continue
		}
		if args[i] == "--help" {
			fmt.Println(strings.TrimSpace(help))
			return nil
		}
		if args[i] == "--" {
			remainArgs = append(remainArgs, args[i+1:]...)
			break
		}
		if strings.HasPrefix(args[i], "-") {
			return fmt.Errorf("unrecognized flag: %v", args[i])
		}
		remainArgs = append(remainArgs, args[i])
	}
	// TODO handle
	_ = some

	return nil
}
