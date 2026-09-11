// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"testing"

	"github.com/lima-vm/lima/v2/pkg/networks/usernet"
	"gotest.tools/v3/assert"
)

func TestUsernetGatewayFlags(t *testing.T) {
	cmd := newUsernetCommand()
	assert.NilError(t, cmd.ParseFlags(usernet.GatewayArgs([]int{8080, 8443})))
	ports, err := cmd.Flags().GetIntSlice("gateway-allowed-port")
	assert.NilError(t, err)
	assert.DeepEqual(t, ports, []int{8080, 8443})
}
