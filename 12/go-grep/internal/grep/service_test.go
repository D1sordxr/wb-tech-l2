package grep

import (
	"strings"
	"testing"
	"wb-tech-l2/12/go-grep/internal/config"
)

func TestService_ProcessLines_BasicMatch(t *testing.T) {
	cfg := &config.Grep{
		Pattern: "hello",
	}
	service := NewService(cfg)

	lines := []string{
		"first line",
		"hello world",
		"third line",
		"another hello",
	}

	result := service.ProcessLines(lines)
	expected := []string{"hello world", "another hello"}

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(result))
	}

	for i, line := range result {
		if line != expected[i] {
			t.Errorf("Expected '%s', got '%s'", expected[i], line)
		}
	}
}

func TestService_ProcessLines_IgnoreCase(t *testing.T) {
	cfg := &config.Grep{
		Pattern:    "HELLO",
		IgnoreCase: true,
	}
	service := NewService(cfg)

	lines := []string{
		"Hello world",
		"HELLO there",
		"goodbye",
	}

	result := service.ProcessLines(lines)
	expected := 2

	if len(result) != expected {
		t.Errorf("Expected %d matches, got %d", expected, len(result))
	}
}

func TestService_ProcessLines_InvertMatch(t *testing.T) {
	cfg := &config.Grep{
		Pattern:     "hello",
		InvertMatch: true,
	}
	service := NewService(cfg)

	lines := []string{
		"hello world",
		"goodbye",
		"another hello",
	}

	result := service.ProcessLines(lines)
	expected := []string{"goodbye"}

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines, got %d", len(expected), len(result))
	}
}

func TestService_ProcessLines_CountOnly(t *testing.T) {
	cfg := &config.Grep{
		Pattern:   "hello",
		CountOnly: true,
	}
	service := NewService(cfg)

	lines := []string{
		"hello world",
		"goodbye",
		"another hello",
	}

	result := service.ProcessLines(lines)
	expected := []string{"2"}

	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Expected count '%s', got '%s'", expected[0], result[0])
	}
}

func TestService_ProcessLines_LineNumbers(t *testing.T) {
	cfg := &config.Grep{
		Pattern:    "hello",
		LineNumber: true,
	}
	service := NewService(cfg)

	lines := []string{
		"first line",
		"hello world",
		"third line",
	}

	result := service.ProcessLines(lines)
	expected := []string{"2:hello world"}

	if len(result) != len(expected) || result[0] != expected[0] {
		t.Errorf("Expected '%s', got '%s'", expected[0], result[0])
	}
}

func TestService_ProcessLines_Context(t *testing.T) {
	cfg := &config.Grep{
		Pattern: "hello",
		Context: 1,
	}
	service := NewService(cfg)

	lines := []string{
		"line 1",
		"line 2",
		"hello world",
		"line 4",
		"line 5",
	}

	result := service.ProcessLines(lines)
	expected := []string{"line 2", "hello world", "line 4"}

	if len(result) != len(expected) {
		t.Errorf("Expected %d lines with context, got %d", len(expected), len(result))
	}
}

func TestService_ProcessLines_FixedString(t *testing.T) {
	cfg := &config.Grep{
		Pattern:     "hello.world",
		FixedString: true,
	}
	service := NewService(cfg)

	lines := []string{
		"hello.world",      // should match
		"hello world",      // should not match
		"hello.world test", // should match
	}

	result := service.ProcessLines(lines)
	expected := 2

	if len(result) != expected {
		t.Errorf("Expected %d fixed string matches, got %d", expected, len(result))
	}
}

func TestService_ProcessLines_Regexp(t *testing.T) {
	cfg := &config.Grep{
		Pattern: `h.llo`, // regexp pattern
	}
	service := NewService(cfg)

	lines := []string{
		"hello",   // match
		"hallo",   // match
		"hxllo",   // match
		"goodbye", // no match
	}

	result := service.ProcessLines(lines)
	expected := 3

	if len(result) != expected {
		t.Errorf("Expected %d regexp matches, got %d", expected, len(result))
	}
}

func TestBuildMatcher_FixedStringCaseSensitive(t *testing.T) {
	cfg := &config.Grep{
		Pattern:     "Hello",
		FixedString: true,
		IgnoreCase:  false,
	}
	service := NewService(cfg)

	matcher, err := service.buildMatcher()
	if err != nil {
		t.Fatalf("Failed to build matcher: %v", err)
	}

	if !matcher("Hello World") {
		t.Error("Should match exact case")
	}
	if matcher("hello world") {
		t.Error("Should not match different case")
	}
}

func TestBuildMatcher_FixedStringIgnoreCase(t *testing.T) {
	cfg := &config.Grep{
		Pattern:     "Hello",
		FixedString: true,
		IgnoreCase:  true,
	}
	service := NewService(cfg)

	matcher, err := service.buildMatcher()
	if err != nil {
		t.Fatalf("Failed to build matcher: %v", err)
	}

	testCases := []string{"Hello World", "hello world", "HELLO WORLD"}
	for _, tc := range testCases {
		if !matcher(tc) {
			t.Errorf("Should match '%s' with ignore case", tc)
		}
	}
}

func TestMatchLines_Invert(t *testing.T) {
	cfg := &config.Grep{
		Pattern:     "hello",
		InvertMatch: true,
	}
	service := NewService(cfg)

	lines := []string{"hello", "world", "hello again"}
	matcher := func(s string) bool { return strings.Contains(s, "hello") }

	matched, count := service.matchLines(lines, matcher)

	if count != 1 {
		t.Errorf("Expected 1 inverted match, got %d", count)
	}
	if !matched[1] { // "world" should be matched
		t.Error("Inverted match failed for 'world'")
	}
	if matched[0] || matched[2] { // "hello" lines should not be matched
		t.Error("Inverted match incorrectly matched hello lines")
	}
}
