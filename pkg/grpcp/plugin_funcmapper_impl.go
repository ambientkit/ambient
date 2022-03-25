package grpcp

import (
	"fmt"
	"net/http"

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
func (d *FuncMapperImpl) Do(requestID string, key string, args []interface{}) (interface{}, error) {
	//d.Log.Warn("grpc-plugin: Do start: %v", requestID)

	// c, ok := d.Map[requestID]
	// if !ok {
	// 	return nil, fmt.Errorf("could not find funcmap for request, len(%v): %v", len(d.Map), requestID)
	// }

	// val, ok := c.FuncMap[key]
	// if !ok {
	// 	return nil, fmt.Errorf("could not get funcmap value for key: %v", key)
	// }

	// FIXME: May not want to use the map here since it may not be set up in
	// another call. Probably better to construct the http.Request

	req, _ := http.NewRequest("GET", "/", nil)
	req = requestuuid.Set(req, requestID)
	//req.Header = headers
	//w := NewResponseWriter()

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
