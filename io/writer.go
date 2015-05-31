package io

import (
	"bitbucket.org/fungrim/go.typeframed"
	"encoding/binary"
	"github.com/golang/protobuf/proto"
	//"bytes"
	"io"
	//"errors"
	//"fmt"
)

type StreamWriter struct {
	writer     io.Writer
	dictionary typeframed.MessageTypeDictionary
	chksum     typeframed.Checksum
}

func NewStreamWriter(writer io.Writer, dictionary typeframed.MessageTypeDictionary, chksum typeframed.Checksum) *StreamWriter {
	return &StreamWriter{writer, dictionary, chksum}
}

func (w *StreamWriter) Write(msg proto.Message, header []byte) error {

	b := newBuffer()

	// 1: type ID
	typeId, err := w.dictionary.GetId(msg)
	if err != nil {
		return err
	}
	b.writeVarInt(typeId)

	// 2: optional header
	if header == nil {
		b.writeVarInt(0)
	} else {
		b.writeVarInt(len(header))
		b.write(header)
	}

	// 3: msg
	msgBytes, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	b.writeVarInt(len(msgBytes))
	b.write(msgBytes)

	// 4: optional checksum
	if w.chksum == nil {
		b.writeVarInt(0)
	} else {
		chksum := w.chksum(msgBytes)
		b.writeVarInt(len(chksum))
		b.write(chksum)
	}

	_, err = w.writer.Write(b.slice()) // TODO check length?

	return err
}

type buffer struct {
	buf      []byte
	position int
}

func newBuffer() *buffer {
	return newBufferWithLength(512)
}

func newBufferWithLength(length int) *buffer {
	return &buffer{make([]byte, length), 0}
}

func (b *buffer) slice() []byte {
	return b.buf[:b.position]
}

func (b *buffer) ensure(length int) {
	if b.remaining() < length {
		b.grow(length)
	}
}

func (b *buffer) writeVarInt(n int) {
	b.ensure(binary.MaxVarintLen32)
	varint := proto.EncodeVarint(uint64(n))
	copy(b.buf[b.position:], varint)
	b.position += len(varint)
}

func (b *buffer) write(bytes []byte) {
	b.ensure(len(bytes))
	copy(b.buf[b.position:], bytes)
	b.position += len(bytes)
}

func (b *buffer) grow(length int) {
	tmp := make([]byte, (cap(b.buf)*2)+length)
	copy(tmp, b.buf)
	b.buf = tmp
}

func (b *buffer) remaining() int {
	return cap(b.buf) - b.position
}
