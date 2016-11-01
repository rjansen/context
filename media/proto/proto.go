package proto

import (
	"errors"
	"farm.e-pedion.com/repo/logger"
	"github.com/golang/protobuf/proto"
	"io"
)

var (
	ErrInvalidProtoMessage = errors.New("Invalid proto message value")
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
	// return proto.NewEncoder(w).Encode(&data)
	return nil
}

//Unmarshal reads a json representation into the struct instance
func Unmarshal(r io.Reader, result interface{}) error {
	// return json.NewDecoder(r).Decode(&result)
	return nil
}

//MarshalBytes writes a json representation of the struct instance
func MarshalBytes(data interface{}) ([]byte, error) {
	msg, err := protoMessage(data)
	if err != nil {
		return nil, err
	}
	protoBytes, err := proto.Marshal(msg)
	logger.Debug("proto.MarshalBytes",
		logger.Bytes("protoBytes", protoBytes),
		logger.Err(err),
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
	logger.Debug("proto.UnmarshalBytes",
		logger.Bool("resultIsNil", result == nil),
		logger.Err(err),
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
