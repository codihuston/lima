// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package usernet

import (
	"testing"

	"github.com/goccy/go-yaml"
	"github.com/lima-vm/lima/v2/pkg/networks"
	"gotest.tools/v3/assert"
)

func TestGatewayConfiguration(t *testing.T) {
	var entry networks.Network
	assert.NilError(t, yaml.UnmarshalWithOptions([]byte("mode: user-v2\ngateway: 192.168.200.1\nnetmask: 255.255.255.0\nblockAllOutbound: true\ngatewayAllowedPorts: [8080, 8443]\n"), &entry, yaml.Strict()))
	assert.DeepEqual(t, entry.GatewayAllowedPorts, []int{8080, 8443})
	config, err := netstackConfig(&GVisorNetstackOpts{Subnet: "192.168.200.0/24", BlockAllOutbound: entry.BlockAllOutbound, GatewayAllowedPorts: entry.GatewayAllowedPorts})
	assert.NilError(t, err)
	assert.DeepEqual(t, config.GatewayAllowedPorts, entry.GatewayAllowedPorts)
	assert.Assert(t, config.BlockAllOutbound)
	assert.Equal(t, config.NAT[config.GatewayIP], "127.0.0.1")
	assert.DeepEqual(t, GatewayArgs(entry.GatewayAllowedPorts), []string{"--gateway-allowed-port", "8080", "--gateway-allowed-port", "8443"})
}

func TestGatewayConfigDefaultNoneAndInvalid(t *testing.T) {
	config, err := netstackConfig(&GVisorNetstackOpts{Subnet: "192.168.200.0/24"})
	assert.NilError(t, err)
	assert.Equal(t, len(config.GatewayAllowedPorts), 0)
	for _, port := range []int{-1, 0, 65536} {
		_, err := netstackConfig(&GVisorNetstackOpts{Subnet: "192.168.200.0/24", GatewayAllowedPorts: []int{port}})
		assert.ErrorContains(t, err, "gatewayAllowedPorts")
	}
}
