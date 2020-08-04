package interfaces_test

import (
	"reflect"
	"testing"

	"github.com/belldata-dx/bdx/interfaces"
	"github.com/stretchr/testify/assert"
)

func TestInrerfaces(t *testing.T) {
	var handler interfaces.BdxHandlerFunc
	handler = func(interfaces.Context) {}
	handlers := interfaces.HandlersChain{
		func(interfaces.Context) {},
		handler,
	}
	last := handlers.Last()
	val1 := reflect.ValueOf(handler)
	val2 := reflect.ValueOf(last)
	assert.Equal(t, val1, val2)
}
