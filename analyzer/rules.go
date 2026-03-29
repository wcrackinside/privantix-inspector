package analyzer

import "privantix-source-inspector/models"

func ApplyRules(file *models.FileProfile) {
	var results []models.RuleResult
	if file.Encoding != "utf-8" && file.Encoding != "utf-8-bom" && file.Encoding != "xlsx-internal" && file.Encoding != "" {
		results = append(results, models.RuleResult{
			RuleName: "encoding_not_utf8",
			Severity: "warning",
			Message:  "File encoding is not UTF-8",
		})
	}
	if !file.HasHeader {
		results = append(results, models.RuleResult{
			RuleName: "missing_header",
			Severity: "warning",
			Message:  "Header was not confidently detected",
		})
	}
	if file.ColumnCount > 100 {
		results = append(results, models.RuleResult{
			RuleName: "too_many_columns",
			Severity: "warning",
			Message:  "File has more than 100 columns",
		})
	}
	for _, c := range file.Columns {
		if c.NullPercentage > 50 {
			results = append(results, models.RuleResult{
				RuleName: "high_null_ratio",
				Severity: "warning",
				Message:  "Column has null ratio above 50%",
				Target:   c.Name,
			})
		}
		if c.InferredType == "email" || c.InferredType == "phone" || c.InferredType == "id_like" {
			results = append(results, models.RuleResult{
				RuleName: "potential_sensitive_data",
				Severity: "info",
				Message:  "Column may contain personal or sensitive identifiers",
				Target:   c.Name,
			})
		}
	}
	file.RulesTriggered = results
}
