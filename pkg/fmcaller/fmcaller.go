package fmcaller

import (
	"errors"
	"fmt"
	"reflect"
	"time"
)

// CallFuncMap calls a FuncMap function and returns an interface and an error.
func CallFuncMap(fn interface{}, args ...interface{}) (interface{}, error) {
	val := reflect.ValueOf(fn)
	t1 := reflect.TypeOf(fn)

	// Return an error if the number of input parameters doesn't match.
	if t1.NumIn() != len(args) {
		return nil, fmt.Errorf("input parameters expected is %v, but passed in %v", t1.NumIn(), len(args))
	}

	if t1.NumOut() > 2 {
		return nil, fmt.Errorf("output parameters expected is max of %v, but passed in %v", 2, t1.NumOut())
	}

	// Build a list of all the input parameters.
	arrIn := make([]reflect.Value, 0)
	for i := 0; i < t1.NumIn(); i++ {
		argInType := reflect.ValueOf(args[i])
		funcInType := t1.In(i)

		// TODO: Figure out a better way to handle this for types that are not
		// primitives.
		if argInType.Type() != funcInType {
			// Try converting to time.
			t, err := time.Parse(time.RFC3339, argInType.String())
			if err == nil {
				argInType = reflect.ValueOf(t)
			}
		}

		arrIn = append(arrIn, argInType)
	}

	// Call the function.
	out := val.Call(arrIn)

	// Return based on the return values from calling the function.
	switch len(out) {
	case 1:
		out1 := out[0].Interface()
		if out1 == nil {
			return nil, nil
		}

		// If the 1st return value is an error, then return as an error.
		if err, ok := out1.(error); ok {
			return nil, err
		}

		return out1, nil
	case 2:
		out1 := out[0].Interface()
		out2 := out[1].Interface()

		// If the 2nd return value is nil, then return nil.
		if out2 == nil {
			return out1, nil
		}

		// If the 1st return value is an error, then return as an error.
		if _, ok := out2.(error); !ok {
			return nil, errors.New("2nd return should be an error")
		}

		return out1, out2.(error)
	default:
		return nil, nil
	}
}
