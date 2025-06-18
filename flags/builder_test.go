package flags

import (
	"strings"
	"testing"
	"time"
)

func TestBuilder_Bool(t *testing.T) {
	var verbose bool
	args := []string{"--verbose", "true"}

	remainArgs, err := Bool("-v,--verbose", &verbose).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if !verbose {
		t.Errorf("Expected verbose=true, got %v", verbose)
	}
	if len(remainArgs) != 1 {
		t.Errorf("Expected 1 remaining args, got %v", remainArgs)
	}
	if remainArgs[0] != "true" {
		t.Errorf("Expected remainArgs=['true'], got %v", remainArgs)
	}
}

func TestBuilder_String(t *testing.T) {
	var verbose string
	args := []string{"--verbose", "debug", "remaining"}

	remainArgs, err := String("-v,--verbose", &verbose).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if verbose != "debug" {
		t.Errorf("Expected verbose='debug', got '%s'", verbose)
	}
	if len(remainArgs) != 1 || remainArgs[0] != "remaining" {
		t.Errorf("Expected remainArgs=['remaining'], got %v", remainArgs)
	}
}

func TestBuilder_StringDoublePointer(t *testing.T) {
	var verbose *string
	args := []string{"--verbose", "debug"}

	remainArgs, err := String("--verbose", &verbose).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if verbose == nil || *verbose != "debug" {
		t.Errorf("Expected verbose='debug', got %v", verbose)
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_Int(t *testing.T) {
	var port int
	args := []string{"--port", "8080", "file.txt"}

	remainArgs, err := String("--unused", new(string)).Int("--port", &port).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if port != 8080 {
		t.Errorf("Expected port=8080, got %d", port)
	}
	if len(remainArgs) != 1 || remainArgs[0] != "file.txt" {
		t.Errorf("Expected remainArgs=['file.txt'], got %v", remainArgs)
	}
}

func TestBuilder_IntDoublePointer(t *testing.T) {
	var port *int
	args := []string{"--port", "9090"}

	remainArgs, err := String("--unused", new(string)).Int("--port", &port).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if port == nil || *port != 9090 {
		t.Errorf("Expected port=9090, got %v", port)
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_Duration(t *testing.T) {
	var timeout time.Duration
	args := []string{"--timeout", "30s"}

	remainArgs, err := String("--unused", new(string)).Duration("--timeout", &timeout).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if timeout != 30*time.Second {
		t.Errorf("Expected timeout=30s, got %v", timeout)
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_DurationDoublePointer(t *testing.T) {
	var timeout *time.Duration
	args := []string{"--timeout", "1m"}

	remainArgs, err := String("--unused", new(string)).Duration("--timeout", &timeout).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if timeout == nil || *timeout != 1*time.Minute {
		t.Errorf("Expected timeout=1m, got %v", timeout)
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_StringSlice(t *testing.T) {
	var files []string
	args := []string{"--files", "file1.txt", "--files", "file2.txt", "remaining"}

	remainArgs, err := String("--unused", new(string)).StringSlice("--files", &files).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	expected := []string{"file1.txt", "file2.txt"}
	if len(files) != len(expected) {
		t.Errorf("Expected files length=%d, got %d", len(expected), len(files))
	}
	for i, f := range files {
		if f != expected[i] {
			t.Errorf("Expected files[%d]='%s', got '%s'", i, expected[i], f)
		}
	}
	if len(remainArgs) != 1 || remainArgs[0] != "remaining" {
		t.Errorf("Expected remainArgs=['remaining'], got %v", remainArgs)
	}
}

func TestBuilder_StringSliceDoublePointer(t *testing.T) {
	var files *[]string
	args := []string{"--files", "file1.txt", "--files", "file2.txt"}

	remainArgs, err := String("--unused", new(string)).StringSlice("--files", &files).Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}
	if files == nil {
		t.Fatal("Expected files to be initialized, got nil")
	}
	expected := []string{"file1.txt", "file2.txt"}
	if len(*files) != len(expected) {
		t.Errorf("Expected files length=%d, got %d", len(expected), len(*files))
	}
	for i, f := range *files {
		if f != expected[i] {
			t.Errorf("Expected files[%d]='%s', got '%s'", i, expected[i], f)
		}
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_Help(t *testing.T) {
	args := []string{"--help"}

	remainArgs, err := String("--unused", new(string)).Help("--help", "This is help text").HelpNoExit().Parse(args)

	if err != ErrHelp {
		t.Errorf("expected ErrHelp, got %v", err)
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_HelpFunc(t *testing.T) {
	var helpCalled bool
	helpFunc := func() {
		helpCalled = true
	}

	args := []string{"--help"}

	remainArgs, err := String("--unused", new(string)).HelpFunc("--help", helpFunc).HelpNoExit().Parse(args)

	if err != ErrHelp {
		t.Errorf("expected ErrHelp, got %v", err)
	}
	if !helpCalled {
		t.Error("Expected help function to be called")
	}
	if len(remainArgs) != 0 {
		t.Errorf("Expected no remaining args, got %v", remainArgs)
	}
}

func TestBuilder_ComplexExample(t *testing.T) {
	var verbose string
	var timeout time.Duration
	var port int
	var files []string
	var helpCalled bool

	helpFunc := func() {
		helpCalled = true
	}

	args := []string{
		"-v", "debug",
		"--timeout", "30s",
		"--port", "8080",
		"--files", "file1.txt",
		"--files", "file2.txt",
		"remaining1", "remaining2",
	}

	remainArgs, err := String("-v,--verbose", &verbose).
		Duration("--timeout", &timeout).
		Int("--port", &port).
		StringSlice("--files", &files).
		HelpFunc("--help", helpFunc).
		Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if verbose != "debug" {
		t.Errorf("Expected verbose='debug', got '%s'", verbose)
	}
	if timeout != 30*time.Second {
		t.Errorf("Expected timeout=30s, got %v", timeout)
	}
	if port != 8080 {
		t.Errorf("Expected port=8080, got %d", port)
	}
	expectedFiles := []string{"file1.txt", "file2.txt"}
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected files length=%d, got %d", len(expectedFiles), len(files))
	}
	for i, f := range files {
		if f != expectedFiles[i] {
			t.Errorf("Expected files[%d]='%s', got '%s'", i, expectedFiles[i], f)
		}
	}

	expectedRemaining := []string{"remaining1", "remaining2"}
	if len(remainArgs) != len(expectedRemaining) {
		t.Errorf("Expected remaining args length=%d, got %d", len(expectedRemaining), len(remainArgs))
	}
	for i, arg := range remainArgs {
		if arg != expectedRemaining[i] {
			t.Errorf("Expected remainArgs[%d]='%s', got '%s'", i, expectedRemaining[i], arg)
		}
	}

	if helpCalled {
		t.Error("Expected help function not to be called")
	}
}

func TestBuilder_UnrecognizedFlag(t *testing.T) {
	var verbose string
	args := []string{"--unknown", "value"}

	_, err := String("--verbose", &verbose).Parse(args)

	if err == nil {
		t.Fatal("Expected error for unrecognized flag")
	}
	if !strings.Contains(err.Error(), "unrecognized flag: --unknown") {
		t.Errorf("Expected 'unrecognized flag' error, got: %v", err)
	}
}

func TestBuilder_MissingValue(t *testing.T) {
	var verbose string
	args := []string{"--verbose"}

	_, err := String("--verbose", &verbose).Parse(args)

	if err == nil {
		t.Fatal("Expected error for missing value")
	}
	if !strings.Contains(err.Error(), "requires a value") {
		t.Errorf("Expected 'requires a value' error, got: %v", err)
	}
}

func TestBuilder_ParseNames(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{"--verbose", []string{"--verbose"}},
		{"-v,--verbose", []string{"-v", "--verbose"}},
		{"-v, --verbose, --debug", []string{"-v", "--verbose", "--debug"}},
		{" -v , --verbose ", []string{"-v", "--verbose"}},
	}

	for _, test := range tests {
		result := parseNames(test.input)
		if len(result) != len(test.expected) {
			t.Errorf("For input '%s', expected length %d, got %d", test.input, len(test.expected), len(result))
			continue
		}
		for i, name := range result {
			if name != test.expected[i] {
				t.Errorf("For input '%s', expected[%d]='%s', got '%s'", test.input, i, test.expected[i], name)
			}
		}
	}
}

func TestBuilder_EqualsFormat(t *testing.T) {
	var verbose string
	var timeout time.Duration
	var port int
	var files []string

	// Test with equals format: --flag=value
	args := []string{
		"--verbose=debug",
		"--timeout=30s",
		"--port=8080",
		"--files=file1.txt",
		"--files=file2.txt",
		"remaining1", "remaining2",
	}

	remainArgs, err := String("--verbose", &verbose).
		Duration("--timeout", &timeout).
		Int("--port", &port).
		StringSlice("--files", &files).
		Parse(args)

	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if verbose != "debug" {
		t.Errorf("Expected verbose='debug', got '%s'", verbose)
	}
	if timeout != 30*time.Second {
		t.Errorf("Expected timeout=30s, got %v", timeout)
	}
	if port != 8080 {
		t.Errorf("Expected port=8080, got %d", port)
	}
	expectedFiles := []string{"file1.txt", "file2.txt"}
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected files length=%d, got %d", len(expectedFiles), len(files))
	}
	for i, f := range files {
		if f != expectedFiles[i] {
			t.Errorf("Expected files[%d]='%s', got '%s'", i, expectedFiles[i], f)
		}
	}

	expectedRemaining := []string{"remaining1", "remaining2"}
	if len(remainArgs) != len(expectedRemaining) {
		t.Errorf("Expected remaining args length=%d, got %d", len(expectedRemaining), len(remainArgs))
	}
	for i, arg := range remainArgs {
		if arg != expectedRemaining[i] {
			t.Errorf("Expected remainArgs[%d]='%s', got '%s'", i, expectedRemaining[i], arg)
		}
	}
}

func TestBuilder_MixedFormats(t *testing.T) {
	var verbose string
	var timeout time.Duration
	var port int
	var files []string

	// Test with mixed format: some with equals, some with space
	args := []string{
		"--verbose=debug",
		"--timeout", "30s",
		"--port=8080",
		"--files", "file1.txt",
		"--files=file2.txt",
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

	if verbose != "debug" {
		t.Errorf("Expected verbose='debug', got '%s'", verbose)
	}
	if timeout != 30*time.Second {
		t.Errorf("Expected timeout=30s, got %v", timeout)
	}
	if port != 8080 {
		t.Errorf("Expected port=8080, got %d", port)
	}
	expectedFiles := []string{"file1.txt", "file2.txt"}
	if len(files) != len(expectedFiles) {
		t.Errorf("Expected files length=%d, got %d", len(expectedFiles), len(files))
	}
	for i, f := range files {
		if f != expectedFiles[i] {
			t.Errorf("Expected files[%d]='%s', got '%s'", i, expectedFiles[i], f)
		}
	}

	if len(remainArgs) != 1 || remainArgs[0] != "remaining" {
		t.Errorf("Expected remainArgs=['remaining'], got %v", remainArgs)
	}
}

func TestTopLevelFunctions(t *testing.T) {
	var verbose string
	var timeout time.Duration
	var port int
	var files []string

	// Test starting with Duration
	args1 := []string{"--timeout", "1m", "--verbose", "debug"}
	remainArgs1, err := Duration("--timeout", &timeout).String("--verbose", &verbose).Parse(args1)
	if err != nil {
		t.Fatalf("Duration start failed: %v", err)
	}
	if timeout != time.Minute || verbose != "debug" {
		t.Errorf("Duration start: expected timeout=1m, verbose=debug, got timeout=%v, verbose=%s", timeout, verbose)
	}
	if len(remainArgs1) != 0 {
		t.Errorf("Duration start: expected no remaining args, got %v", remainArgs1)
	}

	// Test starting with Int
	args2 := []string{"--port", "9090", "--verbose", "info"}
	verbose = "" // reset
	remainArgs2, err := Int("--port", &port).String("--verbose", &verbose).Parse(args2)
	if err != nil {
		t.Fatalf("Int start failed: %v", err)
	}
	if port != 9090 || verbose != "info" {
		t.Errorf("Int start: expected port=9090, verbose=info, got port=%d, verbose=%s", port, verbose)
	}
	if len(remainArgs2) != 0 {
		t.Errorf("Int start: expected no remaining args, got %v", remainArgs2)
	}

	// Test starting with StringSlice
	files = nil // reset
	args3 := []string{"--files", "file1.txt", "--files", "file2.txt", "--verbose", "warn"}
	verbose = "" // reset
	remainArgs3, err := StringSlice("--files", &files).String("--verbose", &verbose).Parse(args3)
	if err != nil {
		t.Fatalf("StringSlice start failed: %v", err)
	}
	expectedFiles := []string{"file1.txt", "file2.txt"}
	if len(files) != len(expectedFiles) || verbose != "warn" {
		t.Errorf("StringSlice start: expected files=%v, verbose=warn, got files=%v, verbose=%s", expectedFiles, files, verbose)
	}
	if len(remainArgs3) != 0 {
		t.Errorf("StringSlice start: expected no remaining args, got %v", remainArgs3)
	}
}

func TestTopLevelHelp(t *testing.T) {
	// Test starting with Help
	args1 := []string{"--help"}
	remainArgs1, err := Help("--help", "Test help text").HelpNoExit().Parse(args1)
	if err != ErrHelp {
		t.Errorf("expected ErrHelp, got %v", err)
	}
	if len(remainArgs1) != 0 {
		t.Errorf("Help start: expected no remaining args, got %v", remainArgs1)
	}

	// Test starting with HelpFunc
	var helpCalled bool
	helpFunc := func() { helpCalled = true }
	args2 := []string{"--help"}
	remainArgs2, err := HelpFunc("--help", helpFunc).HelpNoExit().Parse(args2)
	if err != ErrHelp {
		t.Errorf("expected ErrHelp, got %v", err)
	}
	if !helpCalled {
		t.Error("HelpFunc start: expected help function to be called")
	}
	if len(remainArgs2) != 0 {
		t.Errorf("HelpFunc start: expected no remaining args, got %v", remainArgs2)
	}
}

func TestTopLevelChaining(t *testing.T) {
	var verbose string
	var timeout time.Duration
	var port int
	var files []string

	// Test complex chaining starting with different entry points
	args := []string{
		"--port", "8080",
		"--timeout", "2m",
		"--verbose", "debug",
		"--files", "test.txt",
		"remaining",
	}

	// Start with Int, then chain everything else
	remainArgs, err := Int("--port", &port).
		Duration("--timeout", &timeout).
		String("--verbose", &verbose).
		StringSlice("--files", &files).
		Parse(args)

	if err != nil {
		t.Fatalf("Complex chaining failed: %v", err)
	}

	if port != 8080 {
		t.Errorf("Expected port=8080, got %d", port)
	}
	if timeout != 2*time.Minute {
		t.Errorf("Expected timeout=2m, got %v", timeout)
	}
	if verbose != "debug" {
		t.Errorf("Expected verbose='debug', got '%s'", verbose)
	}
	if len(files) != 1 || files[0] != "test.txt" {
		t.Errorf("Expected files=['test.txt'], got %v", files)
	}
	if len(remainArgs) != 1 || remainArgs[0] != "remaining" {
		t.Errorf("Expected remainArgs=['remaining'], got %v", remainArgs)
	}
}
