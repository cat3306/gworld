package engine

import "github.com/cat3306/goworld/protocol"

type ServerInnerMsg struct {
	ClientId       string            `json:"client_id"`
	Payload        []byte            `json:"payload"`
	ClientMethod   uint32            `json:"client_method"`
	ClientCodeType protocol.CodeType `json:"client_code_type"`
}
type ClientMsg struct {
	Logic    uint32            `json:"logic"`
	Payload  []byte            `json:"payload"`
	Method   uint32            `json:"method"`
	CodeType protocol.CodeType `json:"client_code_type"`
}
