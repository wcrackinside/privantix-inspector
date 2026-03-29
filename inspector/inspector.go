package inspector

import (
	"strings"

	"privantix-source-inspector/analyzer"
	"privantix-source-inspector/archiveanalyzer"
	"privantix-source-inspector/csvanalyzer"
	"privantix-source-inspector/config"
	"privantix-source-inspector/models"
	"privantix-source-inspector/parquetanalyzer"
	"privantix-source-inspector/utils"
	"privantix-source-inspector/xlsxanalyzer"
)

func Analyze(discovered models.FileDiscovered, cfg config.Config) []models.FileProfile {
	var results []models.FileProfile
	var profiles []models.FileProfile

	switch strings.ToLower(discovered.Ext) {
	case ".csv", ".txt", ".tsv", ".rdat", ".dat":
		p := csvanalyzer.Analyze(discovered, cfg.MaxSampleRows)
		profiles = []models.FileProfile{p}
	case ".xlsx":
		p := xlsxanalyzer.Analyze(discovered, cfg.MaxSampleRows)
		profiles = []models.FileProfile{p}
	case ".parquet":
		p := parquetanalyzer.Analyze(discovered, cfg.MaxSampleRows)
		profiles = []models.FileProfile{p}
	case ".7z", ".zip", ".rar":
		profiles = archiveanalyzer.Analyze(discovered, cfg)
	default:
		profiles = []models.FileProfile{{
			Path:       discovered.Path,
			Name:       discovered.Name,
			Extension:  discovered.Ext,
			SizeBytes:  discovered.Size,
			ModifiedAt: discovered.Modified,
			CreatedAt:  discovered.Created,
			Depth:      discovered.Depth,
			Errors:     []string{"unsupported file type"},
		}}
	}

	for i := range profiles {
		// Checksum and ACLs only for direct files (not entries inside archives)
		if len(profiles) == 1 && profiles[i].Path == discovered.Path {
			if cfg.Checksum {
				profiles[i].Checksum = utils.CalculateChecksum(discovered.Path)
			}
			if cfg.Security {
				profiles[i].ACLs = utils.GetFileACLs(discovered.Path)
			}
		}
		analyzer.ApplyRules(&profiles[i])
		results = append(results, profiles[i])
	}
	return results
}
