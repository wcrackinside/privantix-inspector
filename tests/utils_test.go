package tests

import (
	"testing"
	"privantix-source-inspector/utils"
)

func TestDetectEncoding(t *testing.T) {
	cases := []struct {
		name     string
		sample   []byte
		expected string
	}{
		{"utf-8-bom", []byte{0xEF, 0xBB, 0xBF, 'a', 'b', 'c'}, "utf-8-bom"},
		{"utf-16le-bom", []byte{0xFF, 0xFE, 'a', 0x00}, "utf-16le-bom"},
		{"utf-16be-bom", []byte{0xFE, 0xFF, 0x00, 'a'}, "utf-16be-bom"},
		{"utf-8", []byte("hello world"), "utf-8"},
		{"windows-1252", []byte{0xC0, 0xC1, 'a'}, "windows-1252"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := utils.DetectEncoding(tc.sample); got != tc.expected {
				t.Errorf("expected %s, got %s", tc.expected, got)
			}
		})
	}
}
