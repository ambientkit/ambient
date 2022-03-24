package grpcp

import (
	"fmt"
	"reflect"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

// FuncMapper handler.
type FuncMapper interface {
	Do(requestID string, key string, args []interface{}) (value interface{}, err error)
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
			return &protodef.FuncMapperDoResponse{}, fmt.Errorf("grpc-plugin: Do any error: %v", err.Error())
		}
	}

	val, err := m.Impl.Do(req.Requestid, req.Key, params)
	if err != nil {
		return &protodef.FuncMapperDoResponse{}, fmt.Errorf("grpc-plugin: Do error: %v", err.Error())
	}

	//m.Log.Error("Kind: %v %v", reflect.TypeOf(val).Kind(), val)

	arr := make([]*structpb.Struct, 0)
	var anyVal *anypb.Any
	if reflect.TypeOf(val).Kind() == reflect.Slice {
		arr, err = ArrayToProtobufStruct(val)
		if err != nil {
			return &protodef.FuncMapperDoResponse{}, fmt.Errorf("grpc-plugin: Do array error: %v", err.Error())
		}
	} else {
		anyVal, err = InterfaceToProtobufAny(val)
		if err != nil {
			return &protodef.FuncMapperDoResponse{}, fmt.Errorf("grpc-plugin: Do interface error: %v", err.Error())
		}
	}

	return &protodef.FuncMapperDoResponse{
		Value: anyVal,
		Arr:   arr,
	}, err
}
