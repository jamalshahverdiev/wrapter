package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
)

// RenderTemplate renders a template from a URL and writes it to a file
func RenderTemplate(url, templateName, outputDir string, data interface{}) error {
	resp, err := http.Get(url + templateName + ".j2")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Use io.ReadAll instead of ioutil.ReadAll
	templateData, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	tpl, err := template.New(templateName).Parse(string(templateData))
	if err != nil {
		return err
	}

	outputPath := filepath.Join(outputDir, templateName+".tf")
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	return tpl.Execute(outFile, data)
}
