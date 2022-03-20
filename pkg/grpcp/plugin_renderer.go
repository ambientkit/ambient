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
}

// Page handler.
func (l *GRPCRendererPlugin) Page(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	l.Log.Error("grpc-plugin: Page1 hit!")
	// _, err = l.client.Page(context.Background(), &protodef.RendererPageRequest{
	// 	Requestid:    requestID(r),
	// 	Templatename: templateName,
	// })
	return err
}

// PageContent handler.
func (l *GRPCRendererPlugin) PageContent(w http.ResponseWriter, r *http.Request, content string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	l.Log.Error("grpc-plugin: Page2 hit!")

	pvars, err := MapToProtobufStruct(vars)
	if err != nil {
		return err
	}

	_, err = l.client.PageContent(context.Background(), &protodef.RendererPageContentRequest{
		Requestid: requestID(r),
		Content:   content,
		Vars:      pvars,
	})
	return err
}

// Post handler.
func (l *GRPCRendererPlugin) Post(w http.ResponseWriter, r *http.Request, assets embed.FS, templateName string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	l.Log.Error("grpc-plugin: Page3 hit!")
	return nil
}

// PostContent handler.
func (l *GRPCRendererPlugin) PostContent(w http.ResponseWriter, r *http.Request, content string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	l.Log.Error("grpc-plugin: Page4 hit!")
	return nil
}

// Error handler.
func (l *GRPCRendererPlugin) Error(w http.ResponseWriter, r *http.Request, content string, statusCode int,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	l.Log.Error("grpc-plugin: Page5 hit!")
	return nil
}
