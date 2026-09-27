package auth

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strconv"
	"strings"
	"time"
)

const authTokenTTL = 5 * time.Minute

// AuthToken creates a short-lived Bearer token signed by this node.
// Format (base64url of): nodeID|unixExp|hex(sig)
func (id *NodeIdentity) AuthToken() string {
	exp := time.Now().UTC().Add(authTokenTTL).Unix()
	msg := fmt.Sprintf("%s|%d", id.NodeID, exp)
	sig := id.Sign([]byte(msg))
	raw := fmt.Sprintf("%s|%d|%s", id.NodeID, exp, hex.EncodeToString(sig))
	return base64.RawURLEncoding.EncodeToString([]byte(raw))
}

// PeekAuthTokenNodeID extracts the claimed node id without verifying the signature.
func PeekAuthTokenNodeID(token string) (string, error) {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(token))
	if err != nil {
		return "", fmt.Errorf("invalid token encoding")
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 || parts[0] == "" {
		return "", fmt.Errorf("invalid token format")
	}
	return parts[0], nil
}

// ParseAndVerifyAuthToken validates a Bearer token using peerCertPEM.
func ParseAndVerifyAuthToken(token string, peerCertPEM []byte) (nodeID string, err error) {
	raw, err := base64.RawURLEncoding.DecodeString(strings.TrimSpace(token))
	if err != nil {
		return "", fmt.Errorf("invalid token encoding")
	}
	parts := strings.Split(string(raw), "|")
	if len(parts) != 3 {
		return "", fmt.Errorf("invalid token format")
	}
	nodeID = parts[0]
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil {
		return "", fmt.Errorf("invalid token expiry")
	}
	if time.Now().UTC().Unix() > exp {
		return "", fmt.Errorf("token expired")
	}
	sig, err := hex.DecodeString(parts[2])
	if err != nil {
		return "", fmt.Errorf("invalid token signature")
	}
	msg := []byte(fmt.Sprintf("%s|%d", nodeID, exp))
	if !VerifyPeer(peerCertPEM, msg, sig) {
		return "", fmt.Errorf("bad signature")
	}
	return nodeID, nil
}
