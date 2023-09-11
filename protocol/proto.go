package protocol

import (
	"encoding/binary"
	"errors"
	"github.com/valyala/bytebufferpool"
	"io"

	"github.com/panjf2000/gnet/v2"
)

//
// * 0                       4                       8           10
// * +-----------------------+-----------------------+-----------+
// * |   body len    		 |       protocol        | code type |
// * +-----------------------+-----------------------+-----------+
// * |                                   			 			 |
// * +                                       		             +
// * |                   body bytes              		       	 |
// * +                                   						 +
// * |                                  						 |
// * +-----------------------------------------------------------+

const (
	payloadLen  = uint32(4)
	protocolLen = uint32(4)
	codeTypeLen = uint32(2)
)

var (
	packetEndian        = binary.LittleEndian
	ErrIncompletePacket = errors.New("incomplete packet")
	ErrTooLargePacket   = errors.New("too large packet")
)

func Decode(c gnet.Conn) (*Context, error) {

	bodyOffset := int(payloadLen + protocolLen + codeTypeLen)
	buf, err := c.Next(bodyOffset)
	if err != nil {
		return nil, err
	}

	bodyLen := packetEndian.Uint32(buf[:payloadLen])
	protocol := packetEndian.Uint32(buf[payloadLen : payloadLen+protocolLen])
	codeType := packetEndian.Uint16(buf[payloadLen+protocolLen : payloadLen+protocolLen+codeTypeLen])
	msgLen := bodyOffset + int(bodyLen)
	if msgLen > maxByte {
		c.Close()
		return nil, ErrTooLargePacket
	}
	if c.InboundBuffered() < int(bodyLen) {
		return nil, ErrIncompletePacket
	}
	buf, err = c.Next(int(bodyLen))
	if err != nil {
		return nil, err
	}
	buffer := bytebufferpool.Get()
	_, _ = buffer.Write(buf)
	//payload := make([]byte, len(buf))
	//copy(payload, buf)
	packet := &Context{
		Payload:  buffer,
		CodeType: CodeType(codeType),
		Proto:    protocol,
		Conn:     c,
	}
	return packet, nil
}
func Encode(v interface{}, codeType CodeType, proto uint32) *bytebufferpool.ByteBuffer {
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
	headBuffer := make([]byte, payloadLen+protocolLen+codeTypeLen)
	packetEndian.PutUint32(headBuffer, uint32(len(body)))
	packetEndian.PutUint32(headBuffer[payloadLen:], proto)
	packetEndian.PutUint16(headBuffer[payloadLen+protocolLen:], uint16(codeType))
	_, _ = buffer.Write(headBuffer)
	_, _ = buffer.Write(body)
	return buffer
}

func ReadFull(r io.Reader) ([]byte, uint32, uint16, error) {
	preBuff := make([]byte, 10)
	_, err := io.ReadFull(r, preBuff)
	if err != nil {
		return nil, 0, 0, err
	}
	bodyLen := packetEndian.Uint32(preBuff[:payloadLen])
	protocol := packetEndian.Uint32(preBuff[payloadLen : payloadLen+protocolLen])
	codeType := packetEndian.Uint16(preBuff[payloadLen+protocolLen : payloadLen+protocolLen+codeTypeLen])
	payload := make([]byte, bodyLen)
	_, err = io.ReadFull(r, payload)
	if err != nil {
		return nil, 0, 0, err
	}
	return payload, protocol, codeType, nil
}
