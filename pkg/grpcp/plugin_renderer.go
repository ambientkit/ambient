package grpcp

import (
	"fmt"
	"html/template"
	"io/fs"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp/grpcsafe"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
)

// GRPCRendererPlugin .
type GRPCRendererPlugin struct {
	client      protodef.RendererClient
	Log         ambient.Logger
	PluginState *grpcsafe.PluginState
}

// Page handler.
func (l *GRPCRendererPlugin) Page(w http.ResponseWriter, r *http.Request, assets ambient.FileSystemReader, templateName string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	if r == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.Request cannot be nil")}
	} else if w == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.ResponseWriter cannot be nil: %v", r.RequestURI)}
	}

	pvars, err := ObjectToProtobufStruct(vars)
	if err != nil {
		return fmt.Errorf("error on Page struct conversion: %v | %v", err.Error(), pvars)
	}

	keys := make([]string, 0)

	c := &grpcsafe.AssetContainer{
		FS: assets,
	}

	if fm != nil {
		funcMap := fm(r)
		for k := range funcMap {
			keys = append(keys, k)
		}

		c.FuncMap = funcMap
	}

	rid := requestuuid.Get(r)
	l.PluginState.SaveAssets(c, rid)
	defer l.PluginState.DeleteAssets(rid)

	// Save and remove context after 30 seconds.
	_, ok := l.PluginState.Context(rid)
	l.PluginState.SaveContext(r.Context(), rid)
	if !ok {
		// If doesn't exist, then the middleware didn't add it so it needs to be deleted
		// after.
		l.PluginState.DeleteContextDelayed(rid)
	}

	files := make([]*protodef.EmbeddedFile, 0)

	if assets != nil {
		err = fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			b, err := assets.ReadFile(path)
			if err != nil {
				return err
			}

			files = append(files, &protodef.EmbeddedFile{
				Name: path,
				Body: b,
			})
			return nil
		})
		if err != nil {
			return err
		}
	}

	_, err = l.client.Page(r.Context(), &protodef.RendererPageRequest{
		Requestid:    rid,
		Templatename: templateName,
		Vars:         pvars,
		Keys:         keys,
		Files:        files,
	})

	return err
}

// PageContent handler.
func (l *GRPCRendererPlugin) PageContent(w http.ResponseWriter, r *http.Request, content string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	if r == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.Request cannot be nil")}
	} else if w == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.ResponseWriter cannot be nil: %v", r.RequestURI)}
	}

	pvars, err := ObjectToProtobufStruct(vars)
	if err != nil {
		return fmt.Errorf("error on PageContent struct conversion: %v", err.Error())
	}

	keys := make([]string, 0)
	c := &grpcsafe.AssetContainer{}

	if fm != nil {
		funcMap := fm(r)
		for k := range funcMap {
			keys = append(keys, k)
		}

		c.FuncMap = funcMap
	}

	rid := requestuuid.Get(r)
	l.PluginState.SaveAssets(c, rid)
	defer l.PluginState.DeleteAssets(rid)

	// Save and remove context after 30 seconds.
	_, ok := l.PluginState.Context(rid)
	l.PluginState.SaveContext(r.Context(), rid)
	if !ok {
		// If doesn't exist, then the middleware didn't add it so it needs to be deleted
		// after.
		l.PluginState.DeleteContextDelayed(rid)
	}

	_, err = l.client.PageContent(r.Context(), &protodef.RendererPageContentRequest{
		Requestid: rid,
		Content:   content,
		Vars:      pvars,
		Keys:      keys,
	})

	return err
}

