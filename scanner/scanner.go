package scanner

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"privantix-source-inspector/models"
	"privantix-source-inspector/utils"
)

func Scan(root string, extensions []string, recursive bool) ([]models.FileDiscovered, []string, error) {
	var files []models.FileDiscovered
	var errs []string
	extSet := map[string]struct{}{}
	for _, ext := range utils.NormalizeExtensions(extensions) {
		extSet[ext] = struct{}{}
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			errs = append(errs, err.Error())
			return nil
		}
		if d.IsDir() {
			if !recursive && path != root {
				return filepath.SkipDir
			}
			return nil
		}
		info, statErr := d.Info()
		if statErr != nil {
			errs = append(errs, statErr.Error())
			return nil
		}
		ext := strings.ToLower(filepath.Ext(path))
		if _, ok := extSet[ext]; !ok {
			return nil
		}
		owner, perms := utils.GetFileOwnerAndPerms(info, path)
		created := info.ModTime()
		if t, err := utils.GetFileCreationTime(path); err == nil {
			created = t
		}

		files = append(files, models.FileDiscovered{
			Path:        path,
			Name:        info.Name(),
			Ext:         ext,
			Size:        info.Size(),
			Modified:    info.ModTime(),
			Created:     created,
			Depth:       utils.FileDepth(root, path),
			Owner:       owner,
			Permissions: perms,
		})
		return nil
	})
	if os.IsNotExist(err) {
		return nil, errs, err
	}
	return files, errs, err
}
