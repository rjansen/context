package media

import (
	"encoding/json"
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

//JSONMedia is a struct to helps writes and reads of a json representation
type JSONMedia struct {
}

//Marshal writes a json representation of the struct instance
func (m JSONMedia) Marshal(writer io.Writer) error {
	return json.NewEncoder(writer).Encode(&m)
}

//Unmarshal reads a json representation into the struct instance
func (m JSONMedia) Unmarshal(reader io.Reader) error {
	return json.NewDecoder(reader).Decode(&m)
}

//ToBytes writes a json representation of the struct instance
func (m JSONMedia) ToBytes() ([]byte, error) {
	return json.Marshal(m)
}

//FromBytes reads a json representation into the struct instance
func (m JSONMedia) FromBytes(data []byte) error {
	return json.Unmarshal(data, &m)
}
