package archiveanalyzer

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/bodgit/sevenzip"
	"github.com/nwaples/rardecode"

	"privantix-source-inspector/config"
	"privantix-source-inspector/csvanalyzer"
	"privantix-source-inspector/models"
	"privantix-source-inspector/parquetanalyzer"
	"privantix-source-inspector/xlsxanalyzer"
)

var dataExtensions = map[string]bool{
	".csv": true, ".txt": true, ".tsv": true, ".xlsx": true,
	".rdat": true, ".dat": true, ".parquet": true,
}

var archiveExtensions = map[string]bool{
	".zip": true, ".7z": true, ".rar": true,
}

func Analyze(discovered models.FileDiscovered, cfg config.Config) []models.FileProfile {
	ext := strings.ToLower(discovered.Ext)
	switch ext {
	case ".zip":
		return analyzeZip(discovered, cfg)
	case ".7z":
		return analyze7z(discovered, cfg)
	case ".rar":
		return analyzeRar(discovered, cfg)
	default:
		return []models.FileProfile{{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{"unsupported archive format"},
		}}
	}
}

type zipWorkItem struct {
	f       *zip.File
	ext     string
	virtual models.FileDiscovered
}

func analyzeZip(discovered models.FileDiscovered, cfg config.Config) []models.FileProfile {
	zr, err := zip.OpenReader(discovered.Path)
	if err != nil {
		return []models.FileProfile{{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	defer zr.Close()

	basePath := discovered.Path
	var workItems []zipWorkItem
	var nestedProfiles []models.FileProfile
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		if dataExtensions[ext] {
			workItems = append(workItems, zipWorkItem{
				f: f, ext: ext,
				virtual: models.FileDiscovered{
					Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Ext: ext,
					Size: int64(f.UncompressedSize64), Modified: discovered.Modified, Created: discovered.Created,
					Depth: discovered.Depth + 1, Owner: discovered.Owner, Permissions: discovered.Permissions,
				},
			})
		} else if cfg.RecursiveArchives && archiveExtensions[ext] {
			rc, err := f.Open()
			if err != nil {
				nestedProfiles = append(nestedProfiles, models.FileProfile{
					Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Extension: ext,
					SizeBytes: int64(f.UncompressedSize64), ModifiedAt: discovered.Modified,
					Depth: discovered.Depth + 1, Owner: discovered.Owner, Permissions: discovered.Permissions,
					Errors: []string{err.Error()},
				})
				continue
			}
			data, readErr := io.ReadAll(rc)
			rc.Close()
			if readErr != nil {
				continue
			}
			nested := analyzeZipFromBytes(data, basePath+"/"+f.Name, discovered, cfg)
			nestedProfiles = append(nestedProfiles, nested...)
		}
	}

	dataProfiles := processZipWorkItems(workItems, cfg)
	return append(nestedProfiles, dataProfiles...)
}

func analyzeZipFromBytes(data []byte, virtualPath string, parent models.FileDiscovered, cfg config.Config) []models.FileProfile {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return []models.FileProfile{{
			Path: virtualPath, Name: filepath.Base(virtualPath), Extension: filepath.Ext(virtualPath),
			SizeBytes: int64(len(data)), ModifiedAt: parent.Modified, CreatedAt: parent.Created,
			Depth: parent.Depth + 1, Owner: parent.Owner, Permissions: parent.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	var nestedProfiles []models.FileProfile
	for _, f := range zr.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		fullPath := virtualPath + "/" + f.Name
		if dataExtensions[ext] {
			rc, err := f.Open()
			if err != nil {
				nestedProfiles = append(nestedProfiles, models.FileProfile{
					Path: fullPath, Name: filepath.Base(f.Name), Extension: ext,
					SizeBytes: int64(f.UncompressedSize64), ModifiedAt: parent.Modified, CreatedAt: parent.Created,
					Depth: parent.Depth + 2, Owner: parent.Owner, Permissions: parent.Permissions,
					Errors: []string{err.Error()},
				})
				continue
			}
			fileData, readErr := io.ReadAll(rc)
			rc.Close()
			if readErr != nil {
				continue
			}
			virtual := models.FileDiscovered{
				Path: fullPath, Name: filepath.Base(f.Name), Ext: ext,
				Size: int64(f.UncompressedSize64), Modified: parent.Modified, Created: parent.Created,
				Depth: parent.Depth + 2, Owner: parent.Owner, Permissions: parent.Permissions,
			}
			p := analyzeEntry(fileData, ext, virtual, cfg)
			nestedProfiles = append(nestedProfiles, p)
		} else if cfg.RecursiveArchives && archiveExtensions[ext] {
			rc, err := f.Open()
			if err != nil {
				continue
			}
			nestedData, readErr := io.ReadAll(rc)
			rc.Close()
			if readErr != nil {
				continue
			}
			nested := analyzeZipFromBytes(nestedData, fullPath, parent, cfg)
			nestedProfiles = append(nestedProfiles, nested...)
		}
	}
	return nestedProfiles
}

func processZipWorkItems(items []zipWorkItem, cfg config.Config) []models.FileProfile {
	if len(items) == 0 {
		return nil
	}
	var profiles []models.FileProfile
	var mu sync.Mutex
	workers := cfg.Workers
	if workers < 1 {
		workers = 1
	}
	jobs := make(chan zipWorkItem, len(items))
	for _, it := range items {
		jobs <- it
	}
	close(jobs)

	var wg sync.WaitGroup
	for i := 0; i < workers && i < len(items); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for it := range jobs {
				rc, err := it.f.Open()
				if err != nil {
					mu.Lock()
					profiles = append(profiles, models.FileProfile{
						Path: it.virtual.Path, Name: it.virtual.Name, Extension: it.ext,
						SizeBytes: it.virtual.Size, ModifiedAt: it.virtual.Modified, CreatedAt: it.virtual.Created,
						Depth: it.virtual.Depth, Owner: it.virtual.Owner, Permissions: it.virtual.Permissions,
						Errors: []string{err.Error()},
					})
					mu.Unlock()
					continue
				}
				data, err := io.ReadAll(rc)
				rc.Close()
				if err != nil {
					mu.Lock()
					profiles = append(profiles, models.FileProfile{
						Path: it.virtual.Path, Name: it.virtual.Name, Extension: it.ext,
						SizeBytes: it.virtual.Size, ModifiedAt: it.virtual.Modified, CreatedAt: it.virtual.Created,
						Depth: it.virtual.Depth, Owner: it.virtual.Owner, Permissions: it.virtual.Permissions,
						Errors: []string{err.Error()},
					})
					mu.Unlock()
					continue
				}
				p := analyzeEntry(data, it.ext, it.virtual, cfg)
				mu.Lock()
				profiles = append(profiles, p)
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return profiles
}

func analyze7z(discovered models.FileDiscovered, cfg config.Config) []models.FileProfile {
	var profiles []models.FileProfile
	r, err := sevenzip.OpenReader(discovered.Path)
	if err != nil {
		return []models.FileProfile{{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	defer r.Close()

	basePath := discovered.Path
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		if !dataExtensions[ext] && !(cfg.RecursiveArchives && archiveExtensions[ext]) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			profiles = append(profiles, models.FileProfile{
				Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Extension: ext,
				SizeBytes: int64(f.UncompressedSize), ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
				Depth: discovered.Depth + 1, Owner: discovered.Owner, Permissions: discovered.Permissions,
				Errors: []string{err.Error()},
			})
			continue
		}
		data, err := io.ReadAll(rc)
		rc.Close()
		if err != nil {
			profiles = append(profiles, models.FileProfile{
				Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Extension: ext,
				SizeBytes: int64(f.UncompressedSize), ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
				Depth: discovered.Depth + 1, Owner: discovered.Owner, Permissions: discovered.Permissions,
				Errors: []string{err.Error()},
			})
			continue
		}
		virtual := models.FileDiscovered{
			Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Ext: ext,
			Size: int64(f.UncompressedSize), Modified: discovered.Modified, Created: discovered.Created,
			Depth: discovered.Depth + 1, Owner: discovered.Owner, Permissions: discovered.Permissions,
		}
		if cfg.RecursiveArchives && archiveExtensions[ext] {
			nested := analyzeArchiveFromBytes(data, ext, basePath+"/"+f.Name, discovered, cfg)
			profiles = append(profiles, nested...)
		} else {
			profiles = append(profiles, analyzeEntry(data, ext, virtual, cfg))
		}
	}
	return profiles
}

func analyzeArchiveFromBytes(data []byte, ext string, virtualPath string, parent models.FileDiscovered, cfg config.Config) []models.FileProfile {
	switch ext {
	case ".zip":
		return analyzeZipFromBytes(data, virtualPath, parent, cfg)
	case ".7z":
		return analyze7zFromBytes(data, virtualPath, parent, cfg)
	case ".rar":
		return analyzeRarFromBytes(data, virtualPath, parent, cfg)
	default:
		return nil
	}
}

func analyzeRarFromBytes(data []byte, virtualPath string, parent models.FileDiscovered, cfg config.Config) []models.FileProfile {
	rr, err := rardecode.NewReader(bytes.NewReader(data), "")
	if err != nil {
		return []models.FileProfile{{
			Path: virtualPath, Name: filepath.Base(virtualPath), Extension: ".rar",
			SizeBytes: int64(len(data)), ModifiedAt: parent.Modified, CreatedAt: parent.Created,
			Depth: parent.Depth + 1, Owner: parent.Owner, Permissions: parent.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	var result []models.FileProfile
	basePath := virtualPath
	for {
		hdr, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		if hdr.IsDir {
			continue
		}
		ext := strings.ToLower(filepath.Ext(hdr.Name))
		if !dataExtensions[ext] && !(cfg.RecursiveArchives && archiveExtensions[ext]) {
			continue
		}
		fileData, readErr := io.ReadAll(rr)
		if readErr != nil {
			continue
		}
		virtual := models.FileDiscovered{
			Path: basePath + "/" + hdr.Name, Name: filepath.Base(hdr.Name), Ext: ext,
			Size: hdr.UnPackedSize, Modified: parent.Modified, Created: parent.Created,
			Depth: parent.Depth + 2, Owner: parent.Owner, Permissions: parent.Permissions,
		}
		if cfg.RecursiveArchives && archiveExtensions[ext] {
			nested := analyzeArchiveFromBytes(fileData, ext, basePath+"/"+hdr.Name, parent, cfg)
			result = append(result, nested...)
		} else {
			result = append(result, analyzeEntry(fileData, ext, virtual, cfg))
		}
	}
	return result
}

func analyze7zFromBytes(data []byte, virtualPath string, parent models.FileDiscovered, cfg config.Config) []models.FileProfile {
	tmp, err := os.CreateTemp("", "privantix-7z-*.tmp")
	if err != nil {
		return []models.FileProfile{{
			Path: virtualPath, Name: filepath.Base(virtualPath), Extension: ".7z",
			SizeBytes: int64(len(data)), ModifiedAt: parent.Modified, CreatedAt: parent.Created,
			Depth: parent.Depth + 1, Owner: parent.Owner, Permissions: parent.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := tmp.Write(data); err != nil {
		tmp.Close()
		return nil
	}
	tmp.Close()
	r, err := sevenzip.OpenReader(tmpPath)
	if err != nil {
		return []models.FileProfile{{
			Path: virtualPath, Name: filepath.Base(virtualPath), Extension: ".7z",
			SizeBytes: int64(len(data)), ModifiedAt: parent.Modified, CreatedAt: parent.Created,
			Depth: parent.Depth + 1, Owner: parent.Owner, Permissions: parent.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	defer r.Close()
	var result []models.FileProfile
	basePath := virtualPath
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		ext := strings.ToLower(filepath.Ext(f.Name))
		if !dataExtensions[ext] && !(cfg.RecursiveArchives && archiveExtensions[ext]) {
			continue
		}
		rc, err := f.Open()
		if err != nil {
			result = append(result, models.FileProfile{
				Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Extension: ext,
				SizeBytes: int64(f.UncompressedSize), ModifiedAt: parent.Modified, CreatedAt: parent.Created,
				Depth: parent.Depth + 2, Owner: parent.Owner, Permissions: parent.Permissions,
				Errors: []string{err.Error()},
			})
			continue
		}
		fileData, readErr := io.ReadAll(rc)
		rc.Close()
		if readErr != nil {
			continue
		}
		virtual := models.FileDiscovered{
			Path: basePath + "/" + f.Name, Name: filepath.Base(f.Name), Ext: ext,
			Size: int64(f.UncompressedSize), Modified: parent.Modified, Created: parent.Created,
			Depth: parent.Depth + 2, Owner: parent.Owner, Permissions: parent.Permissions,
		}
		if cfg.RecursiveArchives && archiveExtensions[ext] {
			nested := analyzeArchiveFromBytes(fileData, ext, basePath+"/"+f.Name, parent, cfg)
			result = append(result, nested...)
		} else {
			result = append(result, analyzeEntry(fileData, ext, virtual, cfg))
		}
	}
	return result
}

func analyzeRar(discovered models.FileDiscovered, cfg config.Config) []models.FileProfile {
	var profiles []models.FileProfile
	f, err := os.Open(discovered.Path)
	if err != nil {
		return []models.FileProfile{{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}}
	}
	defer f.Close()

	rr, err := rardecode.NewReader(f, "")
	if err != nil {
		return []models.FileProfile{{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}}
	}

	basePath := discovered.Path
	for {
		hdr, err := rr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			profiles = append(profiles, models.FileProfile{
				Path: basePath, Name: discovered.Name, Extension: discovered.Ext,
				SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
				Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
				Errors: []string{err.Error()},
			})
			break
		}
		if hdr.IsDir {
			continue
		}
		ext := strings.ToLower(filepath.Ext(hdr.Name))
		if !dataExtensions[ext] && !(cfg.RecursiveArchives && archiveExtensions[ext]) {
			continue
		}

		data, err := io.ReadAll(rr)
		if err != nil {
			profiles = append(profiles, models.FileProfile{
				Path: basePath + "/" + hdr.Name, Name: filepath.Base(hdr.Name), Extension: ext,
				SizeBytes: hdr.UnPackedSize, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
				Depth: discovered.Depth + 1, Owner: discovered.Owner, Permissions: discovered.Permissions,
				Errors: []string{err.Error()},
			})
			continue
		}

		virtual := models.FileDiscovered{
			Path:        basePath + "/" + hdr.Name,
			Name:        filepath.Base(hdr.Name),
			Ext:         ext,
			Size:        hdr.UnPackedSize,
			Modified:    discovered.Modified,
			Created:     discovered.Created,
			Depth:       discovered.Depth + 1,
			Owner:       discovered.Owner,
			Permissions: discovered.Permissions,
		}

		if cfg.RecursiveArchives && archiveExtensions[ext] {
			nested := analyzeArchiveFromBytes(data, ext, basePath+"/"+hdr.Name, discovered, cfg)
			profiles = append(profiles, nested...)
		} else {
			p := analyzeEntry(data, ext, virtual, cfg)
			profiles = append(profiles, p)
		}
	}
	return profiles
}

func analyzeEntry(data []byte, ext string, virtual models.FileDiscovered, cfg config.Config) models.FileProfile {
	maxRows := cfg.MaxSampleRows
	if maxRows <= 0 {
		maxRows = 200
	}

	switch ext {
	case ".csv", ".txt", ".tsv", ".rdat", ".dat":
		return csvanalyzer.AnalyzeFromReader(bytes.NewReader(data), virtual, maxRows)
	case ".xlsx":
		return xlsxAnalyzeFromBytes(data, virtual, maxRows)
	case ".parquet":
		return parquetAnalyzeFromBytes(data, virtual, maxRows)
	default:
		return models.FileProfile{
			Path: virtual.Path, Name: virtual.Name, Extension: virtual.Ext,
			SizeBytes: virtual.Size, ModifiedAt: virtual.Modified, CreatedAt: virtual.Created,
			Depth: virtual.Depth, Owner: virtual.Owner, Permissions: virtual.Permissions,
			Errors: []string{"unsupported format inside archive"},
		}
	}
}

func xlsxAnalyzeFromBytes(data []byte, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	zr, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return models.FileProfile{
			Path: discovered.Path, Name: discovered.Name, Extension: discovered.Ext,
			SizeBytes: discovered.Size, ModifiedAt: discovered.Modified, CreatedAt: discovered.Created,
			Depth: discovered.Depth, Owner: discovered.Owner, Permissions: discovered.Permissions,
			Errors: []string{err.Error()},
		}
	}
	return xlsxanalyzer.AnalyzeFromZipReader(zr, discovered, maxSampleRows)
}

func parquetAnalyzeFromBytes(data []byte, discovered models.FileDiscovered, maxSampleRows int) models.FileProfile {
	return parquetanalyzer.AnalyzeFromReaderAt(bytes.NewReader(data), int64(len(data)), discovered, maxSampleRows)
}
