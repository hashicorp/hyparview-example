package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlices(t *testing.T) {
	ss := []string{"foo", "bar", "baz"}
	ns := sliceAddrNode(ss)
	ss2 := sliceNodeAddr(ns)
	require.Equal(t, ss, ss2)
}
