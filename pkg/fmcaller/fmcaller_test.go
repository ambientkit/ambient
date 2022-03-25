package fmcaller_test

import (
	"bytes"
	"errors"
	"html/template"
	"testing"
	"time"

	"github.com/ambientkit/ambient/pkg/fmcaller"
	"github.com/stretchr/testify/assert"
)

func TestMain(t *testing.T) {
	fm := make(template.FuncMap)
	fm["none_none"] = func(args ...interface{}) (interface{}, error) {
		fn := func() {}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["one_none"] = func(args ...interface{}) (interface{}, error) {
		fn := func(string) {}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["one_one"] = func(args ...interface{}) (interface{}, error) {
		fn := func(s string) string { return s }
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["onetime_one"] = func(args ...interface{}) (interface{}, error) {
		fn := func(t time.Time) string {
			return t.Format("2006-01-02")
		}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["one_error"] = func(args ...interface{}) (interface{}, error) {
		fn := func(shouldError bool) error {
			if shouldError {
				return errors.New("this is an error")
			}
			return nil
		}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["one_two"] = func(args ...interface{}) (interface{}, error) {
		fn := func(val string) (string, error) {
			return val, nil
		}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["onebool_two"] = func(args ...interface{}) (interface{}, error) {
		fn := func(shouldError bool) (string, error) {
			if shouldError {
				return "", errors.New("this is an error")
			}
			return "", nil
		}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["one_twobad"] = func(args ...interface{}) (interface{}, error) {
		fn := func(shouldError bool) (string, bool) {
			if shouldError {
				return "", shouldError
			}
			return "", shouldError
		}
		return fmcaller.CallFuncMap(fn, args...)
	}
	fm["one_threebad"] = func(args ...interface{}) (interface{}, error) {
		fn := func(shouldError bool) (string, error, error) {
			if shouldError {
				return "", errors.New("this is an error"), nil
			}
			return "", nil, nil
		}
		return fmcaller.CallFuncMap(fn, args...)
	}

	body, err := parseTemplate(t, fm, `{{none_none}}`)
	assert.NoError(t, err)
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{none_none "Foo"}}`)
	assert.Error(t, err) // input parameters expected is 0, but passed in 1
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_none "Foo"}}`)
	assert.NoError(t, err)
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_none "Foo" "Bar"}}`)
	assert.Error(t, err) // input parameters expected is 1, but passed in 2
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_one "Foo"}}`)
	assert.NoError(t, err)
	assert.Equal(t, "Foo", body)

	body, err = parseTemplate(t, fm, `{{onetime_one "2022-03-25T00:00:00Z"}}`)
	assert.NoError(t, err)
	assert.Equal(t, "2022-03-25", body)

	body, err = parseTemplate(t, fm, `{{one_error "Foo"}}`)
	assert.Error(t, err) // Call using string as type bool
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_error "false"}}`)
	assert.Error(t, err) // Call using string as type bool
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_error false}}`)
	assert.NoError(t, err)
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_error true}}`)
	assert.Error(t, err) // this is a an error
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_two "foo"}}`)
	assert.NoError(t, err)
	assert.Equal(t, "foo", body)

	body, err = parseTemplate(t, fm, `{{one_two "foo" true}}`)
	assert.Error(t, err) // input parameters expected is 1, but passed in 2
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{onebool_two "foo" true}}`)
	assert.Error(t, err) // input parameters expected is 1, but passed in 2
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{onebool_two true}}`)
	assert.Error(t, err) // this is an error
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_twobad true}}`)
	assert.Error(t, err) // this is an error
	assert.Equal(t, "", body)

	body, err = parseTemplate(t, fm, `{{one_threebad true}}`)
	assert.Error(t, err) // output parameters expected is max of 2, but passed in 3
	assert.Equal(t, "", body)
}

func parseTemplate(t *testing.T, fm template.FuncMap, content string) (string, error) {
	tmpl, err := template.New("name").Funcs(fm).Parse(content)
	assert.NoError(t, err)

	bs := bytes.NewBufferString("")
	err = tmpl.Execute(bs, nil)

	return bs.String(), err
}
