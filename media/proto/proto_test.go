package proto

import (
	"fmt"
	// "errors"
	"farm.e-pedion.com/repo/logger"
	// "github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Args = append(os.Args, "-ecf", "../../test/etc/context/context.yaml")
	logger.Info("context.media.proto_test.init")
}

func TestProtoMarshalBytes(t *testing.T) {
	p := &Store{
		Id:   1,
		Name: "Proto Buffer Store",
		Data: []*Store_Data{
			&Store_Data{
				Id:    1,
				Name:  "Proto Data Name",
				Email: "Proto Data Email",
			},
		},
	}

	b, e := MarshalBytes(p)
	fmt.Printf("marshal.protoMessage len=%d err=%v\n", len(b), e)
	assert.Nil(t, e)

	e = UnmarshalBytes(b, p)
	fmt.Printf("unmarshal.protoMessage store=%s err=%v\n", p, e)
	assert.Nil(t, e)
}

func TestProtoMarshalBytesError(t *testing.T) {
	p := map[string]interface{}{
		"Id":   1,
		"Name": "Proto Buffer Store",
		"Data": []map[string]interface{}{
			map[string]interface{}{
				"Id":    1,
				"Name":  "Proto Data Name",
				"Email": "Proto Data Email",
			},
		},
	}

	_, e := MarshalBytes(p)
	assert.NotNil(t, e)
}

func TestProtoUnmarshalBytesError(t *testing.T) {
	p := make(map[string]interface{})

	data := []byte{
		0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e,
		0x14, 0xc7, 0x4d, 0x9a, 0x94, 0xf9, 0x26,
		0xdd, 0x91, 0xcc, 0x3f, 0x63, 0x01, 0x00, 0x00,
	}

	e := UnmarshalBytes(data, p)
	assert.NotNil(t, e)
}
