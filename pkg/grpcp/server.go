package grpcp

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/avfs"
	"github.com/ambientkit/ambient/pkg/grpcp/grpcsafe"
	"github.com/ambientkit/ambient/pkg/grpcp/protodef"
	"github.com/ambientkit/ambient/pkg/requestuuid"
	plugin "github.com/hashicorp/go-plugin"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// GRPCServer is the server side implementation.
type GRPCServer struct {
	broker           *plugin.GRPCBroker
	client           protodef.GenericPluginClient
	toolkit          *ambient.Toolkit
	conn             *grpc.ClientConn
	server           *grpc.Server
	serverState      *grpcsafe.ServerState
	funcMapperClient *GRPCFuncMapperServer
}

// PluginName handler.
func (m *GRPCServer) PluginName() string {
	resp, err := m.client.PluginName(context.Background(), &protodef.Empty{})
	if err != nil {
		return ""
	}

	return resp.Name
}

// PluginVersion handler.
func (m *GRPCServer) PluginVersion() string {
	resp, err := m.client.PluginVersion(context.Background(), &protodef.Empty{})
	if err != nil {
		return ""
	}

	return resp.Version
}

// Enable handler.
func (m *GRPCServer) Enable(toolkit *ambient.Toolkit) error {
	//toolkit.Log.Debug("enabled called")

	m.toolkit = toolkit
	loggerServer := &GRPCLoggerServer{
		Impl: toolkit.Log,
	}
	routerServer := &GRPCAddRouterServer{
		Impl:   toolkit.Mux,
		Log:    toolkit.Log,
		broker: m.broker,
		reqmap: m.serverState,
	}
	siteServer := &GRPCSiteServer{
		Impl:   toolkit.Site,
		Log:    toolkit.Log,
		reqmap: m.serverState,
	}
	rendererServer := &GRPCRendererServer{
		Log:    toolkit.Log,
		Impl:   toolkit.Render,
		reqmap: m.serverState,
	}
	funcMapperServer := &GRPCFuncMapperPlugin{
		Impl: &FuncMapperImpl{
			Log: m.toolkit.Log,
		},
		Log: m.toolkit.Log,
	}

	serverFunc := func(opts []grpc.ServerOption) *grpc.Server {
		m.server = grpc.NewServer(opts...)
		protodef.RegisterLoggerServer(m.server, loggerServer)
		protodef.RegisterRouterServer(m.server, routerServer)
		protodef.RegisterSiteServer(m.server, siteServer)
		protodef.RegisterRendererServer(m.server, rendererServer)
		protodef.RegisterFuncMapperServer(m.server, funcMapperServer)

		return m.server
	}

	brokerID := m.broker.NextId()
	go m.broker.AcceptAndServe(brokerID, serverFunc)

	resp, err := m.client.Enable(context.Background(), &protodef.Toolkit{
		Uid: brokerID,
	})
	if err != nil {
		return err
	}

	m.conn, err = m.broker.Dial(resp.Uid)
	if err != nil {
		return err
	}

	routerServer.HandlerClient = &GRPCHandlerServer{
		client: protodef.NewHandlerClient(m.conn),
	}

	rendererServer.FuncMapperClient = &GRPCFuncMapperServer{
		client: protodef.NewFuncMapperClient(m.conn),
		Log:    toolkit.Log,
	}
	m.funcMapperClient = rendererServer.FuncMapperClient

	return nil
}

// Disable handler.
func (m *GRPCServer) Disable() error {
	if m.server != nil {
		_, _ = m.client.Disable(context.Background(), &protodef.Empty{})
		m.server.Stop()
		m.server = nil
	}
	return nil
}

// Routes handler.
func (m *GRPCServer) Routes() {
	if m.server == nil || m.toolkit == nil || m.toolkit.Log == nil {
		return
	}

	_, err := m.client.Routes(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("error calling routes: %v", err)
	}
}

// Assets handler.
func (m *GRPCServer) Assets() ([]ambient.Asset, ambient.FileSystemReader) {
	resp, err := m.client.Assets(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("error calling Assets: %v", err)
		return nil, nil
	}

	var assets []ambient.Asset
	err = ProtobufStructToArray(resp.Assets, &assets)
	if err != nil {
		m.toolkit.Log.Error("error calling Assets conversion: %v", err)
	}

	// Build a file system.
	efs := avfs.NewFS()
	for _, v := range resp.Files {
		efs.AddFile(v.Name, v.Body)
	}

	return assets, efs
}

