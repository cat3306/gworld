package protocol

import (
	"errors"
	"fmt"
)

type rawUint32 struct {
	CodeType CodeType
}

func (r *rawUint32) ToString() string {
	return "uint32"
}
func (r *rawUint32) Unmarshal(b []byte, v interface{}) error {
	if _, ok := v.(*uint32); ok {
		u := packetEndian.Uint32(b)
		v = &u
		return nil
	}
	return errors.New("v type not *int")
}
func (r *rawUint32) Marshal(v interface{}) ([]byte, error) {
	if vv, ok := v.(uint32); ok {
		raw := make([]byte, 4)
		packetEndian.PutUint32(raw, vv)
		return raw, nil
	}
	return nil, fmt.Errorf("v type not string")
}
