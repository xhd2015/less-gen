package flags

import "strings"

// ParseIndex parses a flag from the args slice.
//
// Usage example:
//
// import (
//
//	"fmt"
//	"github.com/xhd2015/less-gen/flags"
//
// )
//
//	func handle(args []string) error {
//	  n := len(args)
//	  var remainArgs []string
//	  for i := 0; i < n; i++ {
//	    flag, value := flags.ParseIndex(args, &i)
//	    if flag == "" {
//	      remainArgs = append(remainArgs, args[i])
//	      continue
//	    }
//	    switch flag {
//	    case "-t", "--timeout":
//	      value, ok := value()
//	      if !ok {
//	        return fmt.Errorf("%s requires a value", flag)
//	      }
//	      _ = value
//	    // ...
//	    default:
//	      return fmt.Errorf("unknown flag: %s", flag)
//	    }
//	  }
//	  return nil
//	}
func ParseIndex(args []string, i *int) (string, func() (string, bool)) {
	arg := args[*i]
	if !strings.HasPrefix(arg, "-") {
		return "", nil
	}
	idx := strings.Index(arg[1:], "=")

	flag := arg
	value := func() (string, bool) {
		if *i+1 >= len(args) {
			return "", false
		}
		*i++
		return args[*i], true
	}
	if idx >= 0 {
		idx++
		flag = arg[:idx]
		value = func() (string, bool) {
			return arg[idx+1:], true
		}
	}
	return flag, value
}
