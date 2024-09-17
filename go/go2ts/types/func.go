package types

import (
	"fmt"
	"go/types"
)

func ValidateRespErr(res *types.Tuple) error {
	if res.Len() == 0 {
		return fmt.Errorf("function no result")
	}
	if res.Len() > 2 {
		return fmt.Errorf("function returns more than 2 results")
	}
	lastRes := res.At(res.Len() - 1)
	var lastResIsErr bool
	if named, ok := lastRes.Type().(*types.Named); ok && named.Obj().Name() == "error" {
		lastResIsErr = true
	}
	if !lastResIsErr {
		return fmt.Errorf("last res is not err: %T %v", lastRes.Type(), lastRes.Type())
	}
	return nil
}
