package utils

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode/utf8"
)

func NormalizeExtensions(items []string) []string {
	set := map[string]struct{}{}
	for _, item := range items {
		item = strings.TrimSpace(strings.ToLower(item))
		if item == "" {
			continue
		}
		if !strings.HasPrefix(item, ".") {
			item = "." + item
		}
		set[item] = struct{}{}
	}
	result := make([]string, 0, len(set))
	for k := range set {
		result = append(result, k)
	}
	sort.Strings(result)
	return result
}

func DetectEncoding(sample []byte) string {
	if len(sample) >= 3 && sample[0] == 0xEF && sample[1] == 0xBB && sample[2] == 0xBF {
		return "utf-8-bom"
	}
	if len(sample) >= 2 {
		if sample[0] == 0xFF && sample[1] == 0xFE {
			return "utf-16le-bom"
		}
		if sample[0] == 0xFE && sample[1] == 0xFF {
			return "utf-16be-bom"
		}
	}
	if utf8.Valid(sample) {
		return "utf-8"
	}
	return "windows-1252"
}

func EnsureDir(path string) error {
	return os.MkdirAll(path, 0o755)
}

func FileDepth(root, path string) int {
	rel, err := filepath.Rel(root, filepath.Dir(path))
	if err != nil || rel == "." {
		return 0
	}
	parts := strings.Split(rel, string(filepath.Separator))
	return len(parts)
}
