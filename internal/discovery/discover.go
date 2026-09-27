package discovery

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/grandcat/zeroconf"
	"github.com/yeixio/yggdrasil-core/internal/config"
	"github.com/yeixio/yggdrasil-core/internal/events"
	"github.com/yeixio/yggdrasil-core/pkg/contracts"
)

// DiscoveredNode is a node found via mDNS.
type DiscoveredNode struct {
	Node    contracts.Node
	Pairing bool
	TXT     map[string]string
}

// Discover browses for peer nodes until ctx cancelled.
func Discover(ctx context.Context, bus *events.Bus, localNodeID string) ([]DiscoveredNode, error) {
	resolver, err := zeroconf.NewResolver(nil)
	if err != nil {
		return nil, err
	}
	entries := make(chan *zeroconf.ServiceEntry)
	var found []DiscoveredNode

	go func() {
		for entry := range entries {
			meta := parseTXT(entry.Text)
			nodeID := meta["node_id"]
			if nodeID == "" || nodeID == localNodeID {
				continue
			}
			addr := ""
			if len(entry.AddrIPv4) > 0 {
				addr = entry.AddrIPv4[0].String()
			} else if len(entry.AddrIPv6) > 0 {
				addr = entry.AddrIPv6[0].String()
			}
			port := entry.Port
			if port <= 0 {
				port = config.DefaultInternalPort
			}
			if addr != "" {
				addr = net.JoinHostPort(addr, strconv.Itoa(port))
			}
			n := contracts.Node{
				ID:      nodeID,
				Name:    meta["name"],
				Status:  contracts.NodeStatusOnline,
				Paired:  false,
				Address: addr,
			}
			pairing := meta["pairing"] == "true"
			found = append(found, DiscoveredNode{Node: n, Pairing: pairing, TXT: meta})
			if bus != nil {
				bus.Publish(events.New(events.NodeDiscovered, map[string]any{
					"node_id": nodeID,
					"name":    n.Name,
					"address": addr,
				}))
			}
		}
	}()

	browseCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()
	err = resolver.Browse(browseCtx, config.ServiceType, config.ServiceDomain, entries)
	if err != nil && ctx.Err() == nil {
		return found, err
	}
	time.Sleep(500 * time.Millisecond)
	return found, nil
}

func parseTXT(lines []string) map[string]string {
	out := make(map[string]string)
	for _, line := range lines {
		k, v, ok := strings.Cut(line, "=")
		if ok {
			out[k] = v
		}
	}
	return out
}
