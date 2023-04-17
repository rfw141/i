package protox

import (
	"github.com/emicklei/proto"
	"github.com/rfw141/i/log"
	"os"
)

func GetDefinition(protoPath string) (*proto.Proto, error) {
	f, err := os.Open(protoPath)
	if err != nil {
		log.Errorf("open file fail: %v", err)
		return nil, err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			log.Errorf("close file fail: %v", err)
		}
	}(f)

	definition, err := proto.NewParser(f).Parse()
	if err != nil {
		log.Errorf("parse proto fail: %v", err)
		return nil, err
	}
	return definition, nil
}
