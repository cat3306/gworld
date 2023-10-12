package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/cat3306/goworld/glog"
	"github.com/panjf2000/gnet/v2"
	"github.com/valyala/bytebufferpool"
)

//
// * 0                       4                       8           10			 14
// * +-----------------------+-----------------------+-----------+-----------+
// * |   body len    		 |       protocol        | code type |  logic    |
// * +-----------------------+-----------------------+-----------+-----------+
// * |                                   			 						 |
// * +                                       		            			 +
// * |                   body bytes              		       	     		 |
// * +                                   						             +
// * |                                  						             |
// * +-----------------------------------------------------------+-----------+

const (
	payloadLen   = uint32(4)
	protocolLen  = uint32(4)
	codeTypeLen  = uint32(2)
	logicLen     = uint32(4)
	maxBufferCap = 1 << 24 //16M
)

var (
	packetEndian        = binary.LittleEndian
	ErrIncompletePacket = errors.New("incomplete packet")
	ErrTooLargePacket   = errors.New("too large packet")
	ErrDiscardedPacket  = errors.New("discarded not equal msg len")
)

func Decode(c gnet.Conn) (*Context, error) {

	bodyOffset := int(payloadLen + protocolLen + codeTypeLen + logicLen)
	headerBuffer, err := c.Peek(bodyOffset)
	if err != nil {
		return nil, err
	}
	if len(headerBuffer) < bodyOffset {
		return nil, ErrIncompletePacket
	}
	bodyLen := packetEndian.Uint32(headerBuffer[:payloadLen])
	protocol := packetEndian.Uint32(headerBuffer[payloadLen : payloadLen+protocolLen])
	codeType := packetEndian.Uint16(headerBuffer[payloadLen+protocolLen : payloadLen+protocolLen+codeTypeLen])
	logic := packetEndian.Uint32(headerBuffer[payloadLen+protocolLen+codeTypeLen : bodyOffset])
	msgLen := bodyOffset + int(bodyLen)
	if msgLen > maxBufferCap {
		c.Close()
		return nil, ErrTooLargePacket
	}

	if c.InboundBuffered() < msgLen {
		return nil, ErrIncompletePacket
	}
	msgBuffer, err := c.Peek(msgLen)
	if err != nil {
		return nil, err
	}
	discarded, err := c.Discard(msgLen)
	if err != nil {
		return nil, err
	}
	if discarded != msgLen {
		glog.Logger.Sugar().Errorf("discarded")
		return nil, ErrDiscardedPacket
	}
	buffer := bytebufferpool.Get()
	_, _ = buffer.Write(msgBuffer[bodyOffset:])
	ctx := &Context{
		Payload:  buffer,
		CodeType: CodeType(codeType),
		Proto:    protocol,
		Conn:     c,
		Logic:    logic,
	}
	return ctx, nil
}
func Encode(v interface{}, codeType CodeType, proto uint32, logic uint32) *bytebufferpool.ByteBuffer {
	if v == nil {
		panic("v nil")
	}
	var (
		body []byte
		err  error
	)
	if tmp, ok := v.([]byte); ok {
		body = tmp
	} else {
		body, err = GameCoder(codeType).Marshal(v)
		if err != nil {
			panic(err)
		}
	}
	//bodyOffset := int(payloadLen + protocolLen + codeTypeLen)
	//msgLen := bodyOffset + len(raw)
	//buffer := *BUFFERPOOL.Get(uint32(msgLen))
	buffer := bytebufferpool.Get()
	headBuffer := make([]byte, payloadLen+protocolLen+codeTypeLen+logicLen)
	packetEndian.PutUint32(headBuffer, uint32(len(body)))
	packetEndian.PutUint32(headBuffer[payloadLen:], proto)
	packetEndian.PutUint16(headBuffer[payloadLen+protocolLen:], uint16(codeType))
	packetEndian.PutUint32(headBuffer[payloadLen+protocolLen+codeTypeLen:], logic)
	_, _ = buffer.Write(headBuffer)
	_, _ = buffer.Write(body)
	return buffer
}
