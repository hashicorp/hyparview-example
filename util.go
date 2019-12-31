package main

import h "github.com/hashicorp/hyparview"

func node(addr string) *h.Node {
	return &h.Node{Addr: addr}
}

func sliceNodeAddr(ns []*h.Node) (ss []string) {
	for _, n := range ns {
		ss = append(ss, n.Addr)
	}
	return ss
}

func sliceAddrNode(ss []string) (ns []*h.Node) {
	for _, n := range ss {
		ns = append(ns, &h.Node{Addr: n})
	}
	return ns
}
