package grpcp

import (
	"encoding/json"
	"errors"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// ObjectToProtobufStruct converts an object to a protobuf struct.
func ObjectToProtobufStruct(obj interface{}) (*structpb.Struct, error) {
	b, err := json.Marshal(obj)
	if err != nil {
		return nil, err
	}
	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// ProtobufStructToObject converts a protobuf struct to an object.
func ProtobufStructToObject(s *structpb.Struct, obj interface{}) error {
	// Check for obj - this is a developer bug.
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value")
	}
	b, err := protojson.Marshal(s)
	if err != nil {
		return err
	}
	return json.Unmarshal(b, &obj)
}
