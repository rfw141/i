package tplx

import (
	"bytes"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
	"github.com/rfw141/i/x/filex"
	"os"
	"path/filepath"
	"text/template"
)

func ParseAndSave(tplFile, outFile string, data interface{}) error {
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}

	var file *os.File
	if !filex.Exist(outFile) {
		if !filex.Exist(filepath.Dir(outFile)) {
			if err = os.MkdirAll(filepath.Dir(outFile), 0755); err != nil {
				return err
			}
		}
		if file, err = os.Create(outFile); err != nil {
			return err
		}
	} else {
		if file, err = os.OpenFile(outFile, os.O_RDWR|os.O_TRUNC, 0755); err != nil {
			return err
		}
	}

	defer func(file *os.File) {
		if err = file.Close(); err != nil {
			log.Errorf("close file error: %v", err)
		}
	}(file)
	if err = tpl.Execute(file, data); err != nil {
		return err
	}
	return nil
}

func MustParseAndSave(tplFile, outFile string, data interface{}) {
	err := ParseAndSave(tplFile, outFile, data)
	errorx.Must(err)
}

func ParseAndGet(tplFile string, data interface{}) (string, error) {
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return "", errorx.New("parse template file(%s) error: %v", tplFile, err)
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, data); err != nil {
		return "", errorx.New("execute template file(%s) error: %v", tplFile, err)
	}
	return content.String(), nil
}

func MustParseAndGet(tplFile string, data interface{}) string {
	content, err := ParseAndGet(tplFile, data)
	errorx.Must(err)
	return content
}

func ParseStrAndGet(tplStr string, data interface{}) (string, error) {
	tpl, err := template.New("tpl").Parse(tplStr)
	if err != nil {
		return "", err
	}
	var content bytes.Buffer
	if err = tpl.Execute(&content, data); err != nil {
		return "", err
	}
	return content.String(), nil
}

func MustParseStrAndGet(tplStr string, data interface{}) string {
	content, err := ParseStrAndGet(tplStr, data)
	errorx.Must(err)
	return content
}
