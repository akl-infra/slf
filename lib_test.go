package slf

import (
	"encoding/json"
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
	for i := 0; i < 7; i++ {
		expected_row := strings.TrimSpace(expected_lines[i])
		out_row := strings.TrimSpace(out_lines[i])
		if expected_row != out_row {
			t.Fatalf("line %d did not match expected value\n  expected: `%s`\n  got: `%s`", i, expected_row, out_row)
		}
	}
}

func TestKeymeowConversion(t *testing.T) {
	l, _ := ReadLayoutFile("test_data/qwerty.json")
	out := l.ToKeymeow()
	b, err := os.ReadFile("test_data/qwerty.keymeow.json")
	if err != nil {
		t.Fatalf("error reading keymeow test data: %v", err)
	}
	var expected KeymeowLayout
	err = json.Unmarshal(b, &expected)
	if err != nil {
		t.Fatalf("error deserializing keymeow test data: %v", err)
	}
	if out.Name != expected.Name {
		t.Fatalf("l.Name should be `QWERTY`, got `%s`", out.Name)
	}
	for i, out_comp := range out.Components {
		expected_comp := expected.Components[i]
		fingers_match := out_comp.Finger[0] == expected_comp.Finger[0]
		keys_match := strings.Join(out_comp.Keys, "") == strings.Join(expected_comp.Keys, "")
		if !fingers_match || !keys_match {
			t.Fatalf("component %d did not match expected value\n  expected:  `%v`\n  got: `%v`", i, expected_comp, out_comp)
		}
	}
}
