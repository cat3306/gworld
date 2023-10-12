package protocol

type CodeType uint16

const (
	CodeNone    = CodeType(0)
	String      = CodeType(1)
	Json        = CodeType(2)
	ProtoBuffer = CodeType(3)
)

var (
	coderSet = map[CodeType]Coder{
		Json:        &jsonCoder{},
		String:      &rawString{},
		ProtoBuffer: &protocBufferCoder{},
	}
)

type Coder interface {
	Unmarshal([]byte, interface{}) error   //解码
	Marshal(v interface{}) ([]byte, error) //编码
	ToString() string
}

func GameCoder(codeType CodeType) Coder {
	coder := coderSet[codeType]
	if coder == nil {
		coder = coderSet[Json]
	}
	return coder
}
