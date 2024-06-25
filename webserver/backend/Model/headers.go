package communications

import (
	"encoding/binary"
	"errors"
)

type Header struct {
	magic  uint32
	length uint32
}

type Communication struct {
	header  Header
	message string
}

func NewCommunication(message string, magic uint32) *Communication {
	return &Communication{
		message: message,
		header:  *NewHeaderFromValue(magic, uint32(len(message))),
	}
}

func (c *Communication) AsByte() []byte {
	var bytes = append((&c.header).AsByte(), []byte(c.message)...)
	return bytes
}

func NewHeaderFromValue(magic uint32, length uint32) *Header {
	return &Header{
		magic:  magic,
		length: length,
	}
}

func NewHeaderFromBytes(bytes []byte) (*Header, error) {
	if len(bytes) < 8 {
		return nil, errors.New("header should be of length 8")
	}

	return &Header{
		magic:  binary.LittleEndian.Uint32(bytes[0:4]),
		length: binary.LittleEndian.Uint32(bytes[4:8]),
	}, nil
}

func (h *Header) AsByte() []byte {
	var bytes [8]byte = [8]byte{}
	binary.LittleEndian.PutUint32(bytes[0:4], h.magic)
	binary.LittleEndian.PutUint32(bytes[4:8], h.length)
	return bytes[:]
}

func (h *Header) GetMessageLength() uint32 {
	return h.length
}

func (h *Header) GetHeaderMagic() uint32 {
	return h.magic
}
