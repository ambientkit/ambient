package grpcp

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"reflect"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// ObjectToProtobufStruct converts an object to a protobuf struct.
func ObjectToProtobufStruct(obj interface{}) (*structpb.Struct, error) {
	if obj == nil {
		return &structpb.Struct{}, nil
	}

	if o, ok := obj.(map[string]interface{}); ok && len(o) == 0 {
		return &structpb.Struct{}, nil
	}

	if o, ok := obj.(http.Header); ok && len(o) == 0 {
		return &structpb.Struct{}, nil
	}

	b, err := json.Marshal(obj)
	if err != nil {
		return &structpb.Struct{}, fmt.Errorf("could not marshal: %v", err.Error())
	}

	s := &structpb.Struct{}
	err = protojson.Unmarshal(b, s)
	if err != nil {
		return &structpb.Struct{}, fmt.Errorf("could not unmarshal: %v ", err.Error())
	}

	return s, nil
}

// ProtobufStructToObject converts a protobuf struct to an object.
func ProtobufStructToObject(s *structpb.Struct, obj interface{}) error {
	// Check for pointer - this is a developer bug if the error occurs.
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

// ArrayToProtobufStruct converts an array to a protobuf struct.
func ArrayToProtobufStruct(obj interface{}) ([]*structpb.Struct, error) {
	// Check for slice - this is a developer bug if the error occurs.
	val := reflect.TypeOf(obj)
	if val.Kind() != reflect.Slice {
		return nil, fmt.Errorf("must pass a slice, not a: %v", val.Kind())
	}

	arr := reflect.ValueOf(obj)
	s := make([]*structpb.Struct, arr.Len())
	for i := 0; i < arr.Len(); i++ {
		b, err := json.Marshal(arr.Index(i).Interface())
		if err != nil {
			return nil, err
		}
		vv := &structpb.Struct{}
		err = protojson.Unmarshal(b, vv)
		if err != nil {
			return nil, err
		}
		s[i] = vv
	}

	return s, nil
}

// ProtobufStructToArray converts a protobuf struct to an array.
func ProtobufStructToArray(s []*structpb.Struct, obj interface{}) error {
	// Check for pointer - this is a developer bug if the error occurs.
	val := reflect.ValueOf(obj)
	if val.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a: value")
	}

	objVal := val.Elem()

	for i, v := range s {
		b, err := protojson.Marshal(v)
		if err != nil {
			return err
		}

		objVal.Set(reflect.Append(objVal, reflect.New(objVal.Type().Elem()).Elem()))
		err = json.Unmarshal(b, objVal.Index(i).Addr().Interface())
		if err != nil {
			return err
		}
	}

	return nil
}

// InterfaceToProtobufAny converts an interface to a protobuf any.
func InterfaceToProtobufAny(obj interface{}) (*anypb.Any, error) {
	nv, err := structpb.NewValue(obj)
	if err != nil {
		return nil, err
	}

	s, err := anypb.New(nv)
	if err != nil {
		return nil, err
	}

	return s, nil
}

// ProtobufAnyToInterface converts a protobuf any to an interface.
func ProtobufAnyToInterface(s *anypb.Any, obj interface{}) error {
	// Check for pointer - this is a developer bug if the error occurs.
	v := reflect.ValueOf(obj)
	if v.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value")
	}

	var val structpb.Value
	err := s.UnmarshalTo(&val)
	if err != nil {
		return err
	}

	b, err := protojson.Marshal(&val)
	if err != nil {
		return err
	}

	return json.Unmarshal(b, &obj)
}

// MapToProtobufStruct converts a map to a protobuf struct.
func MapToProtobufStruct(obj map[string]interface{}) (*structpb.Struct, error) {
	return structpb.NewStruct(obj)
}

// ProtobufStructToMap converts a map to a protobuf struct.
func ProtobufStructToMap(s *structpb.Struct) map[string]interface{} {
	return s.AsMap()
}
