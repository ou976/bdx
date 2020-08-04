package main

import (
	"github.com/belldata-dx/bdx"
	"github.com/belldata-dx/bdx/interfaces"
)

func main() {
	r := bdx.New()
	r.GET("/health", func(c interfaces.Context) {
		c.JSON(200, bdx.B{"data": "test"})
	})
	r.Run()
}
