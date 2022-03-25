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
		argInValue := reflect.ValueOf(args[i])
		funcInType := t1.In(i)

		// If nil, then append and continue. Don't need to do any more work.
		// When 'nil' is passed in the array:
		// Type: <nil>
		// Value: <invalid reflect.Value>
		if !argInValue.IsValid() {
			arrIn = append(arrIn, argInValue)
			continue

			// An example when you pass in nil when it requires a value between
			// this package and the html/template package:
			// This: template: name:1:2: executing "name" at <onetime_one nil>: error calling onetime_one: reflect: call of reflect.Value.Type on zero Value
			// Real: template: name:1:15: executing "name" at <nil>: cannot assign nil to time.Time
		}

		// If someone passed in something unassignable, then try to convert it.
		if !argInValue.Type().AssignableTo(funcInType) {
			switch funcInType {
			case reflect.TypeOf(time.Time{}):
				// Try converting to time.
				t, err := time.Parse(time.RFC3339, argInValue.String())
				if err == nil {
					argInValue = reflect.ValueOf(t)
				}
			default:
				return nil, fmt.Errorf("error with FuncMap: cannot assign %v to %v", argInValue.Type(), funcInType)
			}
		}

		// I could always pass in the zero value if something is not right, but
		// that doesn't seem correct. Should respond the same way as the native
		// template package - with an error.
		// arrIn = append(arrIn, reflect.Zero(funcInType))

		if argInValue.Type() != funcInType {
			// Try converting to time.
			t, err := time.Parse(time.RFC3339, argInValue.String())
			if err == nil {
				argInValue = reflect.ValueOf(t)
			}
		}

		arrIn = append(arrIn, argInValue)
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
