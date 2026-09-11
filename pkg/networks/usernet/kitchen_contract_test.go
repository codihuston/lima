// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package usernet

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/containers/gvisor-tap-vsock/pkg/types"
	"github.com/goccy/go-yaml"
	"github.com/lima-vm/lima/v2/pkg/networks"
	"gotest.tools/v3/assert"
)

// make fork-test supplies fresh output of Kitchen's production EnsureEntry.
// No VM or shared Lima state is involved.
func TestKitchenRenderedNetwork(t *testing.T) {
	dir := os.Getenv("KITCHEN_FORK_FIXTURES")
	if dir == "" {
		t.Skip("run Kitchen make fork-test for the cross-repository contract")
	}
	data, err := os.ReadFile(filepath.Join(dir, "_config", "networks.yaml"))
	assert.NilError(t, err)
	var cfg networks.Config
	assert.NilError(t, yaml.UnmarshalWithOptions(data, &cfg, yaml.Strict()))
	var configs = map[string]*types.Configuration{}
	for name, entry := range cfg.Networks {
		assert.Equal(t, entry.Mode, networks.ModeUserV2)
		_, subnet, err := netmaskToCidr(entry.Gateway, entry.NetMask)
		assert.NilError(t, err)
		config, err := netstackConfig(&GVisorNetstackOpts{MTU: 1500, Subnet: subnet.String(), OutboundAllow: entry.OutboundAllow, BlockAllOutbound: entry.BlockAllOutbound, GatewayAllowedPorts: entry.GatewayAllowedPorts})
		assert.NilError(t, err)
		configs[name] = config
	}
	for _, name := range []string{"kitchen-filter", "kitchen-deny", "kitchen-default"} {
		assert.Assert(t, configs[name] != nil, fmt.Sprintf("missing %s", name))
	}
	data, err = json.Marshal(configs)
	assert.NilError(t, err)
	assert.NilError(t, os.WriteFile(filepath.Join(dir, "stack-configs.json"), data, 0o600))
}
