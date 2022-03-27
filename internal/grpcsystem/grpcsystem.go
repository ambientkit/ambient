// Package grpcsystem manages connecting, loading, monitoring, and disconnecting
// of gRPC plugins.
package grpcsystem

import (
	"sync"
	"time"

	"github.com/ambientkit/ambient"
	"github.com/ambientkit/ambient/pkg/grpcp"
	"github.com/hashicorp/go-plugin"
)

//go:generate go run github.com/vburenin/ifacemaker -f *.go -s GRPCSystem -i GRPCSystem -p ambient -o ../../gen_grpcsystem.go -y "GRPCSystem manages connecting, loading, monitoring, and disconnecting gRPC plugins.." -c "Code generated by ifacemaker. DO NOT EDIT."

// GRPCSystem manages connecting, loading, monitoring, and disconnecting
// gRPC plugins.
type GRPCSystem struct {
	log          ambient.AppLogger
	pluginsystem ambient.PluginSystem
	securesite   ambient.SecureSite

	// pluginClients contains a map for quick lookup of gRPC plugins.
	pluginClients      map[string]*plugin.Client
	pluginClientsMutex sync.RWMutex
	// monitoring will be true when monitoring is turned on.
	monitoring bool
	// MonitoringFrequency is the minimum length of time between checks.
	// Defaults to 2 seconds.
	MonitoringFrequency  time.Duration
	RestartAutomatically bool
}

// New returns a new GRPCSystem.
func New(log ambient.AppLogger, pluginsystem ambient.PluginSystem) *GRPCSystem {
	pluginClients := make(map[string]*plugin.Client, 0)

	return &GRPCSystem{
		log:                  log,
		pluginsystem:         pluginsystem,
		pluginClients:        pluginClients,
		MonitoringFrequency:  2 * time.Second,
		RestartAutomatically: true,
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
func (s *GRPCSystem) ConnectAll() {
	for _, p := range s.pluginsystem.LoaderMiddleware() {
		if p.PluginVersion() == "gRPC" {
			s.Connect(p, true)
		}
	}

	for _, p := range s.pluginsystem.LoaderPlugins() {
		if p.PluginVersion() == "gRPC" {
			s.Connect(p, false)
		}
	}
}

// Connect will connect to a new gRPC plugin, these don't have to be in the
// initial plugin loader.
func (s *GRPCSystem) Connect(p ambient.Plugin, middleware bool) {
	gpb, ok := p.(*ambient.GRPCPluginBase)
	if !ok {
		s.log.Error("ambient: plugin, %v, is not a gRPC plugin", p.PluginName())
		return
	}

	gp, pc, err := grpcp.ConnectPlugin(s.log, gpb.PluginName(), gpb.PluginPath())
	if err != nil {
		s.log.Error("ambient: plugin, %v, could not establish a connection: %v", p.PluginName(), err.Error())
		return
	}

	// Load plugin - does not matter if the plugin already exists because
	// it will be overwritten. If a different type plugin already exists
	// then it will return an error.
	err = s.pluginsystem.LoadPlugin(gp, middleware, true)
	if err != nil {
		s.log.Error("ambient: plugin, %v, could not load: %v", p.PluginName(), err.Error())
		// Kill it since it can't be used.
		pc.Kill()
		return
	}

	// Store reference to the gRPC plugin.
	s.pluginClientsMutex.Lock()
	s.pluginClients[gpb.PluginName()] = pc
	s.pluginClientsMutex.Unlock()
}

// Disconnect stops the gRPC clients.
func (s *GRPCSystem) Disconnect() {
	s.setMonitor(false)
	s.pluginClientsMutex.Lock()
	for _, v := range s.pluginClients {
		v.Kill()
	}
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

	ClientLoop:
		for name, v := range s.pluginClients {
			if v.Exited() {
				s.log.Warn("ambient: detected crashed gRPC plugin: %v", name)

				var plugin ambient.Plugin
				isMiddleware := false

				for _, v := range s.pluginsystem.LoaderMiddleware() {
					if v.PluginName() == name {
						plugin = v
						isMiddleware = true
					}
				}

				if plugin == nil {
					for _, v := range s.pluginsystem.LoaderPlugins() {
						if v.PluginName() == name {
							plugin = v
						}
					}
				}

				if plugin != nil {
					if s.RestartAutomatically {
						s.Connect(plugin, isMiddleware)
						s.securesite.LoadSinglePluginPages(name)
						if isMiddleware {
							s.log.Info("ambient: gRPC middleware restarted: %v", name)
						} else {
							s.log.Info("ambient: gRPC plugin restarted: %v", name)
						}
					} else {
						s.pluginsystem.SetEnabled(name, false)
						plugin.Disable()
						delete(s.pluginClients, name)
						if isMiddleware {
							s.log.Info("ambient: gRPC middleware disabled: %v", name)
						} else {
							s.log.Info("ambient: gRPC plugin disabled: %v", name)
						}
					}

					continue ClientLoop
				}

				s.log.Error("ambient: could not find gRPC plugin: %v", name)
			}
		}
	}
}
