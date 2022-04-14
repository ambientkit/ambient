package grpcp

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"google.golang.org/protobuf/types/known/anypb"
)

// GRPCFuncMapperServer handles server calls for FuncMap.
type GRPCFuncMapperServer struct {
	client protodef.FuncMapperClient
	Log    ambient.Logger
}

// Do handler.
func (l *GRPCFuncMapperServer) Do(r *http.Request, requestID string, key string, args []interface{}, globalFuncMap bool) (interface{}, string, error) {
	var err error
	ctx := r.Context()

	arr := make([]*anypb.Any, len(args))
	for i, v := range args {
		arr[i], err = InterfaceToProtobufAny(v)
		if err != nil {
			return nil, "", fmt.Errorf("Do args error: %v", err.Error())
		}
	}

	sm, err := ObjectToProtobufStruct(r.Header)
	if err != nil {
		return nil, "", fmt.Errorf("Do header conversion error: %v", err.Error())
	}

	body := bytes.NewBuffer(nil)
	_, err = io.Copy(body, r.Body)
	if err != nil {
		return nil, "", fmt.Errorf("Do body copy error: %v", err.Error())
	}
	// Restore body.
	r.Body = ioutil.NopCloser(body)

	resp, err := l.client.Do(ctx, &protodef.FuncMapperDoRequest{
		Globalfm:  globalFuncMap,
		Key:       key,
		Requestid: requestID,
		Params:    arr,
		Method:    r.Method,
		Path:      r.RequestURI,
		Headers:   sm,
		Body:      body.Bytes(),
	})
	if err != nil {
		return nil, "", fmt.Errorf("Do response error:%v", err.Error())
	}

	//l.Log.Error("SERVER: Kind: %v | Value: %v | Valid: %v | Args: %#v", reflect.TypeOf(resp.Value).Kind(), resp.Value, reflect.ValueOf(resp.Value).IsValid(), args)

	var i interface{}
	if resp.Value != nil {
		err = ProtobufAnyToInterface(resp.Value, &i)
	} else {
		err = ProtobufStructToArray(resp.Arr, &i)
	}

	//l.Log.Error("SERVED: Kind: %v | Value: %v | Valid: %v", reflect.TypeOf(i).Kind(), i, reflect.ValueOf(i).IsValid())
	//l.Log.Error("SERVED: Type: %v | Value: %v", reflect.TypeOf(i), i)

	return i, resp.Error, err
}
