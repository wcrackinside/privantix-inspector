package tests

import (
	"testing"
	"privantix-source-inspector/detectors"
)

func TestDetectDelimiter(t *testing.T) {
	cases := []struct {
		name     string
		lines    []string
		expected string
	}{
		{
			"comma separated",
			[]string{"id,name,email", "1,Ana,ana@example.com", "2,Luis,luis@example.com"},
			",",
		},
		{
			"semicolon separated",
			[]string{"id;name;email", "1;Ana;ana@example.com", "2;Luis;luis@example.com"},
			";",
		},
		{
			"tab separated",
			[]string{"id\tname\temail", "1\tAna\tana@example.com", "2\tLuis\tluis@example.com"},
			"\t",
		},
		{
			"pipe separated",
			[]string{"id|name|email", "1|Ana|ana@example.com", "2|Luis|luis@example.com"},
			"|",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := detectors.DetectDelimiter(tc.lines)
			if got != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, got)
			}
		})
	}
}

func TestInferType(t *testing.T) {
	cases := map[string]string{
		"ana@example.com":      "email",
		"+56911112222":         "phone",
		"123":                  "integer",
		"45.6":                 "float",
		"45,6":                 "float",
		"2026-03-05":           "date",
		"05/03/2026":           "date",
		"2026-03-05T10:00:00Z": "datetime",
		"2026-03-05 10:00:00":  "datetime",
		"true":                 "boolean",
		"FALSE":                "boolean",
		"":                     "string",
		"NULL":                 "string",
		"N/A":                  "string",
		"ABC-1234":             "id_like",
		"regular string text":  "string",
	}
	for value, expected := range cases {
		t.Run(value, func(t *testing.T) {
			if got := detectors.InferType(value); got != expected {
				t.Errorf("value=%q expected=%q got=%q", value, expected, got)
			}
		})
	}
}

func TestLooksLikeHeader(t *testing.T) {
	cases := []struct {
		name     string
		first    []string
		second   []string
		expected bool
	}{
		{
			"typical header",
			[]string{"id", "name", "date_of_birth"},
			[]string{"1", "Ana", "2000-01-01"},
			true,
		},
		{
			"no header, just data",
			[]string{"1", "Ana", "2000-01-01"},
			[]string{"2", "Luis", "1990-05-15"},
			false,
		},
		{
			"single row header assumed true if all text",
			[]string{"name", "address", "city"},
			[]string{},
			true,
		},
		{
			"single row data assumed false if not text",
			[]string{"1", "2.5", "2020-01-01"},
			[]string{},
			false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := detectors.LooksLikeHeader(tc.first, tc.second)
			if got != tc.expected {
				t.Errorf("expected %v, got %v", tc.expected, got)
			}
		})
	}
}
