package io

import (
	"bitbucket.org/fungrim/go.typeframed"
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"io"
)

type StreamReader struct {
	reader     *bytes.Reader
	dictionary typeframed.MessageTypeDictionary
	header     typeframed.HeaderCapture
	chksum     typeframed.Checksum
}

func NewStreamReader(reader *bytes.Reader, dictionary typeframed.MessageTypeDictionary, header typeframed.HeaderCapture, chksum typeframed.Checksum) *StreamReader {
	return &StreamReader{reader, dictionary, header, chksum}
}

func (r *StreamReader) Read() (proto.Message, error) {

	// 1: read type
	typeId, err := readVarint(r.reader)
	if err != nil {
		return nil, err
	}

	// 2: read optional header
	headLen, err := readVarint(r.reader)
	if err != nil {
		return nil, err
	} else if headLen > 0 {
		if r.header == nil {
			return nil, errors.New("Incoming envelope holds a header, but no header capture is set in the reader")
		}
		rawHead, err := readBytes(r.reader, int(headLen))
		if err != nil {
			return nil, err
		}
		if err := r.header.Capture(rawHead); err != nil {
			return nil, err
		}
	}

	// 3: read data length and data
	msgLen, err := readVarint(r.reader)
	rawMsg, err := readBytes(r.reader, int(msgLen))
	if err != nil {
		return nil, err
	}

	// 4: read optional checksum
	cksumLen, err := readVarint(r.reader)
	if err != nil {
		return nil, err
	} else if cksumLen > 0 {
		if r.chksum == nil {
			return nil, errors.New("Incoming envelope holds a checksum, but no checksum provider is set in the reader")
		}
		rawChksum, err := readBytes(r.reader, int(cksumLen))
		if err != nil {
			return nil, err
		}
		testChksum := r.chksum(rawMsg)
		if !bytes.Equal(rawChksum, testChksum) {
			return nil, typeframed.NewCorruptedChecksum()
		}
	}

	// 5: create new message and parse
	msg, err := r.dictionary.NewMessageFromId(int(typeId))
	if err != nil {
		return nil, err
	} else {
		err = proto.Unmarshal(rawMsg, msg)
		return msg, err
	}
}

func readBytes(reader io.Reader, length int) ([]byte, error) {
	buf := make([]byte, length)
	_, err := reader.Read(buf) // TODO check len?
	return buf, err
}

/*
 * Copied with variations from: https://github.com/golang/protobuf/blob/master/proto/decode.go
 * Copyright 2010 The Go Authors
 */
func readVarint(reader io.ByteReader) (x uint64, err error) {
	for shift := uint(0); shift < 64; shift += 7 {
		var b byte
		b, err = reader.ReadByte()
		if err != nil {
			return
		}
		x |= (uint64(b) & 0x7F) << shift
		if b < 0x80 {
			return
		}
	}
	// The number is too large to represent in a 64-bit value.
	err = errors.New("proto: integer overflow")
	return
}
