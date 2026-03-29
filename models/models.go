package models

import "time"

type RunMetadata struct {
	StartedAt           time.Time `json:"started_at"`
	CompletedAt         time.Time `json:"completed_at"`
	DurationSeconds     float64   `json:"duration_seconds"`
	Path                string    `json:"path"`
	OutputDir           string    `json:"output_dir"`
	SupportedFiles      int       `json:"supported_files"`
	AnalyzedFiles       int       `json:"analyzed_files"`
	FailedFiles         int       `json:"failed_files"`
	MaxSampleRows       int       `json:"max_sample_rows"`
	Workers             int       `json:"workers"`
	SupportedExtensions []string  `json:"supported_extensions"`
}

type FileDiscovered struct {
	Path        string
	Name        string
	Ext         string
	Size        int64
	Modified    time.Time
	Created     time.Time
	Depth       int
	Owner       string
	Permissions string
	ACLs        []string
}

type RuleResult struct {
	RuleName string `json:"rule_name"`
	Severity string `json:"severity"`
	Message  string `json:"message"`
	Target   string `json:"target,omitempty"`
}

type ColumnProfile struct {
	FilePath       string   `json:"file_path"`
	SheetName      string   `json:"sheet_name,omitempty"`
	Name           string   `json:"name"`
	Position       int      `json:"position"`
	InferredType   string   `json:"inferred_type"`
	NullPercentage float64  `json:"null_percentage"`
	MaxLength      int      `json:"max_length"`
	SampleValues   []string `json:"sample_values"`
}

type FileProfile struct {
	Path             string          `json:"path"`
	Name             string          `json:"name"`
	Extension        string          `json:"extension"`
	SizeBytes        int64           `json:"size_bytes"`
	ModifiedAt       time.Time       `json:"modified_at"`
	CreatedAt        time.Time       `json:"created_at,omitempty"`
	Depth            int             `json:"depth"`
	Owner            string          `json:"owner"`
	Permissions      string          `json:"permissions"`
	ACLs                 []string        `json:"acls,omitempty"`
	TrustedPrincipals    []string        `json:"trusted_principals,omitempty"`
	OutsidePrincipals    []string        `json:"outside_principals,omitempty"`
	ComplianceStatus     string          `json:"compliance_status,omitempty"`
	Checksum             string          `json:"checksum,omitempty"`
	Encoding         string          `json:"encoding"`
	Delimiter        string          `json:"delimiter,omitempty"`
	HasHeader        bool            `json:"has_header"`
	RowCountEstimate int             `json:"row_count_estimate"`
	ColumnCount      int             `json:"column_count"`
	AverageRowWidth  int             `json:"average_row_width"`
	SheetName        string          `json:"sheet_name,omitempty"`
	Columns          []ColumnProfile `json:"columns"`
	RulesTriggered   []RuleResult    `json:"rules_triggered"`
	Errors           []string        `json:"errors,omitempty"`
}

type AnalysisResult struct {
	Run    RunMetadata   `json:"run"`
	Files  []FileProfile `json:"files"`
	Errors []string      `json:"errors"`
}
