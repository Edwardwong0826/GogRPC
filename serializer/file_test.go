package serializer

import (
	wongProto "gogrpc/pb"
	"gogrpc/sample"
	"testing"

	//"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/proto"
)

// install testify package first
// this is like junit assert in Java
// go get github.com/stretchr/testify
// click left button to run test, dont click run test it will hang at ok dont know why

// after run the test and validate the laptop.bin and laptop.json result, this go project generated laptop.bin
// can use in other gRPC project such as Java to test the other generated laptop.json is the same or not
func TestFileSerializer(t *testing.T) {
	t.Parallel()

	binaryFile := "../tmp/laptop.bin"
	jsonFile := "../tmp/laptop.json"

	laptop1 := sample.NewLaptop()

	err := WriteProtobufToBinaryFile(laptop1, binaryFile)
	require.NoError(t, err)

	laptop2 := &wongProto.Laptop{}
	err = ReadProtobufFromBinaryFile(binaryFile, laptop2)
	require.NoError(t, err)
	require.True(t, proto.Equal(laptop1, laptop2))

	err = WriteProtobufToJSONFile(laptop1, jsonFile)
	require.NoError(t, err)

}
