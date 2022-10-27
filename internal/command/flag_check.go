package command

import "github.com/urfave/cli/v2"

func existFlags(ctx *cli.Context, flags ...string) bool {
	for _, f := range flags {
		if !ctx.IsSet(f) {
			return false
		}
	}
	return true
}
