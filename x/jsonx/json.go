package jsonx

import (
	"encoding/json"
	"github.com/rfw141/i/log"
	"github.com/rfw141/i/x/errorx"
)

func Unmarshal(data []byte, v interface{}) error {
	err := json.Unmarshal(data, v)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}
	return nil
}

func Marshal(v interface{}) ([]byte, error) {
	data, err := json.Marshal(v)
	if err != nil {
		log.Errorf("err:%v", err)
		return nil, err
	}
	return data, nil
}

func MustUnmarshal(data []byte, v interface{}) {
	errorx.Must(Unmarshal(data, v))
}

func MustMarshal(v interface{}) []byte {
	data, err := Marshal(v)
	errorx.Must(err)
	return data
}