// Settings handler.
func (m *GRPCServer) Settings() []ambient.Setting {
	resp, err := m.client.Settings(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("error calling Settings: %v", err)
		return nil
	}

	arr := make([]ambient.Setting, 0)

	for _, v := range resp.Settings {
		var i interface{}
		err = ProtobufAnyToInterface(v.Default, &i)
		if err != nil {
			m.toolkit.Log.Error("error calling Settings: %v", err)
		}

		arr = append(arr, ambient.Setting{
			Name: v.Name,
			Type: ambient.SettingType(v.Settingtype),
			Description: ambient.SettingDescription{
				Text: v.Description.Text,
				URL:  v.Description.Url,
			},
			Hide:    v.Hide,
			Default: i,
		})
	}

	if err != nil {
		m.toolkit.Log.Error("error calling Settings conversion: %v", err)
	}
	return arr
}

// GrantRequests handler.
func (m *GRPCServer) GrantRequests() []ambient.GrantRequest {
	resp, err := m.client.GrantRequests(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("error calling GrantRequests: %v", err)
		return nil
	}

	arr := make([]ambient.GrantRequest, 0)
	for _, v := range resp.Grantrequest {
		arr = append(arr, ambient.GrantRequest{
			Grant:       ambient.Grant(v.Grant),
			Description: v.Description,
		})
	}

	return arr
}

// FuncMap handler.
func (m *GRPCServer) FuncMap() func(r *http.Request) template.FuncMap {
	// Return a list of keys for the FuncMap().
	resp, err := m.client.FuncMap(context.Background(), &protodef.Empty{})
	if err != nil {
		m.toolkit.Log.Error("error calling FuncMap: %v", err)
		return nil
	}

	if len(resp.Keys) == 0 {
		return nil
	}

	return func(req *http.Request) template.FuncMap {
		fm := make(template.FuncMap)
		for _, rawV := range resp.Keys {
			// Prevent race conditions.
			v := rawV
			fm[v] = func(args ...interface{}) (interface{}, error) {
				val, errText, err := m.funcMapperClient.Do(req, requestuuid.Get(req), v, args, true)
				if err != nil {
					m.toolkit.Log.Error("error executing FuncMap: %v", err)
				}
				if len(errText) > 0 {
					return val, errors.New(errText)
				}
				return val, nil
			}
		}

		return fm
	}
}

// Middleware handler.
func (m *GRPCServer) Middleware() []func(next http.Handler) http.Handler {
	return []func(next http.Handler) http.Handler{
		func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sm, err := ObjectToProtobufStruct(r.Header)
				if err != nil {
					m.toolkit.Log.Error("error getting Middleware header: %v", err)
				}

				body := bytes.NewBuffer(nil)
				_, err = io.Copy(body, r.Body)
				if err != nil {
					m.toolkit.Log.Error("error getting Middleware body: %v", err)
				}
				// Restore body.
				r.Body = ioutil.NopCloser(body)

				//m.toolkit.Log.Error("body in: %v | %v | %v", r.RequestURI, len(body.Bytes()), body.String())

				uuid := requestuuid.Get(r)
				m.serverState.Save(uuid, &grpcsafe.HTTPContainer{
					Request:  r,
					Response: w,
					FuncMap:  make(template.FuncMap),
				})
				defer m.serverState.Delete(uuid)

				resp, err := m.client.Middleware(context.Background(), &protodef.MiddlewareRequest{
					Requestid: uuid,
					Method:    r.Method,
					Path:      r.RequestURI,
					Headers:   sm,
					Body:      body.Bytes(),
				})

				if err != nil {
					m.toolkit.Log.Error("error calling Middleware: %v", err)
					return
				}

				if resp.Status != 0 {
					var outHeader http.Header
					err = ProtobufStructToObject(resp.Headers, &outHeader)
					if err != nil {
						m.toolkit.Log.Error("error converting Middleware headers: %v", err)
					}

					// Copy over the headers.
					for k, v := range outHeader {
						w.Header()[k] = v
					}

					//m.toolkit.Log.Error("body ou: %v | %v", len(resp.Response), string(resp.Response))

					// If the response has text, then display it.
					if len(resp.Response) > 0 {
						w.WriteHeader(int(resp.Status))
						fmt.Fprint(w, resp.Response)
						return
					}

					// If the status came back as something other than 200, then display it.
					if resp.Status != http.StatusOK {
						w.WriteHeader(int(resp.Status))
						if len(resp.Error) > 0 {
							fmt.Fprint(w, resp.Error)
						}
						return
					}
				}

				h.ServeHTTP(w, r)
			})
		},
	}
}
