package magic

type ComponentFn[Props any] func(ctx Context, props Props) Node

func Component[Props any](compfn ComponentFn[Props]) ComponentFn[Props] {
	return func(ctx Context, props Props) Node {
		subContext := ctx.clone()
		node := compfn(subContext, props)
		return node
	}
}
