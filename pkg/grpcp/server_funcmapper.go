package grpcp

import (
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCFuncMapperServer .
type GRPCFuncMapperServer struct {
	client protodef.FuncMapperClient
}

// Do handler.
func (l *GRPCFuncMapperServer) Do(requestID string, key string, args []interface{}) (interface{}, error) {
	var err error
	ctx := context.Background()

	arr := make([]*anypb.Any, len(args))
	for i, v := range args {
		arr[i], err = InterfaceToProtobufAny(v)
		if err != nil {
			return nil, err
		}
	}

	resp, err := l.client.Do(ctx, &protodef.FuncMapperDoRequest{
		Key:       key,
		Requestid: requestID,
		Params:    arr,
	})
	if err != nil {
		return nil, err
	}

	var i interface{}
	err = ProtobufAnyToInterface(resp.Value, &i)
	return i, err
}
