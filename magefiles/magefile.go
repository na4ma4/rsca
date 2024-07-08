//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	//mage:import
	"github.com/dosquad/mage"
)

// Local update, protoc, format, tidy, lint & test.
func Local(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Protobuf.GenGo)
	mg.CtxDeps(ctx, mage.Protobuf.GenGoGRPC)
	mg.CtxDeps(ctx, mage.Test)
}

var Default = Local
