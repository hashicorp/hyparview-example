package main

import (
	"bytes"
	"fmt"
	"testing"

	h "github.com/hashicorp/hyparview"
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
	bs := peerStatBytes(p)
	require.Equal(t, 228, len(bs))

	// Fake full client
	c := newClient(&clientConfig{
		id:        newID(),
		addr:      "127.0.0.1:4000",
		bootstrap: "127.0.0.1:4000",
	})
	for i := 1; i <= 5; i++ {
		c.hv.AddActive(&h.Node{Addr: fmt.Sprintf("127.0.0.1:40%02d", i)})
	}
	for i := 6; i <= 36; i++ {
		c.hv.AddPassive(&h.Node{Addr: fmt.Sprintf("127.0.0.1:40%02d", i)})
	}

	// Bidirectional encoding
	p1 := c.makePeerStat()
	require.Equal(t, 228, len(peerStatBytes(p1)))

	p2 := bytesPeerStat(peerStatBytes(p1))
	require.Equal(t, p1, p2)
}
