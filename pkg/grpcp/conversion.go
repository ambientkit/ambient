package grpcp

import (
	"encoding/json"
	"net/http"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// HeaderToProtobufStruct converts http.Header to a protobuf struct.
func HeaderToProtobufStruct(m http.Header) (*structpb.Struct, error) {
	b, err := json.Marshal(m)
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

// ProtobufStructToHeader converts a protobuf struct to a http.Header.
func ProtobufStructToHeader(s *structpb.Struct) (http.Header, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(http.Header)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
