package io

import (
	"github.com/golang/protobuf/proto"
	// "bitbucket.org/fungrim/typeframed-go"
	msg "bitbucket.org/fungrim/typeframed-go/test"
	tfd "bitbucket.org/fungrim/typeframed-go"
	"testing"
	"bytes"
	// "fmt"
	"strconv"
)

type TestDictionary struct{}

func (d *TestDictionary) GetId(msg proto.Message) (int, error) {
	return 666, nil
}

func (d *TestDictionary) NewMessageFromId(id int) (proto.Message, error) {
	return &msg.Tell{}, nil
}

func TestSimple(t *testing.T) {
	td := &TestDictionary{}
	var buff bytes.Buffer
	writer := NewStreamWriter(&buff, td, nil)
	tell := &msg.Tell{Msg:proto.String("Hello World!")}
	err := writer.Write(tell, nil)
	if err != nil {
		t.Errorf("failed write with message: %v", err)
	}
	rawMsg := buff.Bytes();
	javaRawMsg := []byte{154, 5, 0, 14, 10, 12, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 0}
	if !bytes.Equal(rawMsg, javaRawMsg) {
		t.Errorf("encoding incompatible with java verification")
	}
	reader := NewStreamReader(bytes.NewReader(rawMsg), td, nil, nil)
	tell2, err := reader.Read()
	if err != nil {
		t.Errorf("failed read with message: %v", err)
	}
	if tell.GetMsg() != tell2.(*msg.Tell).GetMsg() {
		t.Errorf("messages not equal; %v != %v", tell.GetMsg(), tell2.(*msg.Tell).GetMsg())
	}
}

func TestSimpleWithChecksum(t *testing.T) {
	td := &TestDictionary{}
	var buff bytes.Buffer
	writer := NewStreamWriter(&buff, td, tfd.Crc32Checksum)
	tell := &msg.Tell{Msg:proto.String("Hello World!")}
	err := writer.Write(tell, nil)
	if err != nil {
		t.Errorf("failed write with message: %v", err)
	}
	rawMsg := buff.Bytes();
	javaRawMsg := []byte{154, 5, 0, 14, 10, 12, 72, 101, 108, 108, 111, 32, 87, 111, 114, 108, 100, 33, 8, 0, 0, 0, 0, 254, 202, 192, 107}
	if !bytes.Equal(rawMsg, javaRawMsg) {
		t.Errorf("encoding incompatible with java verification")
	}
	reader := NewStreamReader(bytes.NewReader(rawMsg), td, nil, tfd.Crc32Checksum)
	tell2, err := reader.Read()
	if err != nil {
		t.Errorf("failed read with message: %v", err)
	}
	if tell.GetMsg() != tell2.(*msg.Tell).GetMsg() {
		t.Errorf("messages not equal; %v != %v", tell.GetMsg(), tell2.(*msg.Tell).GetMsg())
	}
}

func TestSimpleMultiple(t *testing.T) {
	td := &TestDictionary{}
	var buff bytes.Buffer
	writer := NewStreamWriter(&buff, td, nil)
	for i := 0; i < 5; i++ {
		tell := &msg.Tell{Msg:proto.String(strconv.Itoa(i))}
		err := writer.Write(tell, nil)
		if err != nil {
			t.Errorf("failed write with message: %v", err)
		}
	}
	rawMsg := buff.Bytes();
	reader := NewStreamReader(bytes.NewReader(rawMsg), td, nil, nil)
	for i := 0; i < 5; i++ {
		tell2, err := reader.Read()
		if err != nil {
			t.Errorf("failed read with message: %v", err)
		}
		if strconv.Itoa(i) != tell2.(*msg.Tell).GetMsg() {
			t.Errorf("messages not equal; %v != %v", strconv.Itoa(i), tell2.(*msg.Tell).GetMsg())
		}
	}
}

func TestReadPastEnd(t *testing.T) {
	td := &TestDictionary{}
	var buff bytes.Buffer
	writer := NewStreamWriter(&buff, td, nil)
	tell := &msg.Tell{Msg:proto.String("Hello World!")}
	err := writer.Write(tell, nil)
	if err != nil {
		t.Errorf("failed write with message: %v", err)
	}
	rawMsg := buff.Bytes();
	reader := NewStreamReader(bytes.NewReader(rawMsg), td, nil, nil)
	_, err = reader.Read()
	if err != nil {
		t.Errorf("failed read with message: %v", err)
	}
	_, err = reader.Read()
	if err == nil {
		t.Errorf("expected EOF error")
	}
}