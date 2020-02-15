package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
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
	Active  []string,
	Passive []string,
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

func (c *client) peerStat() *peerStat {
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

func peerStatBytes(p *wireStat) (bs []byte) {
	buf := bytes.NewBuffer(make([]byte, 0, STAT_SIZE))
	binary.Write(buf, binary.LittleEndian, p)
	return buf.Bytes()
}

func bytesPeerStat(bs []byte) *wireStat {
	buf := bytes.NewReader(bs)
	p := &wireStat{}
	binary.Read(buf, binary.LittleEndian, p)
	return p
}

func (p *wireStat) MarshalJSON() ([]byte, error) {
	var bs [ADDR_SIZE]byte
	var addrs []string

	for i := 0; i <= ACTIVE_SIZE; i += ADDR_SIZE {
		j := i + ADDR_SIZE
		copy(bs[:], p.Active[i:j])
		str := parseAddr(bs)
		if str != "" {
			addrs = append(addrs, str)
		}
	}

	return json.Marshal

	return nil, nil
}

func (p *wireStat) UnmarshalJSON() ([]byte, error) {
	return nil, nil
}
