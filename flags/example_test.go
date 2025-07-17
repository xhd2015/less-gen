package flags

import (
	"fmt"
	"testing"
	"time"
)

func ExampleBuilder() {
	// Variables to store the parsed values
	var verbose string
	var timeout time.Duration
	var port int
	var files []string

	// Help function
	help := func() {
		fmt.Println("Usage: myapp [options] [arguments]")
		fmt.Println("Options:")
		fmt.Println("  -v, --verbose  Set verbosity level")
		fmt.Println("  --timeout      Set timeout duration")
		fmt.Println("  --port         Set port number")
		fmt.Println("  --files        Add file to process")
		fmt.Println("  --help         Show this help message")
	}

	// Sample arguments
	args := []string{
		"-v", "debug",
		"--timeout", "30s",
		"--port", "8080",
		"--files", "file1.txt",
		"--files", "file2.txt",
		"remaining1", "remaining2",
	}

	// Parse using fluent builder pattern exactly as requested
	remainArgs, err := String("-v,--verbose", &verbose).
		Duration("--timeout", &timeout).
		Int("--port", &port).
		StringSlice("--files", &files).
		HelpFunc("--help", help).
		Parse(args)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display results
	fmt.Printf("Verbose: %s\n", verbose)
	fmt.Printf("Timeout: %v\n", timeout)
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("Files: %v\n", files)
	fmt.Printf("Remaining args: %v\n", remainArgs)

	// Output:
	// Verbose: debug
	// Timeout: 30s
	// Port: 8080
	// Files: [file1.txt file2.txt]
	// Remaining args: [remaining1 remaining2]
}

func ExampleBuilder_Help() {
	// Variables to store the parsed values
	var verbose string
	var timeout time.Duration
	var port int

	// Sample arguments
	args := []string{
		"-v", "debug",
		"--timeout", "30s",
		"--port", "8080",
		"remaining",
	}

	// Help text as string (new feature)
	helpText := `Usage: myapp [options] [arguments]
Options:
  -v, --verbose  Set verbosity level
  --timeout      Set timeout duration  
  --port         Set port number
  --help         Show this help message`

	// Parse using string-based help
	remainArgs, err := String("-v,--verbose", &verbose).
		Duration("--timeout", &timeout).
		Int("--port", &port).
		Help("--help", helpText).
		Parse(args)

	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	// Display results
	fmt.Printf("Verbose: %s\n", verbose)
	fmt.Printf("Timeout: %v\n", timeout)
	fmt.Printf("Port: %d\n", port)
	fmt.Printf("Remaining args: %v\n", remainArgs)

	// Output:
	// Verbose: debug
	// Timeout: 30s
	// Port: 8080
	// Remaining args: [remaining]
}

func TestExampleWithDoublePointers(t *testing.T) {
	// Test with double pointers as mentioned in requirements
	var verbose *string
	var timeout *time.Duration
	var port *int
	var files *[]string

	args := []string{
		"--verbose", "debug",
		"--timeout", "1m",
		"--port", "9090",
		"--files", "test.txt",
		"remaining",
	}

	remainArgs, err := String("--verbose", &verbose).
		Duration("--timeout", &timeout).
		Int("--port", &port).
		StringSlice("--files", &files).
		Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	// Verify double pointers work correctly
	if verbose == nil || *verbose != "debug" {
		t.Errorf("Expected verbose='debug', got %v", verbose)
	}
	if timeout == nil || *timeout != time.Minute {
		t.Errorf("Expected timeout=1m, got %v", timeout)
	}
	if port == nil || *port != 9090 {
		t.Errorf("Expected port=9090, got %v", port)
	}
	if files == nil || len(*files) != 1 || (*files)[0] != "test.txt" {
		t.Errorf("Expected files=['test.txt'], got %v", files)
	}
	if len(remainArgs) != 1 || remainArgs[0] != "remaining" {
		t.Errorf("Expected remainArgs=['remaining'], got %v", remainArgs)
	}
}