// Post handler.
func (l *GRPCRendererPlugin) Post(w http.ResponseWriter, r *http.Request, assets ambient.FileSystemReader, templateName string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	if r == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.Request cannot be nil")}
	} else if w == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.ResponseWriter cannot be nil: %v", r.RequestURI)}
	}

	pvars, err := ObjectToProtobufStruct(vars)
	if err != nil {
		return fmt.Errorf("error on Post struct conversion: %v", err.Error())
	}

	keys := make([]string, 0)

	c := &grpcsafe.AssetContainer{
		FS: assets,
	}

	if fm != nil {
		funcMap := fm(r)
		for k := range funcMap {
			keys = append(keys, k)
		}

		c.FuncMap = funcMap
	}

	rid := requestuuid.Get(r)
	l.PluginState.SaveAssets(c, rid)
	defer l.PluginState.DeleteAssets(rid)

	// Save and remove context after 30 seconds.
	_, ok := l.PluginState.Context(rid)
	l.PluginState.SaveContext(r.Context(), rid)
	if !ok {
		// If doesn't exist, then the middleware didn't add it so it needs to be deleted
		// after.
		l.PluginState.DeleteContextDelayed(rid)
	}

	files := make([]*protodef.EmbeddedFile, 0)

	if assets != nil {
		err = fs.WalkDir(assets, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				return nil
			}

			b, err := assets.ReadFile(path)
			if err != nil {
				return err
			}

			files = append(files, &protodef.EmbeddedFile{
				Name: path,
				Body: b,
			})
			return nil
		})
		if err != nil {
			return err
		}
	}

	_, err = l.client.Post(r.Context(), &protodef.RendererPostRequest{
		Requestid:    rid,
		Templatename: templateName,
		Vars:         pvars,
		Keys:         keys,
		Files:        files,
	})

	return err
}

// PostContent handler.
func (l *GRPCRendererPlugin) PostContent(w http.ResponseWriter, r *http.Request, content string,
	fm func(r *http.Request) template.FuncMap, vars map[string]interface{}) (err error) {
	if r == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.Request cannot be nil")}
	} else if w == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.ResponseWriter cannot be nil: %v", r.RequestURI)}
	}

	pvars, err := ObjectToProtobufStruct(vars)
	if err != nil {
		return fmt.Errorf("error on PostContent struct conversion: %v", err.Error())
	}

	keys := make([]string, 0)
	c := &grpcsafe.AssetContainer{}

	if fm != nil {
		funcMap := fm(r)
		for k := range funcMap {
			keys = append(keys, k)
		}

		c.FuncMap = funcMap
	}

	rid := requestuuid.Get(r)
	l.PluginState.SaveAssets(c, rid)
	defer l.PluginState.DeleteAssets(rid)

	// Save and remove context after 30 seconds.
	_, ok := l.PluginState.Context(rid)
	l.PluginState.SaveContext(r.Context(), rid)
	if !ok {
		// If doesn't exist, then the middleware didn't add it so it needs to be deleted
		// after.
		l.PluginState.DeleteContextDelayed(rid)
	}

	_, err = l.client.PostContent(r.Context(), &protodef.RendererPostContentRequest{
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
	if r == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.Request cannot be nil")}
	} else if w == nil {
		return ambient.StatusError{Code: http.StatusInternalServerError, Err: fmt.Errorf("htmlengine: http.ResponseWriter cannot be nil: %v", r.RequestURI)}
	}

	pvars, err := ObjectToProtobufStruct(vars)
	if err != nil {
		return fmt.Errorf("error on Error struct conversion: %v", err.Error())
	}

	keys := make([]string, 0)
	c := &grpcsafe.AssetContainer{}

	if fm != nil {
		funcMap := fm(r)
		for k := range funcMap {
			keys = append(keys, k)
		}

		c.FuncMap = funcMap
	}

	rid := requestuuid.Get(r)
	l.PluginState.SaveAssets(c, rid)
	defer l.PluginState.DeleteAssets(rid)

	// Save and remove context after 30 seconds.
	_, ok := l.PluginState.Context(rid)
	l.PluginState.SaveContext(r.Context(), rid)
	if !ok {
		// If doesn't exist, then the middleware didn't add it so it needs to be deleted
		// after.
		l.PluginState.DeleteContextDelayed(rid)
	}

	_, err = l.client.Error(r.Context(), &protodef.RendererErrorRequest{
		Requestid:  rid,
		Content:    content,
		Vars:       pvars,
		Keys:       keys,
		Statuscode: uint32(statusCode),
	})

	return err
}
