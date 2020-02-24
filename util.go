package main

import h "github.com/hashicorp/hyparview"

func node(addr string) *h.Node {
	return &h.Node{ID: addr, Addr: addr}
}

func sliceNodeAddr(ns []*h.Node) []string {
	ss := make([]string, len(ns))
	for i, n := range ns {
		ss[i] = n.Addr
	}
	return ss
}

func sliceAddrNode(ss []string) []*h.Node {
	ns := make([]*h.Node, len(ss))
	for i, n := range ss {
		ns[i] = &h.Node{Addr: n}
	}
	return ns
}
