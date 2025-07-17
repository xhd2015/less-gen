# Flags Package

A fluent flag parsing library for Go that supports chaining and various data types.

# Installation
```sh
go get github.com/xhd2015/less-gen/flags@latest
```

## Features

- Fluent builder pattern for easy flag configuration
- Support for multiple data types: `bool`, `string`, `int`, `int64`, `time.Duration`, `[]string`
- Pointer support: `*T` and `**T` for all types
- Help text and custom help functions
- Multiple flag names (e.g., `-v,--verbose`)

## Quick Start

```go
var verbose bool
var timeout time.Duration
var files []string

const help = `
Usage: myapp [options]

Options:
  --timeout DURATION  set timeout duration
  --file FILE         add files to process, can be repeated
  -v, --verbose       enable verbose output
  -h, --help          show help
`

remainArgs, err := flags.Duration("--timeout", &timeout).
    StringSlice("--files", &files).
    Bool("-v,--verbose", &verbose).
    Help("-h,--help", help).
    Parse(os.Args[1:])
```

## Supported Types

- `*bool`, `**bool` - Boolean flags
- `*string`, `**string` - String values
- `*int`, `**int`, `*int64`, `**int64` - Integer values
- `*time.Duration`, `**time.Duration` - Duration values
- `*[]string`, `**[]string` - String slices (can be repeated)