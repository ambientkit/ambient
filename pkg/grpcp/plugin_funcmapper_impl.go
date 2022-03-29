package grpcp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/fmcaller"
	"github.com/ambientkit/ambient/pkg/grpcp/grpcsafe"
	"github.com/ambientkit/ambient/pkg/requestuuid"
)

// FuncMapperImpl handles the FuncMap logic.
type FuncMapperImpl struct {
	Log         ambient.Logger
	Impl        ambient.MiddlewarePlugin
	PluginState *grpcsafe.PluginState
}

// Do handler.
func (d *FuncMapperImpl) Do(globalFuncMap bool, requestID string, key string, args []interface{}, method string, path string, headers http.Header, body []byte) (interface{}, string, error) {
	//d.Log.Warn("Do start: %v", requestID)

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req = requestuuid.Set(req, requestID)
	req.Header = headers

	// Get the context if saved from middleware.
	ctx, ok := d.PluginState.Context(requestID)
	if ok {
		req = req.WithContext(ctx)
	}

	var callable interface{}
	if globalFuncMap {
		fmc := d.Impl.FuncMap()
		if fmc == nil {
			return nil, "", nil
		}

		fm := fmc(req)
		if fm == nil {
			return nil, "", nil
		}

		var ok bool
		callable, ok = fm[key]
		if !ok {
			//d.Log.Debug("FuncMap field not found: %v", key)
			return nil, "", nil
		}

		//d.Log.Debug("CallFuncMap global: %v | %v | %v | %#v", requestID, key, callable, args)
	} else {
		c, ok := d.PluginState.Assets(requestID)
		if !ok {
			//d.Log.Error("FuncMap not found for request ID: %v", requestID)
			return nil, "", nil
		}
		//d.Log.Debug("CallFuncMap local: %v | %v | %v | %#v", requestID, key, callable, args)

		if c.FuncMap == nil {
			return nil, "", nil
		}

		callable, ok = c.FuncMap[key]
		if !ok {
			//d.Log.Debug("FuncMap field not found: %v", key)
			return nil, "", nil
		}
	}

	//d.Log.Debug("CallFuncMap: %v | %v | %v | %#v", requestID, key, callable, args)

	anyVal, err := fmcaller.CallFuncMap(callable, args...)
	if err != nil {
		_, ok := err.(fmcaller.FMError)
		if !ok {
			return nil, "", fmt.Errorf("CallFuncMap error: %v", err.Error())
		}
		return nil, err.Error(), nil
	}

	se, ok := anyVal.(fmcaller.FMError)
	if ok {
		return se.Error(), "", nil
	}

	return anyVal, "", nil
}
