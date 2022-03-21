package grpcp

import (
	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"golang.org/x/net/context"
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
			return &protodef.FuncMapperDoResponse{}, err
		}
	}

	val, err := m.Impl.Do(req.Requestid, req.Key, params)
	if err != nil {
		return &protodef.FuncMapperDoResponse{}, err
	}

	anyVal, err := InterfaceToProtobufAny(val)
	if err != nil {
		return &protodef.FuncMapperDoResponse{}, err
	}

	return &protodef.FuncMapperDoResponse{
		Value: anyVal,
	}, err
}
