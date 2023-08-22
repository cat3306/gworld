package protocol

import (
	"fmt"
	"google.golang.org/protobuf/proto"
)

type protocBufferCoder struct {
	CoderType CodeType
}

func (p *protocBufferCoder) Marshal(v interface{}) ([]byte, error) {
	vv, ok := v.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("type err:%T", v)
	}
	return proto.Marshal(vv)
}
func (p *protocBufferCoder) Unmarshal(bin []byte, v interface{}) error { //解码
	vv, ok := v.(proto.Message)
	if !ok {
		return fmt.Errorf("type err:%T", v)
	}
	return proto.Unmarshal(bin, vv)
}
func (p *protocBufferCoder) ToString() string {
	return "ProtoBuffer"
}
