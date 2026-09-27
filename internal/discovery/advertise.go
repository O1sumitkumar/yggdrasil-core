package discovery

import (
	"fmt"

	"github.com/grandcat/zeroconf"
	"github.com/yeixio/yggdrasil-core/internal/config"
	"github.com/yeixio/yggdrasil-core/internal/version"
)

// Advertiser publishes mDNS service records.
type Advertiser struct {
	server *zeroconf.Server
}

// StartAdvertise registers _localai._tcp with node metadata.
func StartAdvertise(cfg config.Config, pairingEnabled bool) (*Advertiser, error) {
	txt := []string{
		"node_id=" + cfg.NodeID,
		"name=" + cfg.NodeName,
		"version=" + version.Version,
	}
	if pairingEnabled {
		txt = append(txt, "pairing=true")
	} else {
		txt = append(txt, "pairing=false")
	}
	svc, err := zeroconf.Register(
		cfg.NodeName,
		config.ServiceType,
		config.ServiceDomain,
		cfg.InternalPort,
		txt,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("mdns advertise: %w", err)
	}
	return &Advertiser{server: svc}, nil
}

// Stop shuts down advertisement.
func (a *Advertiser) Stop() {
	if a != nil && a.server != nil {
		a.server.Shutdown()
	}
}
