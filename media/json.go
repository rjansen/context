package media

import (
	"encoding/json"
	"farm.e-pedion.com/repo/logger"
	"io"
)

//JSON is an interface to defines a json media encoder/decoder contract
type JSON interface {
	//Marshal writes a json representation of the struct instance
	Marshal(io.Writer) error
	//Unmarshal reads a json representation into the struct instance
	Unmarshal(io.Reader) error
	//ToBytes writes a json representation of the struct instance
	ToBytes() ([]byte, error)
	//FromBytes reads a json representation into the struct instance
	FromBytes([]byte) error
}

//JSONMarshal writes a json representation of the struct instance
func JSONMarshal(w io.Writer, data interface{}) error {
	return json.NewEncoder(w).Encode(&data)
}

//JSONUnmarshal reads a json representation into the struct instance
func JSONUnmarshal(r io.Reader, result interface{}) error {
	return json.NewDecoder(r).Decode(&result)
}

//ToJSONBytes writes a json representation of the struct instance
func ToJSONBytes(data interface{}) ([]byte, error) {
	jsonBytes, err := json.Marshal(data)
	logger.Debug("media.ToJSONBytes",
		logger.Bytes("jsonBytes", jsonBytes),
		logger.Err(err),
	)
	return jsonBytes, err
}

//FromJSONBytes reads a json representation into the struct instance
func FromJSONBytes(raw []byte, result interface{}) error {
	err := json.Unmarshal(raw, &result)
	logger.Debug("media.FromJSONBytes",
		logger.Bool("resultIsNil", result == nil),
		logger.Err(err),
	)
	return err
}

//JSONMedia is a struct to helps writes and reads of a json representation
type JSONMedia struct {
}

//Marshal writes a json representation of the struct instance
func (m JSONMedia) Marshal(writer io.Writer) error {
	return JSONMarshal(writer, &m)
}

//Unmarshal reads a json representation into the struct instance
func (m JSONMedia) Unmarshal(reader io.Reader) error {
	return JSONUnmarshal(reader, &m)
}

//ToBytes writes a json representation of the struct instance
func (m JSONMedia) ToBytes() ([]byte, error) {
	return ToJSONBytes(&m)
}

//FromBytes reads a json representation into the struct instance
func (m JSONMedia) FromBytes(raw []byte) error {
	return FromJSONBytes(raw, &m)
}
