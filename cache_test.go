package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	cache := initCache("/home/hawk/asd.abc")
	cache.store("a")
	require.True(t, cache.exists("a"))
	require.False(t, cache.exists("b"))
}
