package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

const rootDir = "./testdir"

func TestWalker(t *testing.T) {
	metadata := walker(rootDir)

	for el := range metadata {
		fmt.Println(el)
	}

	require.True(t, false)
}
