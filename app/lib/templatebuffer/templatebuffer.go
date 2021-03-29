// Package templatebuffer parses and executes templates since execute can do
// partial writes on error.
package templatebuffer

import (
	"html/template"

	"github.com/oxtoacart/bpool"
)

// bufpool is used to write out HTML after it's been executed and before it's
// written to the ResponseWriter to catch any partially written templates.
var bufpool *bpool.BufferPool = bpool.NewBufferPool(64)

// ParseTemplate will parse a template and return the string and an error.
func ParseTemplate(body string, data map[string]interface{}) (string, error) {
	// Write temporarily to a buffer pool.
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// Parse the template.
	tmpl, err := template.New("root").Funcs(template.FuncMap{
		//"OK": func() string { return "hello" },
	}).Parse(body)
	if err != nil {
		return "", err
	}

	// Execute the template.
	err = tmpl.Execute(buf, data)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
