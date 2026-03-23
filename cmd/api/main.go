package main

import (
	"context"

	"github.com/k1ender/psf/pkg/api"
)

func main() {
	ctx := context.Background()

	err := api.Run(ctx)
	if err != nil {
		panic(err)
	}
}
