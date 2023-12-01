package serializer

import (
	//"github.com/golang/protobuf/jsonpb"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"
)

// ProtobufToJSON converts protocol buffer message to JSON string
func ProtobufToJSON(message proto.Message) (string, error) {
	jsonString := protojson.Format(message)

	return jsonString, nil
}

// JSONToProtobufMessage converts JSON string to protocol buffer message
func JSONToProtobufMessage(data string, message proto.Message) error {
	bytes1 := []byte(data)
	err := protojson.Unmarshal(bytes1, message)
	return err
}
