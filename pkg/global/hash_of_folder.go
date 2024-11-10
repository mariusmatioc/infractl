package global

import (
	"crypto/md5"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// HashOfFolder computes a hash of files in the folder, 16 hex characters
// This is used to determine if the folder has changed
func HashOfFolder(folder string) (string, error) {
	hashes := []byte{}
	err := filepath.WalkDir(folder, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if strings.HasPrefix(d.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		// We don't want to skip files like .env
		//if strings.HasPrefix(d.Name(), ".") {
		//	return nil
		//}
		// Skip very large files. Most likely not source code
		info, _ := d.Info()
		if info.Size() > 10*MiB {
			return nil
		}
		data, err2 := os.ReadFile(path) // This reads whole file. Perhaps method for stream hashing?
		if err2 != nil {
			// Check if this file is a symlink. If it is, skip it instead of returning an error
			var stat, _ = os.Lstat(path)
			if (stat.Mode() & fs.ModeSymlink) != 0 {
				return nil
			}

			// Not a symlink, so we think it is a file error
			return err2
		}
		hash := md5.Sum(data)
		hashes = append(hashes, hash[:]...)
		return nil
	})

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", md5.Sum(hashes)), nil
}
