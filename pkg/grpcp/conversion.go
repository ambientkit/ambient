package grpcp

import (
	"encoding/json"
	"errors"
	"fmt"
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
	} else if val.Elem().Kind() != reflect.Slice {
		return fmt.Errorf("must pass a slice, not a: %v", val.Elem().Kind())
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
