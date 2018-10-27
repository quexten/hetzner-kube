package clustermanager_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/xetys/hetzner-kube/pkg/clustermanager"
)

func TestGenerateWireguardConf(t *testing.T) {
	nodes := []clustermanager.Node{
		{Name: "node1", IPAddress: "1.1.1.1", PrivateIPAddress: "10.0.0.1", WireGuardKeyPair: clustermanager.WgKeyPair{Private: "node1priv", Public: "node1pub"}},
		{Name: "node2", IPAddress: "1.1.1.2", PrivateIPAddress: "10.0.0.2", WireGuardKeyPair: clustermanager.WgKeyPair{Private: "node2priv", Public: "node2pub"}},
	}

	expectedConf := `[Interface]
Address = 10.0.0.2
PrivateKey = node2priv
ListenPort = 51820

# node1
[Peer]
PublicKey = node1pub
AllowedIps = 10.0.0.1/32
Endpoint = 1.1.1.1:51820
`

	generatedConf := clustermanager.GenerateWireguardConf(nodes[1], nodes)

	if generatedConf != expectedConf {
		t.Errorf("The file was not rendered as expected\n%s\n\n", generatedConf)
	}

}

func TestGenerateKeyPair(t *testing.T) {
	wgKey, err := clustermanager.GenerateKeyPair()
	if err != nil {
		t.Errorf("Unable to generate keypairs")
	}

	if wgKey.Private == "" {
		t.Errorf("Private key is not correctly set")
	}

	if wgKey.Public == "" {
		t.Errorf("Public key is not correctly set")
	}

	privateBytes, err := base64.StdEncoding.DecodeString(wgKey.Private)
	if err != nil {
		t.Errorf("Private key is not correctly Base64 encoded")
	}

	if len(privateBytes) != 32 {
		t.Errorf("Private key is not 32 bytes len")
	}

	publicBytes, err := base64.StdEncoding.DecodeString(wgKey.Public)
	if err != nil {
		t.Errorf("Public key is not correctly Base64 encoded")
	}

	if len(publicBytes) != 32 {
		t.Errorf("Public key is not 32 bytes len")
	}
}

func TestPrivateIPPrefix(t *testing.T) {
	testCases := []struct {
		source   string
		expected string
	}{
		{
			source:   "10.5.3.6",
			expected: "10.5.3",
		},
		{
			source:   "10.20.30.60",
			expected: "10.20.30",
		},
		{
			source:   "250.251.252.253",
			expected: "250.251.252",
		},
	}
	for _, tC := range testCases {
		t.Run(fmt.Sprintf("testing IP: %s", tC.source), func(t *testing.T) {
			generated := clustermanager.PrivateIPPrefix(tC.source)
			if tC.expected != generated {
				t.Errorf("\nParsing IP: %s\nExpected: %s\nGenerated: %s\n", tC.source, tC.expected, generated)
			}
		})
	}
}
