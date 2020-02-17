package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strconv"
	"strings"
)

const (
	ADDR_SIZE    = 6
	ACTIVE_SIZE  = 5 * ADDR_SIZE
	PASSIVE_SIZE = 30 * ADDR_SIZE
	STAT_SIZE    = ACTIVE_SIZE + PASSIVE_SIZE + ADDR_SIZE + 5*3
)

// peerStat has all the stats data, we're going to collect them from the cluster
type peerStat struct {
	Active  []string
	Passive []string
	From    string
	App     int32
	Hops    int32
	Waste   int32
}

// wireStat has all the stats data, we're going to collect them from the cluster
type wireStat struct {
	Active  [ACTIVE_SIZE]byte
	Passive [PASSIVE_SIZE]byte
	From    [ADDR_SIZE]byte
	App     int32
	Hops    int32
	Waste   int32
}

func (c *client) wireStat() *wireStat {
	m := &wireStat{
		Active:  [ACTIVE_SIZE]byte{},
		Passive: [PASSIVE_SIZE]byte{},
		From:    [ADDR_SIZE]byte{},
		App:     c.app.Value,
		Hops:    c.app.Hops,
		Waste:   c.app.Waste,
	}

	// From
	buf := bytes.NewBuffer(make([]byte, 0, ADDR_SIZE))
	writeAddr(buf, c.hv.Self.Addr)
	copy(m.From[:], buf.Bytes()[0:ADDR_SIZE])

	// Active
	buf = bytes.NewBuffer(make([]byte, 0, ACTIVE_SIZE))
	for _, node := range c.hv.Active.Nodes {
		writeAddr(buf, node.Addr)
	}
	copy(m.Active[:], buf.Bytes()[0:ACTIVE_SIZE])

	// Passive
	buf = bytes.NewBuffer(make([]byte, 0, PASSIVE_SIZE))
	for _, node := range c.hv.Passive.Nodes {
		writeAddr(buf, node.Addr)
	}
	copy(m.Passive[:], buf.Bytes()[0:PASSIVE_SIZE])

	return m
}

func writeAddr(buf *bytes.Buffer, addr string) {
	ip_port := strings.Split(addr, ":")
	if len(ip_port) < 2 {
		return
	}
	ip, port := ip_port[0], ip_port[1]

	for _, b := range strings.Split(ip, ".") {
		int, _ := strconv.ParseUint(b, 10, 8)
		binary.Write(buf, binary.LittleEndian, uint8(int))
	}

	int, _ := strconv.ParseUint(port, 10, 16)
	binary.Write(buf, binary.LittleEndian, uint16(int))
}

func parseAddr(bs [ADDR_SIZE]byte) string {
	var ip []string
	for i := 0; i < 4; i++ {
		s := strconv.Itoa(int(bs[i]))
		ip = append(ip, s)
	}

	out := strings.Join(ip, ".")
	p := binary.LittleEndian.Uint16(bs[4:6])
	out = fmt.Sprintf("%s:%d", out, p)

	if out == "0.0.0.0:0" {
		return ""
	}

	return out
}

func (p *wireStat) Bytes() []byte {
	buf := bytes.NewBuffer(make([]byte, 0, STAT_SIZE))
	binary.Write(buf, binary.LittleEndian, p)
	return buf.Bytes()
}

func (p *wireStat) Parse(bs []byte) {
	buf := bytes.NewReader(bs)
	binary.Read(buf, binary.LittleEndian, p)
}

func (p *wireStat) peerStat() *peerStat {
	var addr [ADDR_SIZE]byte
	copy(addr[:], p.From[0:6])
	from := parseAddr(addr)

	addrs := func(size int, src []byte) []string {
		var out []string
		for j := ADDR_SIZE; j <= size; j += ADDR_SIZE {
			i := j - ADDR_SIZE
			copy(addr[:], src[i:j])
			str := parseAddr(addr)
			if str != "" {
				out = append(out, str)
			}
		}
		return out
	}

	return &peerStat{
		Active:  addrs(ACTIVE_SIZE, p.Active[:]),
		Passive: addrs(PASSIVE_SIZE, p.Passive[:]),
		From:    from,
		App:     p.App,
		Hops:    p.Hops,
		Waste:   p.Waste,
	}
}
