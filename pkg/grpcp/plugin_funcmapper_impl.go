package grpcp

import (
	"fmt"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/fmcaller"
)

// FuncMapperImpl handles the FuncMap logic.
type FuncMapperImpl struct {
	Log ambient.Logger
	Map map[string]*FMContainer
}

// Do handler.
func (d *FuncMapperImpl) Do(requestID string, key string, args []interface{}) (interface{}, error) {
	//d.Log.Warn("grpc-plugin: Do start: %v", requestID)

	c, ok := d.Map[requestID]
	if !ok {
		return nil, fmt.Errorf("could not find funcmap for request, len(%v): %v", len(d.Map), requestID)
	}

	val, ok := c.FuncMap[key]
	if !ok {
		return nil, fmt.Errorf("could not get funcmap value for key: %v", key)
	}

	//d.Log.Warn("grpc-plugin: Do start: %v %#v", requestID, args)

	anyVal, err := fmcaller.CallFuncMap(val, args...)
	if err != nil {
		return nil, err
	}

	return anyVal, nil
}
