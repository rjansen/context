package proto

import (
	"fmt"
	// "errors"
	"farm.e-pedion.com/repo/logger"
	// "github.com/golang/protobuf/proto"
	// "github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func init() {
	os.Args = append(os.Args, "-ecf", "../../test/etc/context/context.yaml")
	logger.Info("context.media.proto_test.init")
}

func TestProtoMarshalBytes(t *testing.T) {
	p := &Person{
		Name:  "Proto Buffer Person",
		Id:    1,
		Email: "proto.buffer.person@mock.com",
	}

	b, e := MarshalBytes(p)
	fmt.Printf("marshal.protoMessage len=%d err=%v\n", len(b), e)

	e = UnmarshalBytes(b, p)
	fmt.Printf("unmarshal.protoMessage Person=%s err=%v\n", p, e)
}
