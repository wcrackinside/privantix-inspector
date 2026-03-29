package detectors

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	emailRegex  = regexp.MustCompile(`^[A-Za-z0-9._%+\-]+@[A-Za-z0-9.\-]+\.[A-Za-z]{2,}$`)
	phoneRegex  = regexp.MustCompile(`^\+?[0-9 ()\-]{7,20}$`)
	idLikeRegex = regexp.MustCompile(`^[A-Za-z0-9._/\-]{6,20}$`)
)

func IsNull(value string) bool {
	trimmed := strings.TrimSpace(strings.ToLower(value))
	return trimmed == "" || trimmed == "null" || trimmed == "na" || trimmed == "n/a"
}

func InferType(value string) string {
	v := strings.TrimSpace(value)
	if IsNull(v) {
		return "string"
	}
	if emailRegex.MatchString(v) {
		return "email"
	}
	// Dates before phone: ISO dates (e.g. 2026-03-05) match phoneRegex digit/hyphen pattern.
	if _, ok := parseDateTime(v); ok {
		if strings.Contains(v, ":") {
			return "datetime"
		}
		return "date"
	}
	if phoneRegex.MatchString(v) {
		return "phone"
	}
	if _, err := strconv.Atoi(v); err == nil {
		return "integer"
	}
	if _, err := strconv.ParseFloat(strings.ReplaceAll(v, ",", "."), 64); err == nil && strings.ContainsAny(v, ".,") {
		return "float"
	}
	if strings.EqualFold(v, "true") || strings.EqualFold(v, "false") || v == "0" || v == "1" {
		return "boolean"
	}
	if idLikeRegex.MatchString(v) {
		return "id_like"
	}
	return "string"
}

func parseDateTime(v string) (time.Time, bool) {
	layouts := []string{
		time.RFC3339,
		"2006-01-02",
		"02-01-2006",
		"02/01/2006",
		"2006/01/02",
		"2006-01-02 15:04:05",
		"02/01/2006 15:04:05",
		"02-01-2006 15:04:05",
	}
	for _, layout := range layouts {
		if t, err := time.Parse(layout, v); err == nil {
			return t, true
		}
	}
	return time.Time{}, false
}

func DetectDelimiter(lines []string) string {
	candidates := []rune{',', ';', '\t', '|'}
	bestScore := -1
	best := ','
	for _, c := range candidates {
		score := 0
		lastCount := -1
		consistent := true
		for _, line := range lines {
			if strings.TrimSpace(line) == "" {
				continue
			}
			count := strings.Count(line, string(c))
			if count == 0 {
				consistent = false
				continue
			}
			if lastCount != -1 && count != lastCount {
				consistent = false
			}
			lastCount = count
			score += count
		}
		if consistent {
			score += 100
		}
		if score > bestScore {
			bestScore = score
			best = c
		}
	}
	return string(best)
}

func LooksLikeHeader(first, second []string) bool {
	if len(first) == 0 {
		return false
	}
	if len(second) == 0 {
		textLike := 0
		for _, v := range first {
			t := InferType(v)
			if t == "string" || t == "id_like" || t == "email" {
				textLike++
			}
		}
		return textLike >= len(first)/2
	}

	score := 0.0
	for i := range first {
		t1 := InferType(first[i])
		if t1 == "string" || t1 == "id_like" || t1 == "email" {
			score += 0.5
		}

		if i < len(second) {
			t2 := InferType(second[i])
			if t1 != t2 && (t1 == "string" || t1 == "id_like" || t1 == "email") {
				score += 1.5
			} else if t1 == t2 && t1 != "string" && t1 != "id_like" && t1 != "email" {
				score -= 1.0
			}
		}
	}

	confidence := score / float64(len(first))
	return confidence >= 0.8
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
