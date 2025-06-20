package flags

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"strconv"
	"strings"
	"time"
)

// Builder represents a fluent flag parser builder
type Builder struct {
	flagSpecs  []FlagSpec
	helpFunc   func()
	helpText   string
	helpNoExit bool
	prefixOnly bool
}

// FlagSpec represents a single flag specification
type FlagSpec struct {
	Names    []string    // flag names like ["-v", "--verbose"]
	Target   interface{} // pointer to the target variable
	Type     FlagType    // type of the flag
	HelpText string      // help text for this flag
}

// FlagType represents the type of flag
type FlagType int

const (
	FlagTypeString FlagType = iota
	FlagTypeDuration
	FlagTypeInt
	FlagTypeStringSlice
	FlagTypeBool
)

var ErrHelp = errors.New("help")

func New() *Builder {
	return &Builder{}
}

// String adds a string flag to the builder
func String(names string, target interface{}) *Builder {
	return (&Builder{}).String(names, target)
}
func Bool(names string, target interface{}) *Builder {
	return (&Builder{}).Bool(names, target)
}

// Duration creates a new builder and adds a duration flag
func Duration(names string, target interface{}) *Builder {
	return (&Builder{}).Duration(names, target)
}

// Int creates a new builder and adds an integer flag
func Int(names string, target interface{}) *Builder {
	return (&Builder{}).Int(names, target)
}

// StringSlice creates a new builder and adds a string slice flag
func StringSlice(names string, target interface{}) *Builder {
	return (&Builder{}).StringSlice(names, target)
}

// Help creates a new builder and adds a help flag with help text
func Help(names string, helpText string) *Builder {
	return (&Builder{}).Help(names, helpText)
}

// HelpFunc creates a new builder and adds a help flag with a custom help function
func HelpFunc(names string, helpFunc func()) *Builder {
	return (&Builder{}).HelpFunc(names, helpFunc)
}

// String adds a string flag to the builder
func (b *Builder) String(names string, target interface{}) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names:  flagNames,
		Target: target,
		Type:   FlagTypeString,
	})
	return b
}

func (b *Builder) Bool(names string, target interface{}) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names:  flagNames,
		Target: target,
		Type:   FlagTypeBool,
	})
	return b
}

// Duration adds a duration flag to the builder
func (b *Builder) Duration(names string, target interface{}) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names:  flagNames,
		Target: target,
		Type:   FlagTypeDuration,
	})
	return b
}

// Int adds an integer flag to the builder
func (b *Builder) Int(names string, target interface{}) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names:  flagNames,
		Target: target,
		Type:   FlagTypeInt,
	})
	return b
}

// StringSlice adds a string slice flag to the builder
func (b *Builder) StringSlice(names string, target interface{}) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names:  flagNames,
		Target: target,
		Type:   FlagTypeStringSlice,
	})
	return b
}

// Help adds a help flag to the builder with help text
func (b *Builder) Help(names string, helpText string) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names: flagNames,
		Type:  FlagTypeBool,
	})
	b.helpText = helpText
	return b
}

// HelpFunc adds a help flag to the builder with a custom help function
func (b *Builder) HelpFunc(names string, helpFunc func()) *Builder {
	flagNames := parseNames(names)
	b.flagSpecs = append(b.flagSpecs, FlagSpec{
		Names: flagNames,
		Type:  FlagTypeBool,
	})
	b.helpFunc = helpFunc
	return b
}
func (b *Builder) HelpNoExit() *Builder {
	b.helpNoExit = true
	return b
}

// PrefixOnly stops parsing flags after the first non-flag argument
func (b *Builder) PrefixOnly() *Builder {
	b.prefixOnly = true
	return b
}

