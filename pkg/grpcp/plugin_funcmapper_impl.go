package grpcp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/fmcaller"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	"golang.org/x/net/context"
)

// FuncMapperImpl handles the FuncMap logic.
type FuncMapperImpl struct {
	Log        ambient.Logger
	Map        map[string]*FMContainer
	Impl       ambient.MiddlewarePlugin
	ContextMap map[string]context.Context
}

// Do handler.
func (d *FuncMapperImpl) Do(globalFuncMap bool, requestID string, key string, args []interface{}, method string, path string, headers http.Header, body []byte) (interface{}, error) {
	//d.Log.Warn("grpc-plugin: Do start: %v", requestID)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req = requestuuid.Set(req, requestID)
	req.Header = headers

	// Get the context if saved from middleware.
	ctx, ok := d.ContextMap[requestID]
	if ok {
		req = req.WithContext(ctx)
	}

	var callable interface{}
	if globalFuncMap {
		fmc := d.Impl.FuncMap()
		if fmc == nil {
			return nil, nil
		}

		fm := fmc(req)
		if fm == nil {
			return nil, nil
		}

		var ok bool
		callable, ok = fm[key]
		if !ok {
			//d.Log.Debug("grpc-plugin: FuncMap field not found: %v", key)
			return nil, nil
		}

		//d.Log.Debug("grpc-plugin: CallFuncMap global: %v | %v | %v | %#v", requestID, key, callable, args)
	} else {
		c, ok := d.Map[requestID]
		if !ok {
			//d.Log.Error("grpc-plugin: FuncMap not found for request ID: %v", requestID)
			return nil, nil
		}
		//d.Log.Debug("grpc-plugin: CallFuncMap local: %v | %v | %v | %#v", requestID, key, callable, args)

		if c.FuncMap == nil {
			return nil, nil
		}

		callable, ok = c.FuncMap[key]
		if !ok {
			//d.Log.Debug("grpc-plugin: FuncMap field not found: %v", key)
			return nil, nil
		}
	}

	//d.Log.Debug("grpc-plugin: CallFuncMap: %v | %v | %v | %#v", requestID, key, callable, args)

	anyVal, err := fmcaller.CallFuncMap(callable, args...)
	if err != nil {
		return nil, fmt.Errorf("grpc-plugin: CallFuncMap error: %v", err.Error())
	}

	return anyVal, nil
}
