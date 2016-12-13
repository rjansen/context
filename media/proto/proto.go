package proto

import (
	"bytes"
	"errors"
	"github.com/golang/protobuf/proto"
	"github.com/rjansen/l"
	"io"
)

var (
	//ErrInvalidProtoMessage is returned when a invalid message is provided
	ErrInvalidProtoMessage = errors.New("Invalid proto message value")
	//ErrEmptyInput is returned when a empty input is provided
	ErrEmptyInput = errors.New("The provided input (reader, []bytes) is empty")
	//ContentType is a constant to hold the protocol buffer content type value
	ContentType = "application/octet-stream"
)

func protoMessage(val interface{}) (proto.Message, error) {
	msg, ok := val.(proto.Message)
	if !ok {
		return nil, ErrInvalidProtoMessage
	}
	return msg, nil
}

//Marshal writes a json representation of the struct instance
func Marshal(w io.Writer, data interface{}) error {
	dataBytes, err := MarshalBytes(data)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(dataBytes)
	_, err = buf.WriteTo(w)
	return err
}

//Unmarshal reads a json representation into the struct instance
func Unmarshal(r io.Reader, result interface{}) error {
	var buf bytes.Buffer
	if reads, err := buf.ReadFrom(r); err != nil {
		return err
	} else if reads <= 0 {
		return ErrEmptyInput
	}
	return UnmarshalBytes(buf.Bytes(), result)
}

//MarshalBytes writes a json representation of the struct instance
func MarshalBytes(data interface{}) ([]byte, error) {
	msg, err := protoMessage(data)
	if err != nil {
		return nil, err
	}
	protoBytes, err := proto.Marshal(msg)
	l.Debug("proto.MarshalBytes",
		l.Int("len", len(protoBytes)),
		l.Err(err),
	)
	return protoBytes, err
}

//UnmarshalBytes reads a json representation into the struct instance
func UnmarshalBytes(raw []byte, result interface{}) error {
	msg, err := protoMessage(result)
	if err != nil {
		return err
	}
	err = proto.Unmarshal(raw, msg)
	l.Debug("proto.UnmarshalBytes",
		l.Bool("nilResult", result == nil),
		l.Err(err),
	)
	return err
}

//Media is a struct to helps writes and reads of a json representation
type Media struct {
}

//Marshal writes a json representation of the struct instance
func (Media) Marshal(writer io.Writer, val interface{}) error {
	return Marshal(writer, val)
}

//Unmarshal reads a json representation into the struct instance
func (Media) Unmarshal(reader io.Reader, ref interface{}) error {
	return Unmarshal(reader, ref)
}

//MarshalBytes writes a json representation of the struct instance
func (Media) MarshalBytes(val interface{}) ([]byte, error) {
	return MarshalBytes(val)
}

//UnmarshalBytes reads a json representation into the struct instance
func (Media) UnmarshalBytes(raw []byte, ref interface{}) error {
	return UnmarshalBytes(raw, ref)
}
