package main

import (
	"context"
	"fmt"
	"os"

	"github.com/oseau/web/cmd/http"
)

// check https://npf.io/2016/10/reusable-commands/
func run() int {
	ctx := context.Background()
	if err := http.Run(ctx, os.Stdout); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		return 1
	}
	return 0
}
