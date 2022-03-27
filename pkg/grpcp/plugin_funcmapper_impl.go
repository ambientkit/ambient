package grpcp

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/fmcaller"
	"github.com/ambientkit/ambient/pkg/requestuuid"
)

// FuncMapperImpl handles the FuncMap logic.
type FuncMapperImpl struct {
	Log  ambient.Logger
	Map  map[string]*FMContainer
	Impl ambient.MiddlewarePlugin
}

// Do handler.
func (d *FuncMapperImpl) Do(requestID string, key string, args []interface{}, method string, path string, headers http.Header, body []byte) (interface{}, error) {
	//d.Log.Warn("grpc-plugin: Do start: %v", requestID)

	// FIXME: May not want to use the map here since it may not be set up in
	// another call. Probably better to construct the http.Request

	req := httptest.NewRequest(method, path, bytes.NewReader(body))
	req = requestuuid.Set(req, requestID)
	req.Header = headers

	// FIXME: This is ignoring what is passed in - this will behave differently
	// from the global funcmap.
	fmc := d.Impl.FuncMap()
	fm := fmc(req)
	val := fm[key]

	//d.Log.Warn("grpc-plugin: CallFuncMap: %v | %v | %v | %#v", requestID, key, val, args)

	anyVal, err := fmcaller.CallFuncMap(val, args...)
	if err != nil {
		return nil, fmt.Errorf("grpc-plugin: CallFuncMap error: %v", err.Error())
	}

	return anyVal, nil
}
