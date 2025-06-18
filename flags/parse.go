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
//	      return fmt.Errorf("unrecognized flag: %s", flag)
//	    }
//	  }
//	  return nil
//	}
func ParseIndex(args []string, i *int) (string, func() (string, bool)) {
	flag, fn := parseIndex(args, i)
	var cb func() (string, bool)
	if fn != nil {
		cb = func() (string, bool) {
			return fn(false)
		}
	}
	return flag, cb
}

func parseIndex(args []string, i *int) (string, func(boolOnly bool) (string, bool)) {
	arg := args[*i]
	if !strings.HasPrefix(arg, "-") {
		return "", nil
	}
	eqIdx := strings.Index(arg[1:], "=")

	flag := arg
	value := func(boolOnly bool) (string, bool) {
		if boolOnly {
			return "", true
		}
		if *i+1 >= len(args) {
			return "", false
		}
		*i++
		return args[*i], true
	}
	if eqIdx >= 0 {
		eqIdx++
		flag = arg[:eqIdx]
		value = func(boolOnly bool) (string, bool) {
			return arg[eqIdx+1:], true
		}
	}
	return flag, value
}
