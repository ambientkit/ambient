package grpcp

import (
	"fmt"
	"net/http"
	"reflect"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// FuncMapper handler.
type FuncMapper interface {
	Do(globalFuncMap bool, requestID string, key string, args []interface{}, method string, path string, headers http.Header, body []byte) (value interface{}, errText string, err error)
}

// GRPCFuncMapperPlugin is the gRPC server that GRPCClient talks to.
type GRPCFuncMapperPlugin struct {
	Impl FuncMapper
	Log  ambient.Logger
}

// Do handler.
func (m *GRPCFuncMapperPlugin) Do(ctx context.Context, req *protodef.FuncMapperDoRequest) (resp *protodef.FuncMapperDoResponse, err error) {
	params := make([]interface{}, len(req.Params))
	for i, v := range req.Params {
		err = ProtobufAnyToInterface(v, &params[i])
		if err != nil {
			return &protodef.FuncMapperDoResponse{}, fmt.Errorf("Do any error: %v", err.Error())
		}
	}

	var headers http.Header
	err = ProtobufStructToObject(req.Headers, &headers)
	if err != nil {
		return nil, fmt.Errorf("Do header conversion error: %v", err.Error())
	}

	val, errText, err := m.Impl.Do(req.Globalfm, req.Requestid, req.Key, params, req.Method, req.Path, headers, req.Body)
	if err != nil {
		return &protodef.FuncMapperDoResponse{}, fmt.Errorf("Do error: %v", err.Error())
	}

	// if val != nil {
	// 	m.Log.Error("PLUGIN: Kind: %v | Value: %v | Valid: %v", reflect.TypeOf(val).Kind(), val, reflect.ValueOf(val).IsValid())
	// } else {
	// 	m.Log.Error("PLUGIN: Type: %v | Value: %v", reflect.TypeOf(val), val)
	// }

	arr := make([]*structpb.Struct, 0)
	var anyVal *anypb.Any
	if !reflect.ValueOf(val).IsValid() {
		return &protodef.FuncMapperDoResponse{
			Value: nil,
			Arr:   nil,
			Error: errText,
		}, nil
	} else if reflect.TypeOf(val).Kind() == reflect.Slice {
		arr, err = ArrayToProtobufStruct(val)
		if err != nil {
			return &protodef.FuncMapperDoResponse{}, fmt.Errorf("Do array error: %v", err.Error())
		}
	} else {
		anyVal, err = InterfaceToProtobufAny(val)
		if err != nil {
			return &protodef.FuncMapperDoResponse{}, fmt.Errorf("Do interface error: %v", err.Error())
		}
	}

	return &protodef.FuncMapperDoResponse{
		Value: anyVal,
		Arr:   arr,
		Error: errText,
	}, err
}
