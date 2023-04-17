package tplx

import (
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/filex"
	"os"
	"text/template"
)

func ParseAndSave(tplFile, outFile string, data interface{}) error {
	tpl, err := template.ParseFiles(tplFile)
	if err != nil {
		return err
	}

	var file *os.File
	if !filex.Exist(outFile) {
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
