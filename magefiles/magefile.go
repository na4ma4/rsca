//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	//mage:import
	"github.com/dosquad/mage"
)

// Test update, protoc, format, tidy, lint & test.
func Test(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Protobuf.GenGo)
	mg.CtxDeps(ctx, mage.Protobuf.GenGoGRPC)
	mg.CtxDeps(ctx, mage.Golang.Lint)
	mg.CtxDeps(ctx, mage.Golang.Test)
}

var Default = Test
