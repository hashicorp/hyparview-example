package main

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddrBinary(t *testing.T) {
	// Encode a string
	a := "127.0.0.1:4000"
	buf := bytes.NewBuffer(make([]byte, 0, ADDR_SIZE))
	writeAddr(buf, a)

	// Decode 6 bytes
	arr := [ADDR_SIZE]byte{}
	copy(arr[:], buf.Bytes()[0:6])
	b := parseAddr(arr)
	require.Equal(t, a, b)
}

func TestAddrBinaryEmpty(t *testing.T) {
	arr := [ADDR_SIZE]byte{0}
	b := parseAddr(arr)
	require.Equal(t, "", b)
}

func TestStatBinary(t *testing.T) {
	// Empty
	p := &wireStat{}
	bs := p.Bytes()
	require.Equal(t, 228, len(bs))

	// Fake full client
	c := newClient(&clientConfig{
		id:        newID(),
		addr:      "127.0.0.1:4000",
		bootstrap: "127.0.0.1:4000",
	})
	for i := 1; i <= 5; i++ {
		c.hv.AddActive(node(fmt.Sprintf("127.0.0.1:40%02d", i)))
	}
	for i := 6; i <= 36; i++ {
		c.hv.AddPassive(node(fmt.Sprintf("127.0.0.1:40%02d", i)))
	}

	// Bidirectional encoding
	p1 := c.wireStat()
	p1.App = 5
	p1.Hops = 7
	p1.Waste = 9
	require.Equal(t, 228, len(p1.Bytes()))

	p2 := &wireStat{}
	p2.Parse(p1.Bytes())
	require.Equal(t, p1, p2)

	// Pretty print conversion for json encoding
	pp := p2.peerStat()
	require.Equal(t, "127.0.0.1:4001", pp.Active[0])
	require.Equal(t, "127.0.0.1:4005", pp.Active[4])
	require.Equal(t, "127.0.0.1:4006", pp.Passive[0])
	require.Equal(t, "127.0.0.1:4036", pp.Passive[29])
	require.Equal(t, int32(5), pp.App)
	require.Equal(t, int32(7), pp.Hops)
	require.Equal(t, int32(9), pp.Waste)
}
