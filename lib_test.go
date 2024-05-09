package slf

import (
	"os"
	"strings"
	"testing"
)

func TestParseLayout(t *testing.T) {
	l, err := ReadLayoutFile("test_data/qwerty.json")
	if err != nil {
		t.Fatalf("Error reading layout file: %v", err)
	}
	if l.Name != "QWERTY" {
		t.Fatalf("l.Name should be `QWERTY`, got %s", l.Name)
	}
}

func TestGenkeyConversion(t *testing.T) {
	l, _ := ReadLayoutFile("test_data/qwerty.json")
	converted, err := l.ToGenkey()
	if err != nil {
		t.Fatalf("l.ToGenkey() should have no error for qwerty, but got %v", err)
	}
	out_lines := strings.Split(converted, "\n")
	b, err := os.ReadFile("test_data/qwerty.genkey")
	if err != nil {
		t.Fatalf("error reading genkey test data: %v", err)
	}
	expected_lines := strings.Split(string(b), "\n")
	for i:=0;i<7;i++ {
		expected_row := strings.TrimSpace(expected_lines[i])
		out_row := strings.TrimSpace(out_lines[i])
		if expected_row != out_row {
			t.Fatalf("line %d did not match expected value\n  expected: `%s`\n  got: `%s`", i, expected_row, out_row)
		}
	}
}

