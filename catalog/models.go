package catalog

import "time"

// CatalogEntry represents a single dataset in the catalog.
type CatalogEntry struct {
	Path             string          `json:"path"`
	Name             string          `json:"name"`
	Extension        string          `json:"extension"`
	SizeBytes        int64           `json:"size_bytes"`
	ModifiedAt       time.Time       `json:"modified_at"`
	Depth            int             `json:"depth"`
	Owner            string          `json:"owner"`
	Permissions      string          `json:"permissions"`
	Encoding         string          `json:"encoding"`
	Delimiter        string          `json:"delimiter"`
	HasHeader        bool            `json:"has_header"`
	RowCountEstimate int             `json:"row_count_estimate"`
	ColumnCount      int             `json:"column_count"`
	AverageRowWidth  int             `json:"average_row_width"`
	Checksum         string          `json:"checksum,omitempty"`
	HasACLs          bool            `json:"has_acls"`
	Columns          []ColumnSummary `json:"columns"`
	RulesTriggered   []RuleSummary   `json:"rules_triggered"`
	GovernanceFlags  GovernanceFlags `json:"governance"`
}

// ColumnSummary is a condensed column profile for the catalog.
type ColumnSummary struct {
	Name           string  `json:"name"`
	Position       int     `json:"position"`
	InferredType   string  `json:"inferred_type"`
	NullPercentage float64 `json:"null_percentage"`
	MaxLength      int     `json:"max_length"`
	HasSamples     bool    `json:"has_samples"`
}

// RuleSummary represents a triggered rule.
type RuleSummary struct {
	RuleName string `json:"rule_name"`
	Severity string `json:"severity"`
	Target   string `json:"target,omitempty"`
}

// GovernanceFlags indicates governance-related metadata.
type GovernanceFlags struct {
	HasChecksum   bool `json:"has_checksum"`
	HasACLs       bool `json:"has_acls"`
	SamplesHidden bool `json:"samples_hidden"`
}

// CatalogResult is the full catalog output.
type CatalogResult struct {
	GeneratedAt   time.Time      `json:"generated_at"`
	SourcePath    string         `json:"source_path"`
	SourceRun     string         `json:"source_run"`
	TotalDatasets int            `json:"total_datasets"`
	TotalColumns  int            `json:"total_columns"`
	TotalRules    int            `json:"total_rules"`
	Datasets      []CatalogEntry `json:"datasets"`
}
