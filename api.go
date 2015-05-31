package typeframed

import (
	"github.com/golang/protobuf/proto"
	"strconv"
)

type MessageTypeDictionary interface {
	GetId(msg proto.Message) (int, error)
	NewMessageFromId(id int) (proto.Message, error)
}

type StreamReader interface {
	Read() (proto.Message, error)
}

type StreamWriter interface {
	Write(msg proto.Message, header []byte) error
}

type NoSuchTypeError struct {
	Id int
}

type UnknownMessageError struct {
	Msg string
}

func (e *UnknownMessageError) Error() string {
	return e.Msg
}

func (e *NoSuchTypeError) Error() string {
	return "No such type: " + strconv.Itoa(e.Id)
}

func NewNoSuchType(id int) error {
	return &NoSuchTypeError{id}
}

func NewUnknownMessage(msg string) error {
	return &UnknownMessageError{msg}
}
