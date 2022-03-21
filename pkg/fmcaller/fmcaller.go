package fmcaller

import (
	"errors"
	"fmt"
	"reflect"
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
		arrIn = append(arrIn, reflect.ValueOf(args[i]))
	}

	// Call the function.
	out := val.Call(arrIn)
	//fmt.Printf("Parameter count. In: %v | Out: %v\n", t1.NumIn(), len(out))

	// return nil, fmt.Errorf("Output: %v", out[0])
	//fmt.Println("Output:", out[0])

	// // Return an error if there are too many returns values.
	// if len(out) > 2 {
	// 	return nil, fmt.Errorf("can only have up to 2 return values, but found: %v", len(out))
	// }

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

		return out[0], nil
	case 2:
		// If the 2nd return value is nil, then return nil.
		if out[1].Interface() == nil {
			return out[0], nil
		}

		// If the 1st return value is an error, then return as an error.
		if _, ok := out[1].Interface().(error); !ok {
			return nil, errors.New("2nd return should be an error")
		}

		return out[0], out[1].Interface().(error)
	default:
		return nil, nil
	}
}
