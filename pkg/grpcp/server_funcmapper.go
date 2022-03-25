package grpcp

import (
	"fmt"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCFuncMapperServer handles server calls for FuncMap.
type GRPCFuncMapperServer struct {
	client protodef.FuncMapperClient
	Log    ambient.Logger
}

// Do handler.
func (l *GRPCFuncMapperServer) Do(requestID string, key string, args []interface{}) (interface{}, error) {
	var err error
	ctx := context.Background()

	arr := make([]*anypb.Any, len(args))
	for i, v := range args {
		arr[i], err = InterfaceToProtobufAny(v)
		if err != nil {
			return nil, fmt.Errorf("grpc-server: Do args error: %v", err.Error())
		}
	}

	resp, err := l.client.Do(ctx, &protodef.FuncMapperDoRequest{
		Key:       key,
		Requestid: requestID,
		Params:    arr,
	})
	if err != nil {
		return nil, fmt.Errorf("grpc-server: Do response error:%v", err.Error())
	}

	//l.Log.Error("SERVER: Kind: %v | Value: %v | Valid: %v | Args: %#v", reflect.TypeOf(resp.Value).Kind(), resp.Value, reflect.ValueOf(resp.Value).IsValid(), args)

	// switch resp.vl.(type) {

	// }

	var i interface{}
	if resp.Value != nil {
		err = ProtobufAnyToInterface(resp.Value, &i)
	} else {
		err = ProtobufStructToArray(resp.Arr, &i)
	}

	//l.Log.Error("SERVED: Kind: %v | Value: %v | Valid: %v", reflect.TypeOf(i).Kind(), i, reflect.ValueOf(i).IsValid())
	//l.Log.Error("SERVED: Type: %v | Value: %v", reflect.TypeOf(i), i)

	return i, err
}
