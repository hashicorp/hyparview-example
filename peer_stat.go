package main

import (
	"bytes"
	"encoding/binary"
	"strconv"
	"strings"
)

const (
	ACTIVESIZE  = 30  // 5 * 6 bytes
	PASSIVESIZE = 180 // 30 * 6 bytes
)

// peerStat has all the stats data, we're going to collect them from the cluster
type peerStat struct {
	Active  [ACTIVESIZE]byte
	Passive [PASSIVESIZE]byte
	From    byte
	App     int32
	Hops    int32
	Waste   int32
}

func (c *client) makePeerStat() *peerStat {
	m := &peerStat{
		App:   c.app.Value,
		Hops:  c.app.Hops,
		Waste: c.app.Waste,
	}

	// From
	buf := bytes.NewBuffer(make([]byte, 0, 6))
	writeAddr(buf, c.hv.Self.Addr)
	m.From = buf.Bytes()

	// Active
	buf = bytes.NewBuffer(make([]byte, 0, ACTIVESIZE))
	for _, node := range c.hv.Active.Nodes {
		writeAddr(buf, node.Addr)
	}
	m.Active = buf.Bytes()

	// Passive
	buf = bytes.NewBuffer(make([]byte, 0, PASSIVESIZE))
	for _, node := range c.hv.Passive.Nodes {
		writeAddr(buf, node.Addr)
	}
	m.Passive = buf.Bytes()

	return m
}

func writeAddr(bytes []byte, addr string) {
	ip, port := strings.Split(addr, ":")
	for _, b := range strings.Split(ip, ".") {
		int := uint8(strconv.ParseUint(b))
		binary.Write(buf, binary.LittleEndian, int)
	}
	int := uint16(strconv.ParseUint(port))
	binary.Write(buf, binary.LittleEndian, int)
}

func bytesPeerStatMessage(bs []byte) (p *peerStat) {
	buf := bytes.NewReader(bs)
	p := &peerStat{}
	binary.Read(buf, binary.LittleEndian, p)
	return p
}
