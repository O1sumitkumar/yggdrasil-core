package discovery

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestProbeStaticPeers(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/internal/v1/node", func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]string{
			"id": "peer-1", "name": "Peer", "cert_pem": "x",
		})
	})
	srv := httptest.NewServer(mux)
	defer srv.Close()

	host := srv.Listener.Addr().String()
	found := ProbeStaticPeers(context.Background(), []string{host}, "local")
	if len(found) != 1 {
		t.Fatalf("found=%d", len(found))
	}
	if found[0].Node.ID != "peer-1" || found[0].Node.Address != host {
		t.Fatalf("%+v", found[0].Node)
	}
}
