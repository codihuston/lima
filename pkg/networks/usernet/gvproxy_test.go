// SPDX-FileCopyrightText: Copyright The Lima Authors
// SPDX-License-Identifier: Apache-2.0

package usernet

import (
	"bufio"
	"os"
	"path"
	"testing"

	"gotest.tools/v3/assert"
)

func TestSearchDomain(t *testing.T) {
	t.Run("search domain", func(t *testing.T) {
		resolvFile := path.Join(t.TempDir(), "resolv.conf")
		createResolveFile(t, resolvFile, `
search test.com lima.net
nameserver 192.168.0.100
nameserver 8.8.8.8`)

		dns := resolveSearchDomain(resolvFile)
		assert.DeepEqual(t, dns, []string{"test.com", "lima.net"})
	})

	t.Run("empty search domain", func(t *testing.T) {
		resolvFile := path.Join(t.TempDir(), "resolv.conf")
		createResolveFile(t, resolvFile, `
nameserver 192.168.0.100
nameserver 8.8.8.8`)

		dns := resolveSearchDomain(resolvFile)
		var expected []string
		assert.DeepEqual(t, dns, expected)
	})
}

func TestGatewayPortAllowUint16(t *testing.T) {
	t.Run("empty input returns nil", func(t *testing.T) {
		var expected []uint16
		assert.DeepEqual(t, gatewayPortAllowUint16(nil), expected)
	})

	t.Run("valid ports pass through", func(t *testing.T) {
		got := gatewayPortAllowUint16([]int{80, 443, 65535, 0})
		assert.DeepEqual(t, got, []uint16{80, 443, 65535, 0})
	})

	t.Run("out-of-range ports are dropped, not truncated", func(t *testing.T) {
		got := gatewayPortAllowUint16([]int{80, -1, 65536, 443})
		assert.DeepEqual(t, got, []uint16{80, 443})
	})
}

func createResolveFile(t *testing.T, file, content string) {
	f, err := os.Create(file)
	assert.NilError(t, err)
	t.Cleanup(func() { _ = f.Close() })
	writer := bufio.NewWriter(f)
	_, err = writer.WriteString(content)
	assert.NilError(t, err)
	err = writer.Flush()
	assert.NilError(t, err)
}
