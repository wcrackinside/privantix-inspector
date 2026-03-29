package catalog

import (
	"encoding/json"
	"os"

	"privantix-source-inspector/models"
)

// ParseAnalysis reads analysis.json and returns the inspector result.
func ParseAnalysis(path string) (*models.AnalysisResult, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var result models.AnalysisResult
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}
