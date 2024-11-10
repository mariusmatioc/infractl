package global

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// RemoveQuotes removes quotes from a string if present and returns the string and true if quotes were removed
func RemoveQuotes(s string) (string, bool) {
	s = strings.TrimSpace(s)
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1], true
	}
	return s, false
}

func TrimAndRemoveQuotes(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > 1 && s[0] == '"' && s[len(s)-1] == '"' {
		return s[1 : len(s)-1]
	}
	return s
}

func SetRootFolder(args []string) {
	var root string
	if len(args) > 0 {
		root = args[0]
		root = os.ExpandEnv(root)
		root, _ = filepath.Abs(root)
	} else {
		root, _ = os.Getwd()
	}
	RootParent = root
	RootFolder = filepath.Join(root, InfraCtlRoot)
}

func SetDefaultRootFolder() {
	root, _ := os.Getwd()
	RootParent = root
	RootFolder = filepath.Join(root, InfraCtlRoot)
}

func AddQuotes(s string) string {
	return `"` + s + `"`
}

func QuotedArray(arr []string) string {
	items := []string{}
	for i := range arr {
		items = append(items, AddQuotes(arr[i]))
	}
	return "[" + strings.Join(items, ",") + "]"
}

func AdjustAwsString(s string) string {
	s = strings.ReplaceAll(s, "_", "-")
	s = strings.TrimSpace(s)
	return s
}

// GetCraftPath returns the craft absolute path from the command line arguments and sets the root folder
func GetCraftPath(args []string) (path string, err error) {
	craftName := ""
	if len(args) == 2 {
		craftName = args[1]
		SetRootFolder(args)
	} else {
		craftName = args[0]
		SetRootFolder([]string{})
	}
	ext := filepath.Ext(craftName)
	if ext == "" {
		craftName += ".yml"
	}
	path, err = GetAbsoluteCraftPath(craftName)
	return
}

func GetAbsoluteCraftPath(path string) (string, error) {
	if !filepath.IsAbs(path) {
		path = filepath.Join(RootFolder, CratsSubFolder, path)
	}
	if !FileExists(path) {
		err := fmt.Errorf("craft file %s does not exist", path)
		return "", err
	}
	return path, nil
}

//func GetCraftFilePathFromName(craftName string) string {
//	ext := filepath.Ext(craftName)
//	if ext == "" {
//		craftName += ".yml"
//	}
//	return filepath.Join(RootFolder, CratsSubFolder, craftName)
//}

//func ToAbsPath(files *[]string) {
//	for i := range *files {
//		(*files)[i], _ = filepath.Abs((*files)[i])
//	}
//}

func ToAbsPathBasedOn(base, path string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(base, path)
}

func RemoveFileExtension(fileName string) string {
	return fileName[:len(fileName)-len(filepath.Ext(fileName))]
}

func NameOnly(path string) string {
	return RemoveFileExtension(filepath.Base(path))
}

func GetBuildFolder(craftName string) (buildFolder string, err error) {
	buildFolder = filepath.Join(RootFolder, BuildSubFolder, craftName)
	err = os.MkdirAll(buildFolder, FilePerm)
	return
}

// FileExists also works for folders
func FileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func DeleteFiles(folder, wildcard string) (err error) {
	files, err := filepath.Glob(filepath.Join(folder, wildcard))
	if err != nil {
		return
	}
	for _, f := range files {
		if err = os.Remove(f); err != nil {
			return
		}
	}
	return
}

func WriteStringToFile(path, content string) (err error) {
	f, err := os.Create(path)
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	_, err = f.WriteString(content)
	return
}
