// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package networks

import "net"

type Config struct {
	Paths    Paths              `yaml:"paths" json:"paths"`
	Group    string             `yaml:"group,omitempty" json:"group,omitempty"` // default: "admin"
	Networks map[string]Network `yaml:"networks" json:"networks"`
}

type Paths struct {
	SocketVMNet string `yaml:"socketVMNet" json:"socketVMNet"`
	VarRun      string `yaml:"varRun" json:"varRun"`
	Sudoers     string `yaml:"sudoers,omitempty" json:"sudoers,omitempty"`
}

const (
	ModeUserV2  = "user-v2"
	ModeHost    = "host"
	ModeShared  = "shared"
	ModeBridged = "bridged"
)

var Modes = []string{
	ModeUserV2,
	ModeHost,
	ModeShared,
	ModeBridged,
}

type Network struct {
	Mode      string `yaml:"mode" json:"mode"`                               // "user-v2", "host", "shared", or "bridged"
	Interface string `yaml:"interface,omitempty" json:"interface,omitempty"` // only used by "bridged" networks
	Gateway   net.IP `yaml:"gateway,omitempty" json:"gateway,omitempty"`     // only used by "user-v2", "host" and "shared" networks
	DHCPEnd   net.IP `yaml:"dhcpEnd,omitempty" json:"dhcpEnd,omitempty"`     // default: same as Gateway, last byte is 254
	NetMask   net.IP `yaml:"netmask,omitempty" json:"netmask,omitempty"`     // default: 255.255.255.0

	// OutboundAllow is a list of regex patterns for allowed outbound domains, threaded
	// into gvproxy's Configuration.OutboundAllow. Only used by "user-v2" networks.
	OutboundAllow []string `yaml:"outboundAllow,omitempty" json:"outboundAllow,omitempty"`

	// BlockAllOutbound blocks all guest-initiated outbound TCP/UDP connections,
	// threaded into gvproxy's Configuration.BlockAllOutbound. Only used by "user-v2" networks.
	BlockAllOutbound bool `yaml:"blockAllOutbound,omitempty" json:"blockAllOutbound,omitempty"`

	// GatewayPortAllow lists the host ports on the gateway address that
	// remain reachable from the guest while OutboundAllow is active,
	// threaded into gvproxy's Configuration.GatewayPortAllow. Default:
	// empty, meaning no port on the gateway address is reachable. Only used
	// by "user-v2" networks.
	GatewayPortAllow []int `yaml:"gatewayPortAllow,omitempty" json:"gatewayPortAllow,omitempty"`
}
