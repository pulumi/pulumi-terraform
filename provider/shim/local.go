package shim

import (
	"context"

	"github.com/zclconf/go-cty/cty"
)

type LocalStateReferenceInputs struct {
	Path string
}

func LocalStateReferenceRead(ctx context.Context, args LocalStateReferenceInputs) (map[string]any, error) {
	return StateReferenceRead(ctx, "local", "", map[string]cty.Value{
		"path": cty.StringVal(args.Path),
	})
}
