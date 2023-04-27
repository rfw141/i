package filex

import (
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
	"os"
	"path/filepath"
)

func Exist(filename string) bool {
	if _, err := os.Stat(filename); err != nil {
		if os.IsNotExist(err) {
			return false
		}
		log.Errorf("err:%v", err)
	}
	return true
}

func MustExist(filename string) error {
	if _, err := os.Stat(filename); err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func Read(filename string) ([]byte, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return content, nil
}

func MustRead(filename string) []byte {
	content, err := Read(filename)
	errorx.Must(err)
	return content
}

func MustReadStr(filename string) string {
	content, err := Read(filename)
	errorx.Must(err)
	return string(content)
}

func Write(filename, content string) error {
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		return err
	}
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("close file fail: %v", err)
		}
	}(f)
	if _, err = f.WriteString(content); err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func MustWrite(filename, content string) {
	err := Write(filename, content)
	errorx.Must(err)
}

func Path(paths ...string) string {
	return filepath.ToSlash(filepath.Join(paths...))
}
