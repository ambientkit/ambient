package grpcp

import (
	"encoding/json"
	"net/http"

	"github.com/ambientkit/ambient"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
)

// HeaderToProtobufStruct converts a http.Header to a protobuf struct.
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

// GrantBoolMapToProtobufStruct converts a Grant map to a protobuf struct.
func GrantBoolMapToProtobufStruct(m map[ambient.Grant]bool) (*structpb.Struct, error) {
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

// ProtobufStructToGrantBoolMap converts a protobuf struct to a Grant map.
func ProtobufStructToGrantBoolMap(s *structpb.Struct) (map[ambient.Grant]bool, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[ambient.Grant]bool)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}

// PluginDataMapToProtobufStruct converts a PluginData map to a protobuf struct.
func PluginDataMapToProtobufStruct(m map[string]ambient.PluginData) (*structpb.Struct, error) {
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

// ProtobufStructToPluginDataMap converts a protobuf struct to a PluginData map.
func ProtobufStructToPluginDataMap(s *structpb.Struct) (map[string]ambient.PluginData, error) {
	b, err := protojson.Marshal(s)
	if err != nil {
		return nil, err
	}
	m := make(map[string]ambient.PluginData)
	err = json.Unmarshal(b, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
