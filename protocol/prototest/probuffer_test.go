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
	if err != nil {
		t.Fatalf("TestGameCoder err:%s", err.Error())
	}
	t.Log(len(jBin))
	bin, err := coder.Marshal(s)
	if err != nil {
		t.Fatalf("TestGameCoder err:%s", err.Error())
	}
	t.Log(len(bin))
}

func Test1(t *testing.T) {

	sss := &ServerMsg{
		Logic:    1,
		Payload:  []byte(`[:]`),
		Method:   2,
		CodeType: 4,
	}
	c:=protocol.GameCoder(protocol.ProtoBuffer)
	raw, err := c.Marshal(sss)
	if err != nil {
		t.Logf(err.Error())
	}
	t.Logf(string(raw))
	ssss := &ServerMsg{}
	err = protocol.GameCoder(protocol.ProtoBuffer).Unmarshal(raw, ssss)
	t.Log(ssss)
}