//go:build mage

package main

import (
	"context"

	"github.com/magefile/mage/mg"

	//mage:import
	"github.com/dosquad/mage"
)

// TestLocal update, protoc, format, tidy, lint & test.
func TestLocal(ctx context.Context) {
	mg.CtxDeps(ctx, mage.Test)
}

var Default = TestLocal
