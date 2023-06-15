package magic

type ComponentFn[Props any] func(s Socket, p Props) AppliedView
type DeferredComponentFn[Props any] func(s Socket, p Props) func() AppliedView

func Component[Props any](compfn ComponentFn[Props]) ComponentFn[Props] {
	return func(s Socket, p Props) AppliedView {
		sref := s.clone()
		res := compfn(sref, p)
		return res
	}
}

func DeferredComponent[Props any](compfn ComponentFn[Props]) DeferredComponentFn[Props] {
	return func(s Socket, p Props) func() AppliedView {
		return func() AppliedView {
			sref := s.clone()
			res := compfn(sref, p)
			return res
		}
	}
}
