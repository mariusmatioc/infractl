package global

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

func (cfg *Config) BuildSliceFromTemplate(templateString, fileName string, slice []any) (err error) {
	f, err := os.Create(filepath.Join(cfg.BuildFolder, fileName))
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	content := ""
	for _, obj := range slice {
		s, err := StringFromTemplate(templateString, obj)
		if err != nil {
			return err
		}
		content += "\n"
		content += s
	}
	_, err = fmt.Fprint(f, content)
	if err == nil {
		fmt.Println("Wrote", fileName)
	}
	return
}

func (cfg *Config) BuildFromServiceTemplate(templateString, fileName string, services Services) (err error) {
	slice := []any{}
	for _, serv := range services {
		slice = append(slice, serv)
	}
	return cfg.BuildSliceFromTemplate(templateString, fileName, slice)
}

func (cfg *Config) BuildFromRecipeTemplate(templateString, fileName string) (err error) {
	f, err := os.Create(filepath.Join(cfg.BuildFolder, fileName))
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	content := ""
	s, err := StringFromTemplate(templateString, cfg.Recipe)
	if err != nil {
		return err
	}
	content += "\n"
	content += s
	_, err = fmt.Fprint(f, content)
	if err == nil {
		fmt.Println("Wrote", fileName)
	}
	return
}

func WriteFargateFile(buildFolder, file, content string) (err error) {
	f, err := os.Create(filepath.Join(buildFolder, file))
	if err != nil {
		return
	}
	defer func() { _ = f.Close() }()

	_, err = fmt.Fprint(f, content)
	if err == nil {
		fmt.Println("Wrote", file)
	}
	return
}

func (cfg *Config) writeFargateFile(file, content string) (err error) {
	return WriteFargateFile(cfg.BuildFolder, file, content)
}

func StringFromTemplate(templateString string, config any) (s string, err error) {
	parsed, err := template.New("template").Parse(templateString)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = parsed.Execute(&buf, config)
	if err != nil {
		return
	}
	s = buf.String()
	return
}

func BuildFromTemplate(templateString, buildFolder, fileName string, config any) (err error) {
	parsed, err := template.New("template").Parse(templateString)
	if err != nil {
		return
	}
	var buf bytes.Buffer
	err = parsed.Execute(&buf, config)
	if err != nil {
		return
	}
	err = WriteFargateFile(buildFolder, fileName, buf.String())
	return
}

func (cfg *Config) BuildFromTemplate(templateString, fileName string, config any) (err error) {
	return BuildFromTemplate(templateString, cfg.BuildFolder, fileName, config)
}
