package filex

import (
	"github.com/rfw141/i/log"
	"os"
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

func Write(filename, content string) error {
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
