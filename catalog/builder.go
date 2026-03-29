package catalog

import (
	"time"

	"privantix-source-inspector/models"
)

// Build creates a CatalogResult from an AnalysisResult.
func Build(result *models.AnalysisResult) *CatalogResult {
	catalog := &CatalogResult{
		GeneratedAt:   time.Now(),
		SourcePath:    result.Run.Path,
		SourceRun:     result.Run.StartedAt.Format("2006-01-02 15:04:05"),
		TotalDatasets: len(result.Files),
		Datasets:      make([]CatalogEntry, 0, len(result.Files)),
	}

	for _, f := range result.Files {
		entry := CatalogEntry{
			Path:             f.Path,
			Name:             f.Name,
			Extension:        f.Extension,
			SizeBytes:        f.SizeBytes,
			ModifiedAt:       f.ModifiedAt,
			Depth:            f.Depth,
			Owner:            f.Owner,
			Permissions:      f.Permissions,
			Encoding:         f.Encoding,
			Delimiter:        f.Delimiter,
			HasHeader:        f.HasHeader,
			RowCountEstimate: f.RowCountEstimate,
			ColumnCount:      f.ColumnCount,
			AverageRowWidth:  f.AverageRowWidth,
			Checksum:         f.Checksum,
			HasACLs:          len(f.ACLs) > 0,
			GovernanceFlags: GovernanceFlags{
				HasChecksum:   f.Checksum != "",
				HasACLs:       len(f.ACLs) > 0,
				SamplesHidden: samplesHidden(f.Columns),
			},
		}

		for _, c := range f.Columns {
			entry.Columns = append(entry.Columns, ColumnSummary{
				Name:           c.Name,
				Position:       c.Position,
				InferredType:   c.InferredType,
				NullPercentage: c.NullPercentage,
				MaxLength:      c.MaxLength,
				HasSamples:     len(c.SampleValues) > 0,
			})
			catalog.TotalColumns++
		}

		for _, r := range f.RulesTriggered {
			entry.RulesTriggered = append(entry.RulesTriggered, RuleSummary{
				RuleName: r.RuleName,
				Severity: r.Severity,
				Target:   r.Target,
			})
			catalog.TotalRules++
		}

		catalog.Datasets = append(catalog.Datasets, entry)
	}

	return catalog
}

func samplesHidden(cols []models.ColumnProfile) bool {
	for _, c := range cols {
		if len(c.SampleValues) > 0 {
			return false
		}
	}
	return len(cols) > 0
}
