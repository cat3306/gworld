package prototest

import (
	"github.com/cat3306/goworld/protocol"
	"testing"
)

func TestGameCoder(t *testing.T) {
	coder := protocol.GameCoder(protocol.ProtoBuffer)
	s := &Student{
		Name:   "cat13",
		Male:   true,
		Scores: []int32{12, 21},
	}
	j := protocol.GameCoder(protocol.Json)
	jBin, err := j.Marshal(s)
	if err!=nil{
		t.Fatalf("TestGameCoder err:%s", err.Error())
	}
	t.Log(len(jBin))
	bin, err := coder.Marshal(s)
	if err != nil {
		t.Fatalf("TestGameCoder err:%s", err.Error())
	}
	t.Log(len(bin))
	ss := Student{}
	err = coder.Unmarshal(bin, &ss)
	if err != nil {
		t.Fatalf("TestGameCoder err:%s", err.Error())
	}
	t.Log(ss)
}
