package grpcp

import (
	"context"
	"embed"
	"html/template"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
)

// GRPCRendererPlugin .
type GRPCRendererPlugin struct {
	client protodef.RendererClient
	Log    ambient.Logger
	Map    map[string]*FMContainer
}

// FMContainer .
type FMContainer struct {
	FuncMap template.FuncMap
	FS      *embed.FS
}

// Page handler.
func (l *GRPCRendererPlugin) Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	pvars, err := MapToProtobufStruct(vars)
	if err != nil {
		return err
	}

	funcMap := fm(nil)
	keys := make([]string, 0)
	for k := range funcMap {
		keys = append(keys, k)
	}

	rid := requestID(r)
	l.Map[rid] = &FMContainer{
		FuncMap: fm(r),
		FS:      &assets,
	}
	defer delete(l.Map, rid)

	_, err = l.client.Page(context.Background(), &protodef.RendererPageRequest{
		Requestid:    rid,
		Templatename: templateName,
		Vars:         pvars,
		Keys:         keys,
	})

	return err
}

// PageContent handler.
func (l *GRPCRendererPlugin) PageContent(w http.ResponseWriter, r *http.Request, content string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	pvars, err := MapToProtobufStruct(vars)
	if err != nil {
		return err
	}

	funcMap := fm(nil)
	keys := make([]string, 0)
	for k := range funcMap {
		keys = append(keys, k)
	}

	rid := requestID(r)
	l.Map[rid] = &FMContainer{
		FuncMap: fm(r),
		FS:      nil,
	}
	defer delete(l.Map, rid)

	_, err = l.client.PageContent(context.Background(), &protodef.RendererPageContentRequest{
		Requestid: rid,
		Content:   content,
		Vars:      pvars,
		Keys:      keys,
	})

	return err
}

// Post handler.
func (l *GRPCRendererPlugin) Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	pvars, err := MapToProtobufStruct(vars)
	if err != nil {
		return err
	}

	funcMap := fm(nil)
	keys := make([]string, 0)
	for k := range funcMap {
		keys = append(keys, k)
	}

	rid := requestID(r)
	l.Map[rid] = &FMContainer{
		FuncMap: fm(r),
		FS:      &assets,
	}
	defer delete(l.Map, rid)

	_, err = l.client.Post(context.Background(), &protodef.RendererPostRequest{
		Requestid:    rid,
		Templatename: templateName,
		Vars:         pvars,
		Keys:         keys,
	})

	return err
}

// PostContent handler.
func (l *GRPCRendererPlugin) PostContent(w http.ResponseWriter, r *http.Request, content string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	pvars, err := MapToProtobufStruct(vars)
	if err != nil {
		return err
	}

	funcMap := fm(nil)
	keys := make([]string, 0)
	for k := range funcMap {
		keys = append(keys, k)
	}

	rid := requestID(r)
	l.Map[rid] = &FMContainer{
		FuncMap: fm(r),
		FS:      nil,
	}
	defer delete(l.Map, rid)

	_, err = l.client.PostContent(context.Background(), &protodef.RendererPostContentRequest{
		Requestid: rid,
		Content:   content,
		Vars:      pvars,
		Keys:      keys,
	})

	return err
}

// Error handler.
func (l *GRPCRendererPlugin) Error(w http.ResponseWriter, r *http.Request, content string, statusCode int,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	pvars, err := MapToProtobufStruct(vars)
	if err != nil {
		return err
	}

	funcMap := fm(nil)
	keys := make([]string, 0)
	for k := range funcMap {
		keys = append(keys, k)
	}

	rid := requestID(r)
	l.Map[rid] = &FMContainer{
		FuncMap: fm(r),
		FS:      nil,
	}
	defer delete(l.Map, rid)

	_, err = l.client.Error(context.Background(), &protodef.RendererErrorRequest{
		Requestid:  rid,
		Content:    content,
		Vars:       pvars,
		Keys:       keys,
		Statuscode: uint32(statusCode),
	})

	return err
}