// Parse parses the arguments using the configured flags
func (b *Builder) Parse(args []string) ([]string, error) {
	var remainArgs []string
	n := len(args)

	for i := 0; i < n; i++ {
		if args[i] == "--" {
			remainArgs = append(remainArgs, args[i+1:]...)
			break
		}
		flag, getValue := parseIndex(args, &i)
		if flag == "" {
			if b.prefixOnly {
				remainArgs = append(remainArgs, args[i:]...)
				break
			}
			remainArgs = append(remainArgs, args[i])
			continue
		}

		spec := b.findFlagSpec(flag)
		if spec == nil {
			return nil, fmt.Errorf("unrecognized flag: %s", flag)
		}

		// Handle help flag
		if spec.Type == FlagTypeBool && (b.helpFunc != nil || b.helpText != "") {
			for _, name := range spec.Names {
				if name == flag {
					if b.helpFunc != nil {
						b.helpFunc()
					} else if b.helpText != "" {
						txt := strings.TrimPrefix(b.helpText, "\n")
						fmt.Print(txt)
						if !strings.HasSuffix(txt, "\n") {
							fmt.Println()
						}
					}
					if !b.helpNoExit {
						os.Exit(0)
					}
					return remainArgs, ErrHelp
				}
			}
		}

		// Get the value for non-bool flags
		value, hasValue := getValue(spec.Type == FlagTypeBool)
		if !hasValue {
			return nil, fmt.Errorf("%s requires a value", flag)
		}

		// Set the value based on type
		err := b.setValue(spec, value)
		if err != nil {
			return nil, fmt.Errorf("error setting value for %s: %v", flag, err)
		}
	}

	return remainArgs, nil
}

// findFlagSpec finds the flag specification for a given flag name
func (b *Builder) findFlagSpec(flagName string) *FlagSpec {
	for i := range b.flagSpecs {
		spec := &b.flagSpecs[i]
		for _, name := range spec.Names {
			if name == flagName {
				return spec
			}
		}
	}
	return nil
}

// setValue sets the value to the target based on the flag type
func (b *Builder) setValue(spec *FlagSpec, value string) error {
	target := spec.Target

	switch spec.Type {
	case FlagTypeBool:
		return setBoolValue(target, value)
	case FlagTypeString:
		return setStringValue(target, value)
	case FlagTypeDuration:
		return setDurationValue(target, value)
	case FlagTypeInt:
		return setIntValue(target, value)
	case FlagTypeStringSlice:
		return setStringSliceValue(target, value)
	default:
		return fmt.Errorf("unsupported flag type")
	}
}

func setBoolValue(target interface{}, value string) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}
	elem := v.Elem()
	if elem.Kind() == reflect.Bool {
		elem.SetBool(value == "" || value == "true")
	} else {
		return fmt.Errorf("target must be a pointer to a bool")
	}
	return nil
}

// setStringValue sets a string value to a target (supports *string and **string)
func setStringValue(target interface{}, value string) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := v.Elem()
	if elem.Kind() == reflect.Ptr {
		// Handle **string
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
		elem.Elem().SetString(value)
	} else if elem.Kind() == reflect.String {
		// Handle *string
		elem.SetString(value)
	} else {
		return fmt.Errorf("target must be *string or **string")
	}

	return nil
}

// setDurationValue sets a duration value to a target
func setDurationValue(target interface{}, value string) error {
	duration, err := time.ParseDuration(value)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := v.Elem()
	if elem.Kind() == reflect.Ptr {
		// Handle **time.Duration
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
		elem.Elem().Set(reflect.ValueOf(duration))
	} else if elem.Type() == reflect.TypeOf(time.Duration(0)) {
		// Handle *time.Duration
		elem.Set(reflect.ValueOf(duration))
	} else {
		return fmt.Errorf("target must be *time.Duration or **time.Duration")
	}

	return nil
}

// setIntValue sets an integer value to a target
func setIntValue(target interface{}, value string) error {
	intVal, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := v.Elem()
	if elem.Kind() == reflect.Ptr {
		// Handle **int
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
		elem.Elem().SetInt(int64(intVal))
	} else if elem.Kind() == reflect.Int {
		// Handle *int
		elem.SetInt(int64(intVal))
	} else {
		return fmt.Errorf("target must be *int or **int")
	}

	return nil
}

// setStringSliceValue sets a string slice value to a target
func setStringSliceValue(target interface{}, value string) error {
	v := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr {
		return fmt.Errorf("target must be a pointer")
	}

	elem := v.Elem()
	if elem.Kind() == reflect.Ptr {
		// Handle **[]string
		if elem.IsNil() {
			elem.Set(reflect.New(elem.Type().Elem()))
		}
		slice := elem.Elem()
		newVal := reflect.Append(slice, reflect.ValueOf(value))
		slice.Set(newVal)
	} else if elem.Kind() == reflect.Slice && elem.Type().Elem().Kind() == reflect.String {
		// Handle *[]string
		newVal := reflect.Append(elem, reflect.ValueOf(value))
		elem.Set(newVal)
	} else {
		return fmt.Errorf("target must be *[]string or **[]string")
	}

	return nil
}

// parseNames parses comma-separated flag names
func parseNames(names string) []string {
	parts := strings.Split(names, ",")
	var result []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}
