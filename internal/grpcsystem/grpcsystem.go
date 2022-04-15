// Package grpcsystem manages connecting, loading, monitoring, and disconnecting
// of gRPC plugins.
package grpcsystem

import (
	"context"
	"sync"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp"
	hplugin "github.com/hashicorp/go-plugin"
)

//go:generate go run github.com/vburenin/ifacemaker -f *.go -s GRPCSystem -i GRPCSystem -p ambient -o ../../gen_grpcsystem.go -y "GRPCSystem manages connecting, loading, monitoring, and disconnecting gRPC plugins.." -c "Code generated by ifacemaker. DO NOT EDIT."

// GRPCSystem manages connecting, loading, monitoring, and disconnecting
// gRPC plugins.
type GRPCSystem struct {
	log          ambient.AppLogger
	pluginsystem ambient.PluginSystem
	securesite   ambient.SecureSite

	// pluginClients contains a map for quick lookup of gRPC plugins.
	pluginClients         map[string]*hplugin.Client
	pluginClientsProtocol map[string]hplugin.ClientProtocol
	pluginClientsMutex    sync.RWMutex
	// monitoring will be true when monitoring is turned on.
	monitoring bool
	// MonitoringFrequency is the minimum length of time between checks.
	// Defaults to 2 seconds.
	MonitoringFrequency  time.Duration
	RestartAutomatically bool
}

// New returns a new GRPCSystem.
func New(log ambient.AppLogger, pluginsystem ambient.PluginSystem) *GRPCSystem {
	return &GRPCSystem{
		log:                   log,
		pluginsystem:          pluginsystem,
		pluginClients:         make(map[string]*hplugin.Client),
		pluginClientsProtocol: make(map[string]hplugin.ClientProtocol),
		MonitoringFrequency:   2 * time.Second,
		RestartAutomatically:  true,
	}
}

// Monitor starts monitoring the gRPC plugins.
func (s *GRPCSystem) Monitor(securesite ambient.SecureSite) {
	s.securesite = securesite
	if !s.isMonitoring() && len(s.pluginClients) > 0 {
		go s.monitorGRPCClients()
	}
}

// ConnectAll will connect to all initial gRPC plugins in the plugin system.
func (s *GRPCSystem) ConnectAll(ctx context.Context) {
	for _, p := range s.pluginsystem.LoaderMiddleware() {
		if p.PluginVersion(ctx) == "gRPC" {
			s.Connect(ctx, p, true)
		}
	}

	for _, p := range s.pluginsystem.LoaderPlugins() {
		if p.PluginVersion(ctx) == "gRPC" {
			s.Connect(ctx, p, false)
		}
	}
}

// Connect will connect to a new gRPC plugin, these don't have to be in the
// initial plugin loader.
func (s *GRPCSystem) Connect(ctx context.Context, p ambient.Plugin, middleware bool) {
	name := p.PluginName(ctx)
	gpb, ok := p.(*ambient.GRPCPluginBase)
	if !ok {
		s.log.Error("plugin, %v, is not a gRPC plugin", name)
		return
	}

	gp, pc, cp, err := grpcp.ConnectPlugin(s.log.Named("grpc-server"), name, gpb.PluginPath())
	if err != nil {
		s.log.Error("plugin, %v, could not establish a connection: %v", name, err.Error())
		return
	}

	// Load plugin - does not matter if the plugin already exists because
	// it will be overwritten. If a different type plugin already exists
	// then it will return an error.
	err = s.pluginsystem.LoadPlugin(ctx, gp, middleware, true)
	if err != nil {
		s.log.Error("plugin, %v, could not load: %v", name, err.Error())
		// Kill it since it can't be used.
		pc.Kill()
		cp.Close()
		return
	}

	// Store reference to the gRPC plugin.
	s.pluginClientsMutex.Lock()
	s.pluginClients[name] = pc
	s.pluginClientsProtocol[name] = cp
	s.pluginClientsMutex.Unlock()
}

// Disconnect stops the gRPC clients.
func (s *GRPCSystem) Disconnect() {
	s.setMonitor(false)
	s.pluginClientsMutex.Lock()
	for _, v := range s.pluginClientsProtocol {
		v.Close()
		//v.Kill()
	}
	hplugin.CleanupClients()
	s.pluginClientsMutex.Unlock()
}

func (s *GRPCSystem) isMonitoring() bool {
	s.pluginClientsMutex.RLock()
	b := s.monitoring
	s.pluginClientsMutex.RUnlock()
	return b
}

func (s *GRPCSystem) setMonitor(val bool) {
	s.pluginClientsMutex.Lock()
	s.monitoring = val
	s.pluginClientsMutex.Unlock()
}

// monitorGRPCClients will restart clients if they crash.
func (s *GRPCSystem) monitorGRPCClients() {
	if s.isMonitoring() {
		return
	}
	s.setMonitor(true)

	for s.isMonitoring() {
		<-time.After(s.MonitoringFrequency)
		// Break if monitoring changes while waiting.
		if !s.isMonitoring() {
			break
		}
		ctx := context.Background()
	ClientLoop:
		for name, v := range s.pluginClients {
			if v.Exited() {
				s.log.Warn("detected crashed gRPC plugin: %v", name)

				err := s.pluginClientsProtocol[name].Close()
				if err != nil {
					s.log.Error("error closing gRPC connection: %v", err.Error())
				}
				//v.Kill()

				var plugin ambient.Plugin
				isMiddleware := false

				for _, v := range s.pluginsystem.LoaderMiddleware() {
					if v.PluginName(ctx) == name {
						plugin = v
						isMiddleware = true
					}
				}

				if plugin == nil {
					for _, v := range s.pluginsystem.LoaderPlugins() {
						if v.PluginName(ctx) == name {
							plugin = v
						}
					}
				}

				if plugin != nil {
					delete(s.pluginClients, name)
					if s.RestartAutomatically {
						s.Connect(ctx, plugin, isMiddleware)
						s.securesite.LoadSinglePluginPages(name)
						if isMiddleware {
							s.log.Info("gRPC middleware restarted: %v", name)
						} else {
							s.log.Info("gRPC plugin restarted: %v", name)
						}
					} else {
						s.pluginsystem.SetEnabled(name, false)
						plugin.Disable()
						if isMiddleware {
							s.log.Info("gRPC middleware disabled: %v", name)
						} else {
							s.log.Info("gRPC plugin disabled: %v", name)
						}
					}

					continue ClientLoop
				}

				s.log.Error("could not find gRPC plugin: %v", name)
			}
		}
	}
}
