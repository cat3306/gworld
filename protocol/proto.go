package protocol

import (
	"encoding/binary"
	"errors"
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
	payload := make([]byte, len(buf))
	copy(payload, buf)
	packet := &Context{
		Payload:  payload,
		CodeType: CodeType(codeType),
		Proto:    protocol,
		Conn:     c,
	}
	return packet, nil
}
func Encode(v interface{}, codeType CodeType, proto uint32) []byte {
	if v == nil {
		panic("v nil")
	}
	raw, err := GameCoder(codeType).Marshal(v)
	if err != nil {
		panic(err)
	}
	bodyOffset := int(payloadLen + protocolLen + codeTypeLen)
	msgLen := bodyOffset + len(raw)
	buffer := *BUFFERPOOL.Get(uint32(msgLen))
	//data := make([]byte, msgLen)
	packetEndian.PutUint32(buffer, uint32(len(raw)))
	packetEndian.PutUint32(buffer[payloadLen:], proto)
	packetEndian.PutUint16(buffer[payloadLen+protocolLen:], uint16(codeType))
	copy(buffer[bodyOffset:msgLen], raw)
	return buffer[:msgLen]
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
