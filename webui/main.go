package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"rogchap.com/v8go"
)

func init() {
	// Verbose logging with file name and line number.
	log.SetFlags(log.Lshortfile)
}

func main() {
	r, err := NewServerSideRenderer("./dist/static/js/index.js")
	if err != nil {
		log.Fatalln(err.Error())
	}
	output, err := r.Render()
	if err != nil {
		log.Fatalln(err.Error())
	}

	fmt.Println(output)
}

// Renderer renders a React application to HTML.
type Renderer struct {
	Path string

	scriptSource string

	isolate        *v8go.Isolate
	global         *v8go.ObjectTemplate
	reactssrObject *v8go.ObjectTemplate
}

// NewServerSideRenderer creates a new server side renderer from a JavaScript bundle file.
func NewServerSideRenderer(path string) (*Renderer, error) {
	// TODO: perform validation(s) on path.

	iso, err := v8go.NewIsolate()
	if err != nil {
		return nil, err
	}
	// The "global" global injected into the v8 isolate's global namespace.
	global, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}
	// The "reactssr" global injected into the v8 isolate's global namespace.
	reactssrObj, err := v8go.NewObjectTemplate(iso)
	if err != nil {
		return nil, err
	}
	r := &Renderer{
		Path: path,

		isolate:        iso,
		global:         global,
		reactssrObject: reactssrObj,
	}
	return r, r.loadScriptSource()
}

// Render renders the provided path to HTML.
func (r *Renderer) Render() (string, error) {
	return r.render()
}

func (r *Renderer) loadScriptSource() error {
	buf, err := ioutil.ReadFile(r.Path)
	if err != nil {
		return err
	}
	r.scriptSource = string(buf)
	return nil
}

// render renders the CRA output bundle.
func (r *Renderer) render() (string, error) {
	reactssr, err := v8go.NewObjectTemplate(r.isolate)
	if err != nil {
		return "", err
	}
	outputHTML := ""
	render, err := v8go.NewFunctionTemplate(r.isolate, func(info *v8go.FunctionCallbackInfo) *v8go.Value {
		args := info.Args()
		if len(args) > 0 {
			outputHTML = args[0].String()
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	reactssr.Set("render", render)
	r.global.Set("reactssr", reactssr)
	ctx, err := v8go.NewContext(r.isolate, r.global)
	if err != nil {
		return "", err
	}
	if _, err := ctx.RunScript(r.scriptSource, r.Path); err != nil {
		return "", err
	}
	return outputHTML, nil
}
