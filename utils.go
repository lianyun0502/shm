package shm

import (
	"bytes"
	"encoding/gob"
)

func GetStruct[T any](binData []byte) (*T, error) {
	obj := new(T)
	reader := bytes.NewReader(binData)
	dec := gob.NewDecoder(reader)
	if err := dec.Decode(obj); err != nil {
		return nil, err
	}
	return obj, nil
}

func GetBinary(obj any) ([]byte, error) {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	if err := enc.Encode(obj); err != nil {
		return nil, err
	}
	return b.Bytes(), nil
}
