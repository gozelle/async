package group

import (
	"context"
	"golang.org/x/sync/errgroup"
)

func NewGroup() *Group {
	return &Group{}
}

func WithContext(ctx context.Context) (*Group, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	return (*Group)(g), ctx
}

type Group errgroup.Group

func (g *Group) Wait() error {
	return (*errgroup.Group)(g).Wait()
}

func (g *Group) Go(f func() error) {
	(*errgroup.Group)(g).Go(f)
}

func (g *Group) TryGo(f func() error) bool {
	return (*errgroup.Group)(g).TryGo(f)
}

func (g *Group) SetLimit(n int) {
	(*errgroup.Group)(g).SetLimit(n)
}
