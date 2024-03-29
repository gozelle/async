package parallel

import (
	"context"
)

func NewGroup() *Group {
	return &Group{}
}

type Group struct {
	cancel func()
	fns    []func() error
	limit  uint
	ctx    context.Context
}

func (g *Group) SetContext(ctx context.Context) {
	g.ctx = ctx
}

func (g *Group) SetLimit(limit uint) {
	g.limit = limit
}

func (g *Group) Wait() error {
	var runners []Runner[Null]
	for _, v := range g.fns {
		f := v
		runners = append(runners, func(ctx context.Context) (Null, error) {
			return nil, f()
		})
	}
	if g.limit == 0 {
		g.limit = 10
	}
	res := Run[Null](g.ctx, g.limit, runners)
	return Wait[Null](res, nil)
}

func (g *Group) Go(f func() error) {
	g.fns = append(g.fns, f)
}
