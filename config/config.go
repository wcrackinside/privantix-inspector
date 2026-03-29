package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	SupportedExtensions []string
	MaxSampleRows       int
	Workers             int
	Security            bool
	Checksum            bool
	Recursive           bool
	RecursiveArchives   bool
	Incremental         bool
	HideSampleValues    bool
}

func Default() Config {
	return Config{
		SupportedExtensions: []string{".csv", ".txt", ".tsv", ".xlsx", ".rdat", ".dat", ".parquet", ".7z", ".zip", ".rar"},
		MaxSampleRows:       200,
		Workers:             4,
		Security:            false,
		Checksum:            false,
		Recursive:           true,
		RecursiveArchives:   false,
		Incremental:         false,
		HideSampleValues:    false,
	}
}

// Load parses a very small YAML-like config format for MVP purposes.
func Load(path string) (Config, error) {
	cfg := Default()
	if path == "" {
		return cfg, nil
	}
	f, err := os.Open(path)
	if err != nil {
		return cfg, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	currentKey := ""
	var exts []string
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if strings.HasPrefix(line, "-") && currentKey == "supported_extensions" {
			ext := strings.TrimSpace(strings.TrimPrefix(line, "-"))
			ext = strings.Trim(ext, `"'"`)
			exts = append(exts, ext)
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		val := strings.TrimSpace(parts[1])
		currentKey = key
		switch key {
		case "max_sample_rows":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.MaxSampleRows = n
			}
		case "workers":
			if n, err := strconv.Atoi(val); err == nil && n > 0 {
				cfg.Workers = n
			}
		case "supported_extensions":
			if val != "" {
				items := strings.Split(val, ",")
				cfg.SupportedExtensions = nil
				for _, item := range items {
					item = strings.Trim(strings.TrimSpace(item), `"'\[]`)
					if item != "" {
						cfg.SupportedExtensions = append(cfg.SupportedExtensions, item)
					}
				}
			}
		case "hide_sample_values":
			cfg.HideSampleValues = strings.EqualFold(val, "true") || val == "1"
		case "recursive_archives":
			cfg.RecursiveArchives = strings.EqualFold(val, "true") || val == "1"
		case "incremental":
			cfg.Incremental = strings.EqualFold(val, "true") || val == "1"
		}
	}
	if len(exts) > 0 {
		cfg.SupportedExtensions = exts
	}
	return cfg, scanner.Err()
}
